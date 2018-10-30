package main

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/afex/hystrix-go/hystrix/metric_collector"
	ginHystrixMiddleware "github.com/eldad87/go-boilerplate/src/pgk/gin/middleware"
	"github.com/ibm-developer/generator-ibm-core-golang-gin/generators/app/templates/plugins"

	"github.com/Bose/go-gin-opentracing"
	ginController "github.com/eldad87/go-boilerplate/src/internal/http/gin"
	jaegerLogrus "github.com/eldad87/go-boilerplate/src/pgk/uber/jaeger-client-go/log/logrus"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

	machinery "github.com/RichardKnop/machinery/v1"
	machineryConfigBuilder "github.com/RichardKnop/machinery/v1/config"
	machineryLog "github.com/RichardKnop/machinery/v1/log"

	"github.com/eldad87/go-gin-prometheus"
	"github.com/evalphobia/logrus_sentry"
	"github.com/gin-contrib/gzip"

	ginlogrus "github.com/Bose/go-gin-logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-gonic/gin"

	healthcheck "github.com/heptiolabs/healthcheck"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/promrus"

	"github.com/eldad87/go-boilerplate/src/config"
	"os"
	"time"
)

func main() {
	/*
	 * PreRequisite: Config
	 * **************************** */
	conf, err := config.GetConfig(os.Getenv("BUILD_ENV"), nil)
	if err != nil {
		panic(err) // Nothing we can do
	}

	/*
	 * PreRequisite: Hystrix
	 * **************************** */
	// Prometheus
	collector := plugins.InitializePrometheusCollector(plugins.PrometheusCollectorConfig{
		Namespace: conf.GetString("app.name"),
	})
	metricCollector.Registry.Register(collector.NewPrometheusCollector)

	/*
	 * PreRequisite: Health Check + Prometheus gauge
	 * **************************** */
	healthChecker := healthcheck.NewMetricsHandler(prometheus.DefaultRegisterer, "health_check")

	/*
	 * PreRequisite: Logrus
	 * **************************** */
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	logLevel, _ := log.ParseLevel(conf.GetString("log.level"))
	log.SetLevel(logLevel)

	// Prometheus
	hook := promrus.MustNewPrometheusHook()
	log.AddHook(hook)

	// Sentry
	if conf.GetString("sentry.dsn") != "" {
		hook, err := logrus_sentry.NewWithTagsSentryHook(conf.GetString("sentry.dsn"),
			map[string]string{"app.name": conf.GetString("app.name")},
			[]log.Level{log.PanicLevel, log.FatalLevel, log.ErrorLevel})

		if err != nil {
			log.Error("Error integrating Sentry: " + err.Error())
			healthChecker.AddReadinessCheck("sentry", func() error { return err }) // Permanent, take us down.
		} else {
			// Enabling Stacktraces
			hook.StacktraceConfiguration.Enable = true
			log.AddHook(hook)
		}
	}

	/*
	 * PreRequisite: Jaeger
	 * **************************** */
	jaegConfig := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: conf.GetString("opentracing.jaeger.host") + ":" + conf.GetString("opentracing.jaeger.port"),
		},
	}

	closer, err := jaegConfig.InitGlobalTracer(conf.GetString("app.name"), jaegercfg.Logger(jaegerLogrus.NewLogger(log.StandardLogger())))
	if err != nil {
		log.Error("Can't start jaeger: " + err.Error())
		healthChecker.AddReadinessCheck("jaeger", func() error { return err }) // Permanent, take us down.
	}
	defer closer.Close()

	/*
	 * PreRequisite: Gin
	 * **************************** */
	ginRouter := gin.New()
	if "development" != conf.GetString("build.env") {
		gin.SetMode(gin.ReleaseMode)
	}

	// Configure hystrix, later enforced by a middleware (HystrixHandler)
	hystrix.ConfigureCommand(conf.GetString("app.name"), hystrix.CommandConfig{
		Timeout:                conf.GetInt("app.request.timeout"),
		MaxConcurrentRequests:  conf.GetInt("app.request.max_conn"),
		RequestVolumeThreshold: conf.GetInt("app.request.vol_threshold"),
		SleepWindow:            conf.GetInt("app.request.sleep_window"),
		ErrorPercentThreshold:  conf.GetInt("app.request.err_percent_threshold"),
	})

	ginRouter.Use(ginlogrus.WithTracing(log.StandardLogger(),
		false,
		time.RFC3339,
		true,
		"trace-id",
		[]byte("uber-trace-id"),
		[]byte("tracing-context")),
		gin.Recovery(),
		ginHystrixMiddleware.HystrixHandler(conf.GetString("app.name")),
		gzip.Gzip(gzip.BestCompression))

	// Health check support
	healthRoute := ginRouter.Group(conf.GetString("health_check.route.group"))
	{ // Is there a better way?
		healthRoute.GET(conf.GetString("health_check.route.live"), gin.WrapF(healthChecker.LiveEndpoint))
		healthRoute.GET(conf.GetString("health_check.route.ready"), gin.WrapF(healthChecker.ReadyEndpoint))
	}

	// Prometheus support
	ginprometheus.NewPrometheus("gin")
	h := promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{DisableCompression: true}),
	)
	ginRouter.GET(conf.GetString("prometheus.route"), func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})

	// OpenTracing
	p := ginopentracing.OpenTracer([]byte(conf.GetString("app.name") + "-"))
	ginRouter.Use(p)

	/*
	 * PreRequisite: Machinery
	 * **************************** */
	var machineryConf = &machineryConfigBuilder.Config{
		Broker:        conf.GetString("machinery.broker_dsn"),
		DefaultQueue:  conf.GetString("machinery.default_queue"),
		ResultBackend: conf.GetString("machinery.result_backend_dsn"),
		NoUnixSignals: true, // We will manage our signaling
		AMQP: &machineryConfigBuilder.AMQPConfig{
			Exchange:      conf.GetString("machinery.exchange"),
			ExchangeType:  conf.GetString("machinery.exchange_type"),
			BindingKey:    conf.GetString("machinery.binding_key"),
			PrefetchCount: conf.GetInt("machinery.prefetch_count"),
		},
	}

	server, err := machinery.NewServer(machineryConf)
	if err != nil {
		log.Error("Can't start machinery: " + err.Error())
		healthChecker.AddReadinessCheck("machinery", func() error { return err }) // Permanent, take us down.
	}

	// Configure machinery to use logrus
	machineryLog.Set(log.StandardLogger())

	// Start our consumer
	if 1 == conf.GetInt("machinery.consumer.enable") {
		server.RegisterTask("repeat", func(str string) (string, error) {
			return str, nil
		})

		worker := server.NewWorker(conf.GetString("machinery.consumer.tag"), conf.GetInt("machinery.consumer.concurrent_tasks"))
		errorsChan := make(chan error)
		worker.LaunchAsync(errorsChan)
		go func() {
			for err := range errorsChan {
				log.Error("Machinery consumer: " + err.Error())
			}
		}()
	}

	/*
	 * Gin: Handlers
	 * **************************** */
	echo := ginController.NewEcho(log.StandardLogger(), server)
	ginRouter.GET("/echo", echo.Repeat)

	ginRouter.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Listen for incoming HTTP requests
	ginRouter.Run(":" + conf.GetString("app.port"))
}

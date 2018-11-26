package main

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/afex/hystrix-go/hystrix/metric_collector"
	reHystrix "github.com/eldad87/go-boilerplate/src/pkg/concurrency/hystrix"
	ginHystrixMiddleware "github.com/eldad87/go-boilerplate/src/pkg/gin/middleware"
	"github.com/ibm-developer/generator-ibm-core-golang-gin/generators/app/templates/plugins"

	jaegerprom "github.com/jaegertracing/jaeger-lib/metrics/prometheus"
	"github.com/opentracing/opentracing-go"

	"github.com/Bose/go-gin-opentracing"
	ginController "github.com/eldad87/go-boilerplate/src/internal/http/gin"
	jaegerLogrus "github.com/eldad87/go-boilerplate/src/pkg/uber/jaeger-client-go/log/logrus"
	"github.com/uber/jaeger-client-go"

	"github.com/RichardKnop/machinery/v1"
	machineryConfigBuilder "github.com/RichardKnop/machinery/v1/config"
	machineryLog "github.com/RichardKnop/machinery/v1/log"
	machineryProducer "github.com/eldad87/go-boilerplate/src/pkg/task/producer/machinery"

	"github.com/evalphobia/logrus_sentry"
	"github.com/gin-contrib/gzip"
	"github.com/zsais/go-gin-prometheus"

	"github.com/Bose/go-gin-logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-gonic/gin"

	"github.com/heptiolabs/healthcheck"
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
	if "development" == conf.GetString("build.env") {
		conf.Debug()
	}

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
	// Add jaeger metrics and reporting to prometheus route
	logAdapt := jaegerLogrus.NewLogger(log.StandardLogger())
	factory := jaegerprom.New()
	metrics := jaeger.NewMetrics(factory, map[string]string{"lib": "jaeger"})

	//Add tracing to application
	transport, err := jaeger.NewUDPTransport(conf.GetString("opentracing.jaeger.host")+":"+conf.GetString("opentracing.jaeger.port"), 0)
	if err != nil {
		healthChecker.AddReadinessCheck("jaeger", func() error { return err }) // Permanent, take us down.
		log.Error(err)
	}

	reporter := jaeger.NewCompositeReporter(
		jaeger.NewLoggingReporter(logAdapt),
		jaeger.NewRemoteReporter(transport,
			jaeger.ReporterOptions.Metrics(metrics),
			jaeger.ReporterOptions.Logger(logAdapt),
		),
	)
	defer reporter.Close()

	sampler := jaeger.NewConstSampler(true)
	tracer, closer := jaeger.NewTracer(conf.GetString("app.name"),
		sampler,
		reporter,
		jaeger.TracerOptions.Metrics(metrics),
	)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	/*
	 * PreRequisite: Gin
	 * **************************** */
	ginRouter := gin.New()
	if "development" != conf.GetString("build.env") {
		gin.SetMode(gin.ReleaseMode)
	}

	ginRouter.Use(ginlogrus.WithTracing(log.StandardLogger(),
		false,
		time.RFC3339,
		true,
		"trace-id",
		[]byte("uber-trace-id"),
		[]byte("tracing-context")),
		gin.Recovery(),
		ginHystrixMiddleware.HystrixHandler(conf.GetString("app.name"), hystrix.CommandConfig{
			Timeout:                conf.GetInt("app.request.timeout"),
			MaxConcurrentRequests:  conf.GetInt("app.request.max_conn"),
			RequestVolumeThreshold: conf.GetInt("app.request.vol_threshold"),
			SleepWindow:            conf.GetInt("app.request.sleep_window"),
			ErrorPercentThreshold:  conf.GetInt("app.request.err_per_threshold"),
		}),
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

	// Produce
	producer := machineryProducer.NewProducer(server,
		hystrix.CommandConfig{
			Timeout:                conf.GetInt("machinery.broker.timeout"),
			MaxConcurrentRequests:  conf.GetInt("machinery.broker.max_conn"),
			RequestVolumeThreshold: conf.GetInt("machinery.broker.vol_threshold"),
			SleepWindow:            conf.GetInt("machinery.broker.sleep_window"),
			ErrorPercentThreshold:  conf.GetInt("machinery.broker.err_per_threshold"),
		},
		reHystrix.RetryConfig{Attempts: conf.GetInt("machinery.broker.retries"), Delay: conf.GetDuration("machinery.broker.retry_delay")})

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
	echo := ginController.NewEcho(log.StandardLogger(), producer)
	ginRouter.GET("/echo", echo.Repeat)

	ginRouter.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Listen for incoming HTTP requests
	ginRouter.Run(":" + conf.GetString("app.port"))
}

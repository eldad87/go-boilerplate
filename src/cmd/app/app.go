package main

import (
	"github.com/Bose/go-gin-opentracing"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/afex/hystrix-go/hystrix/metric_collector"
	"github.com/eldad87/go-boilerplate/src/app"
	"github.com/eldad87/go-boilerplate/src/app/proto"
	ginController "github.com/eldad87/go-boilerplate/src/internal/http/gin"
	reHystrix "github.com/eldad87/go-boilerplate/src/pkg/concurrency/hystrix"
	ginHystrixMiddleware "github.com/eldad87/go-boilerplate/src/pkg/gin/middleware"
	grpcGatewayError "github.com/eldad87/go-boilerplate/src/pkg/grpc-gateway/error"
	jaegerLogrus "github.com/eldad87/go-boilerplate/src/pkg/uber/jaeger-client-go/log/logrus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ibm-developer/generator-ibm-core-golang-gin/generators/app/templates/plugins"
	jaegerprom "github.com/jaegertracing/jaeger-lib/metrics/prometheus"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc"
	"net"

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

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/gin-gonic/gin"

	healthChecks "github.com/eldad87/go-boilerplate/src/pkg/healthcheck"
	"github.com/heptiolabs/healthcheck"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/promrus"

	"context"
	"github.com/eldad87/go-boilerplate/src/config"
	"os"
	"strings"
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
	// Expose CB Prometheus metrics
	metricCollector.Registry.Register(collector.NewPrometheusCollector)

	/*
	 * PreRequisite: Health Check + Expose status Prometheus metrics gauge
	 * **************************** */
	healthChecker := healthcheck.NewMetricsHandler(prometheus.DefaultRegisterer, "health_check")
	healthChecker.AddLivenessCheck("Goroutine Threshold", healthcheck.GoroutineCountCheck(conf.GetInt("health_check.goroutine_threshold")))

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

	// Expose log level Prometheus metrics
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
	factory := jaegerprom.New() // By default uses prometheus.DefaultRegisterer
	metrics := jaeger.NewMetrics(factory, map[string]string{"lib": "jaeger"})

	// Add tracing to application
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
	 * PreRequisite: Database
	 * **************************** */
	db, err := gorm.Open("mssql", "sqlserver://username:password@localhost:1433?database=dbname")
	defer db.Close()

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
		ginHystrixMiddleware.HystrixHandler("http", hystrix.CommandConfig{
			Timeout:                conf.GetInt("app.request.timeout"),
			MaxConcurrentRequests:  conf.GetInt("app.request.max_conn"),
			RequestVolumeThreshold: conf.GetInt("app.request.vol_threshold"),
			SleepWindow:            conf.GetInt("app.request.sleep_window"),
			ErrorPercentThreshold:  conf.GetInt("app.request.err_per_threshold"),
		}),
		gzip.Gzip(gzip.BestCompression),
		// a workaround in order to support swagger running on a different port/domain
		func(c *gin.Context) {
			if strings.HasPrefix(c.Request.RequestURI, "/swaggerui/swagger.json") {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
				c.Header("Access-Control-Allow-Origin", "*")
			}
		},
	)

	// Health check handlers
	healthRoute := ginRouter.Group(conf.GetString("health_check.route.group"))
	{ // Is there a better way?
		healthRoute.GET(conf.GetString("health_check.route.live"), gin.WrapF(healthChecker.LiveEndpoint))
		healthRoute.GET(conf.GetString("health_check.route.ready"), gin.WrapF(healthChecker.ReadyEndpoint))
	}

	// Prometheus handlers
	ginprometheus.NewPrometheus("gin")
	h := promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{DisableCompression: true}),
	)
	ginRouter.GET(conf.GetString("prometheus.route"), func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})

	// OpenTracing middleware
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

	// Healthcheck: Redis
	if strings.HasPrefix(conf.GetString("machinery.result_backend_dsn"), "redis") {
		healthChecker.AddLivenessCheck("Redis",
			healthcheck.Async(
				healthChecks.RedisCheck(conf.GetString("machinery.result_backend_dsn"), time.Second),
				time.Second*5)) // Check interval
	}

	// Healthcheck: AMQP
	if strings.HasPrefix(conf.GetString("machinery.broker_dsn"), "amqp") {
		healthChecker.AddLivenessCheck("AMQP",
			healthcheck.Async(
				healthChecks.AMQPCheck(conf.GetString("machinery.broker_dsn"), time.Second),
				time.Second*5)) // Check interval
	}

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
	 * PreRequisite: gRPC
	 * **************************** */
	lis, err := net.Listen("tcp", ":"+conf.GetString("app.grpc.port"))
	if err != nil || lis == nil {
		log.Error("gRPC failed to listen: %v", err)
	}
	log.Println("gRPC Listening on:", conf.GetString("app.grpc.port"))
	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()

	// Visit Service
	visitService := pb.VisitService{VisitService: &app.VisitServiceDemo{}}
	pb.RegisterVisitServiceServer(grpcServer, &visitService)

	// Start listening to gRPC requests
	go func() {
		log.Info("gRPC start to listen on: ", conf.GetString("app.grpc.port"))
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("gRPC failed to serve: %v", err)
			healthChecker.AddReadinessCheck("gRPC", func() error { return err }) // Permanent, take us down.
		}
	}()

	/*
	 * gRPC: gRPC as HTTP Gateway
	 * **************************** */
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Customize our error response
	runtime.HTTPError = grpcGatewayError.CustomHTTPError
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterVisitServiceHandlerFromEndpoint(ctx, mux, ":"+conf.GetString("app.grpc.port"), opts)
	if err != nil {
		log.Error("failed to register Visit Service: %v", err)
	}
	ginRouter.Any(conf.GetString("app.grpc.http_route_prefix")+"/*any", gin.WrapF(mux.ServeHTTP))
	// Expose our Swagger-UI .json file
	ginRouter.StaticFile("/swaggerui/swagger.json", "src/app/proto/visit_service.swagger.json")

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

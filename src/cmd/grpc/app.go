package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/TheZeroSlave/zapsentry"
	metricCollector "github.com/afex/hystrix-go/hystrix/metric_collector"
	service "github.com/eldad87/go-boilerplate/src/app/mysql"
	"github.com/eldad87/go-boilerplate/src/config"

	//grpcGatewayError "github.com/eldad87/go-boilerplate/src/pkg/grpc-gateway/error"
	grpc_status_validator "github.com/eldad87/go-boilerplate/src/pkg/grpc/middleware/status/validator.v10"
	grpc_validator "github.com/eldad87/go-boilerplate/src/pkg/grpc/middleware/validator/protoc_gen_validate"
	promZap "github.com/eldad87/go-boilerplate/src/pkg/uber/zap"
	grpcTransport "github.com/eldad87/go-boilerplate/src/transport/grpc"
	pb "github.com/eldad87/go-boilerplate/src/transport/grpc/proto"
	swaggerui "github.com/jmattheis/go-packr-swagger-ui"

	null_v4_validation "github.com/eldad87/go-boilerplate/src/pkg/validator/custom/guregu/null-v4"
	v10validator "github.com/go-playground/validator/v10"

	sqlLogger "github.com/eldad87/go-boilerplate/src/pkg/go-sql-driver/logger"
	databaseDriver "github.com/go-sql-driver/mysql"
	"github.com/gobuffalo/packr"
	migrate "github.com/rubenv/sql-migrate"

	sqlmwInterceptor "github.com/eldad87/go-boilerplate/src/pkg/ngrok/sqlmw"
	"github.com/ngrok/sqlmw"

	//"github.com/eldad87/go-boilerplate/src/pkg/crypto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/heptiolabs/healthcheck"
	"github.com/ibm-developer/generator-ibm-core-golang-gin/generators/app/templates/plugins"
	jaegerZap "github.com/jaegertracing/jaeger-client-go/log/zap"
	jaegerprom "github.com/jaegertracing/jaeger-lib/metrics/prometheus"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	/*
	 * PreRequisite: Config
	 * **************************** */
	conf, err := config.GetConfig(os.Getenv("BUILD_ENV"), nil)
	if err != nil {
		panic(err) // Nothing we can do
	}
	if conf.GetString("environment") != "production" {
		conf.Debug()
	}

	/*
	 * PreRequisite: Prometheus
	 * **************************** */
	collector := plugins.InitializePrometheusCollector(plugins.PrometheusCollectorConfig{
		Namespace: conf.GetString("app.name"),
	})
	http.Handle(conf.GetString("prometheus.route"), promhttp.Handler())

	/*
	 * PreRequisite: Hystrix
	 * **************************** */
	// Expose CB Prometheus metrics
	metricCollector.Registry.Register(collector.NewPrometheusCollector)

	/*
	 * PreRequisite: Health Check + Expose status Prometheus metrics gauge
	 * **************************** */
	healthChecker := healthcheck.NewMetricsHandler(prometheus.DefaultRegisterer, "health_check")
	healthChecker.AddLivenessCheck("Goroutine Threshold", healthcheck.GoroutineCountCheck(conf.GetInt("health_check.goroutine_threshold")))

	// Expose to HTTP
	http.HandleFunc(conf.GetString("health_check.route.group")+conf.GetString("health_check.route.live"), healthChecker.LiveEndpoint)
	http.HandleFunc(conf.GetString("health_check.route.group")+conf.GetString("health_check.route.ready"), healthChecker.ReadyEndpoint)

	/*
	 * PreRequisite: Logger
	 * **************************** */
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level.UnmarshalText([]byte(conf.GetString("log.level")))
	zapConfig.Development = conf.GetString("environment") != "production"
	// Expose log level Prometheus metrics
	hook := promZap.MustNewPrometheusHook([]zapcore.Level{zapcore.DebugLevel, zapcore.InfoLevel, zapcore.WarnLevel,
		zapcore.ErrorLevel, zapcore.FatalLevel, zapcore.PanicLevel, zapcore.DebugLevel})
	logger, _ := zapConfig.Build(zap.Hooks(hook))

	// Sentry
	if conf.GetString("sentry.dsn") != "" {
		atom := zap.NewAtomicLevel()
		err := atom.UnmarshalText([]byte(conf.GetString("sentry.log_level")))
		if err != nil {
			logger.Fatal("Failed to parse Zap-Sentry log level", zap.String("sentry.log_level", conf.GetString("sentry.log_level")))
		}

		cfg := zapsentry.Configuration{
			Level: atom.Level(), //when to send message to sentry
			Tags: map[string]string{
				"component": conf.GetString("app.name"),
			},
		}
		core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromDSN(conf.GetString("sentry.dsn")))
		//in case of err it will return noop core. so we can safely attach it
		if err != nil {
			logger.Fatal("failed to init sentry / zap")
		}
		logger = zapsentry.AttachCoreToLogger(core, logger)
	}
	defer logger.Sync()

	/*
	 * PreRequisite: Jaeger
	 * **************************** */
	// Add jaeger metrics and reporting to prometheus route
	logAdapt := jaegerZap.NewLogger(logger)
	factory := jaegerprom.New() // By default uses prometheus.DefaultRegisterer
	metrics := jaeger.NewMetrics(factory, map[string]string{"lib": "jaeger"})

	// Add tracing to application
	transport, err := jaeger.NewUDPTransport(conf.GetString("opentracing.jaeger.host")+":"+conf.GetString("opentracing.jaeger.port"), 0)
	if err != nil {
		healthChecker.AddReadinessCheck("jaeger", func() error { return err }) // Permanent, take us down.
		logger.Sugar().Errorf("%+v", err)
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
	 * PreRequisite: Other
	 * **************************** */
	//passwordHandler := crypto.Hash{}

	/*
	 * PreRequisite: DataBase
	 * **************************** */
	// Logger
	databaseDriver.SetLogger(sqlLogger.NewLogger(logger))
	// Tracer
	mysqlInterceptor := sqlmwInterceptor.Interceptor{Tracer: tracer}
	sql.Register("instrumented-mysql", sqlmw.Driver(databaseDriver.MySQLDriver{}, mysqlInterceptor))
	db, err := sql.Open("instrumented-mysql", conf.GetString("database.dsn"))
	if err != nil {
		logger.Sugar().Fatal("Database failed to listen: %v. Due to error: %v", conf.GetString("database.dsn"), err)
	}

	if err := db.Ping(); err != nil {
		logger.Sugar().Errorf("Database failed to Ping: %v. Due to error: %v", conf.GetString("database.dsn"), err)
	}
	// Our app is not ready if we can't connect to our database (`var db *sql.DB`) in <1s.
	healthChecker.AddReadinessCheck(conf.GetString("database.driver"), healthcheck.DatabasePingCheck(db, 1*time.Second))

	// Migration
	if conf.GetString("database.auto_migrate") == "on" {
		migrations := &migrate.PackrMigrationSource{
			Box: packr.NewBox("../../../src/migration"),
		}
		n, err := migrate.Exec(db, conf.GetString("database.driver"), migrations, migrate.Up)
		if err != nil {
			logger.Error("Error applying migration:", zap.Error(err))
		}
		logger.Debug("Applied migrations:", zap.Int("attempt", n))
	}

	/*
	 * PreRequisite: gRPC
	 * **************************** */
	lis, err := net.Listen("tcp", ":"+conf.GetString("app.grpc.port"))
	if err != nil || lis == nil {
		logger.Sugar().Errorf("gRPC failed to listen: %v", err)
	}
	logger.Info("gRPC is about to start listening for incoming requests", zap.String("port", conf.GetString("app.grpc.port")))

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(tracer)),
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(logger),
			grpc_recovery.StreamServerInterceptor(),
			grpc_validator.StreamServerInterceptor(),
			grpc_status_validator.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_recovery.UnaryServerInterceptor(),
			grpc_validator.UnaryServerInterceptor(),
			grpc_status_validator.UnaryServerInterceptor(),
		)),
	)
	defer grpcServer.GracefulStop()

	validator := v10validator.New()
	// Register valuer for guregu/null.v4
	null_v4_validation.RegisterSQLNullValuer(validator)

	// Visit Service
	visitService := service.NewVisitService(db, validator)
	grpcVisitServer := grpcTransport.VisitServer{VisitService: visitService}

	pb.RegisterVisitServer(grpcServer, &grpcVisitServer)

	// Start listening to gRPC requests
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("gRPC failed listening for incoming requests",
				zap.String("port", conf.GetString("app.grpc.port")),
				zap.String("error", err.Error()),
			)

			healthChecker.AddReadinessCheck("gRPC", func() error { return err }) // Permanent, take us down.
		} else {
			logger.Info("gRPC is listening for incoming requests", zap.String("port", conf.GetString("app.grpc.port")))
		}
	}()

	/*
	 * gRPC: gRPC as HTTP Gateway
	 * **************************** */
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithMetadata(
			func(ctx context.Context, r *http.Request) metadata.MD {
				span := opentracing.SpanFromContext(ctx)

				// Meta-data from span
				ctxSpan := span.Context()
				carrier := make(map[string]string)

				span.Tracer().Inject(
					ctxSpan,
					opentracing.TextMap,
					opentracing.TextMapCarrier(carrier),
				)

				return metadata.New(carrier)
			},
		),
		// Customize our error response
		// runtime.WithErrorHandler(grpcGatewayError.CustomHTTPError),
	)

	// https://github.com/grpc-ecosystem/grpc-gateway/issues/348
	muxHandlerFunc := nethttp.MiddlewareFunc(
		tracer,
		mux.ServeHTTP,
		nethttp.OperationNameFunc(func(r *http.Request) string {
			return fmt.Sprintf("HTTP-gRPC %s %s", r.Method, r.URL.String())
		}),
	)

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterVisitHandlerFromEndpoint(ctx, mux, ":"+conf.GetString("app.grpc.port"), opts)
	if err != nil {
		logger.Sugar().Errorf("Failed to register Visit Service %+v", err)
	}

	http.HandleFunc(conf.GetString("app.grpc.http_route_prefix")+"/", muxHandlerFunc)

	// Swagger
	if conf.GetString("environment") == "development" {
		// Serve swagger.json
		box := swaggerui.GetBox()
		http.Handle(conf.GetString("swagger.ui.route.group"), http.StripPrefix(conf.GetString("swagger.ui.route.group"), http.FileServer(box)))

		fs := http.FileServer(http.Dir("src/transport/grpc/proto"))
		http.Handle(conf.GetString("swagger.json.route.group")+"/", http.StripPrefix(conf.GetString("swagger.json.route.group"), fs))
	}

	/*
	 * Start listening for incoming HTTP requests
	 * **************************** */
	logger.Info("Starting on port " + conf.GetString("app.port"))
	http.ListenAndServe(":"+conf.GetString("app.port"), nil)
}

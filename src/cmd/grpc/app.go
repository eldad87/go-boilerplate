package main

import (
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix/metric_collector"
	"github.com/eldad87/go-boilerplate/src/app"
	"github.com/eldad87/go-boilerplate/src/app/proto"
	"github.com/eldad87/go-boilerplate/src/config"
	grpcGatewayError "github.com/eldad87/go-boilerplate/src/pkg/grpc-gateway/error"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
	"net/http"
	"os"
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
	 * PreRequisite: Logger
	 * **************************** */
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level.UnmarshalText([]byte(conf.GetString("log.level")))
	zapConfig.Development = conf.GetString("environment") != "production"
	logger, _ := zapConfig.Build()
	defer logger.Sync()
	// TODO: Zap + Sentry

	/*
	 * PreRequisite: Prometheus
	 * **************************** */
	collector := plugins.InitializePrometheusCollector(plugins.PrometheusCollectorConfig{
		Namespace: conf.GetString("app.name"),
	})
	http.Handle(conf.GetString("prometheus.route"), promhttp.Handler())
	// TODO: Zap

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

	// TODO: Mux:
	// Max CONN, Rate, TimeOut

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
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	defer grpcServer.GracefulStop()

	// Visit Service
	visitService := pb.VisitService{VisitService: &app.VisitServiceDemo{}}
	pb.RegisterVisitServiceServer(grpcServer, &visitService)

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
	// Customize our error response
	runtime.HTTPError = grpcGatewayError.CustomHTTPError
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
	)

	// https://github.com/grpc-ecosystem/grpc-gateway/issues/348
	muxHandlerFunc := nethttp.MiddlewareFunc(
		tracer,
		mux.ServeHTTP,
		nethttp.OperationNameFunc(func(r *http.Request) string {
			return fmt.Sprintf("HTTP-gRPC %s %s", r.Method, r.URL.String())
		}),
	)

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	err = pb.RegisterVisitServiceHandlerFromEndpoint(ctx, mux, ":"+conf.GetString("app.grpc.port"), opts)
	if err != nil {
		logger.Sugar().Errorf("Failed to register Visit Service %+v", err)
	}
	http.HandleFunc(conf.GetString("app.grpc.http_route_prefix")+"/", muxHandlerFunc)

	/*
	 * Start listening for incoming HTTP requests
	 * **************************** */
	http.ListenAndServe(":"+conf.GetString("app.port"), nil)
}

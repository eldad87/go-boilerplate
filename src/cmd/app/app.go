package main

import (
	machinery "github.com/RichardKnop/machinery/v1"
	machineryConfigBuilder "github.com/RichardKnop/machinery/v1/config"
	machineryLog "github.com/RichardKnop/machinery/v1/log"

	"github.com/eldad87/go-gin-prometheus"
	"github.com/evalphobia/logrus_sentry"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-gonic/gin"

	healthcheck "github.com/heptiolabs/healthcheck"
	log "github.com/sirupsen/logrus"

	"github.com/eldad87/go-connect/src/config"
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
	 * PreRequisite: Gin
	 * **************************** */
	ginRouter := gin.New()
	if "development" != conf.GetString("build.env") {
		gin.SetMode(gin.ReleaseMode)
	}

	// Configure middlewares: logrus, Auto recovery and Gzip support
	ginRouter.Use(ginrus.Ginrus(log.StandardLogger(), time.RFC3339, true), gin.Recovery(), gzip.Gzip(gzip.BestCompression))

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

	/*
	 * PreRequisite: Machinery
	 * **************************** */
	var machineryConf = &machineryConfigBuilder.Config{
		Broker:        conf.GetString("machinery.broker_dsn"),
		DefaultQueue:  conf.GetString("machinery.default_queue"),
		ResultBackend: conf.GetString("machinery.result_backend_dsn"),
		NoUnixSignals: true, // We will manage our signaling
		AMQP: &machineryConfigBuilder.AMQPConfig{
			Exchange:     conf.GetString("machinery.exchange"),
			ExchangeType: conf.GetString("machinery.exchange_type"),
			BindingKey:   conf.GetString("machinery.binding_key"),
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
	if true == conf.GetBool("machinery.consumer.enable") {
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
	ginRouter.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Listen for incoming HTTP requests
	ginRouter.Run(":" + conf.GetString("app.port"))
}

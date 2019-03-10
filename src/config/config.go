package config

import (
	"github.com/spf13/viper"
	"strings"
)

// explicit call to Set > flag > env > config > key/value store > default
func GetConfig(env string, confFiles map[string]string) (*viper.Viper, error) {
	conf := viper.New()
	conf.SetDefault("environment", env)

	// Defaults: App
	conf.SetDefault("app.name", "default")
	conf.SetDefault("app.port", "8080")
	conf.SetDefault("app.grpc.port", "8082")
	conf.SetDefault("app.grpc.http_route_prefix", "/v1")

	conf.SetDefault("app.request.timeout", 100)
	conf.SetDefault("app.request.max_conn", 10)
	conf.SetDefault("app.request.vol_threshold", 20)
	conf.SetDefault("app.request.sleep_window", 5000)
	conf.SetDefault("app.request.err_per_threshold", 50)

	// Defaults: Monitoring
	conf.SetDefault("log.level", "debug")
	conf.SetDefault("health_check.route.group", "/health")
	conf.SetDefault("health_check.route.live", "/live")
	conf.SetDefault("health_check.route.ready", "/ready")
	conf.SetDefault("health_check.goroutine_threshold", 100) // Alert when there are more than X Goroutine

	// Defaults: Prometheus
	conf.SetDefault("prometheus.route", "/metrics")

	// Defaults: Sentry
	conf.SetDefault("sentry.dsn", "")
	conf.SetDefault("sentry.log_level", "error")

	// Swagger
	conf.SetDefault("swagger.ui.route.group", "/swaggerui/")
	conf.SetDefault("swagger.json.route.group", "/swagger")

	// Defaults: DataBase
	conf.SetDefault("database.driver", "")
	conf.SetDefault("database.dsn", "") // If you use the MySQL driver with existing database client, you must create the client with parameter multiStatements=true:
	conf.SetDefault("database.auto_migrate", "off")

	// Defaults: Machinery
	conf.SetDefault("machinery.broker_dsn", "")

	conf.SetDefault("machinery.broker.retries", 2)
	conf.SetDefault("machinery.broker.retry_delay", 10)
	conf.SetDefault("machinery.broker.timeout", 100)
	conf.SetDefault("machinery.broker.max_conn", 10)
	conf.SetDefault("machinery.broker.vol_threshold", 20)
	conf.SetDefault("machinery.broker.sleep_window", 5000)
	conf.SetDefault("machinery.broker.err_per_threshold", 50)

	conf.SetDefault("machinery.default_queue", "")
	conf.SetDefault("machinery.result_backend_dsn", "")
	conf.SetDefault("machinery.exchange", "")
	conf.SetDefault("machinery.exchange_type", "")
	conf.SetDefault("machinery.binding_key", "")
	conf.SetDefault("machinery.consumer.enable", 0)
	conf.SetDefault("machinery.consumer.tag", "")
	conf.SetDefault("machinery.consumer.concurrent_tasks", 10)
	conf.SetDefault("machinery.consumer.prefetch_count", 1)

	// Conf Env
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "_", "__")) // APP_DATA__BASE_PASS -> app.data_base.pass
	conf.AutomaticEnv()                                              // Automatically load Env variables

	// Conf Files
	//conf.SetConfigType("yaml") 				// We're using yaml
	conf.SetConfigName(env)                   // Search for a config file that matches our environment
	conf.AddConfigPath("./src/config/" + env) // look for config in the working directory
	conf.ReadInConfig()                       // Find and read the config file

	// Read additional files
	for confFile := range confFiles {
		conf.SetConfigName(confFile)
		conf.MergeInConfig()
	}

	return conf, nil
}

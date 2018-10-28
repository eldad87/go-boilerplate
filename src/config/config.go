package config

import (
	"github.com/spf13/viper"
	"strings"
)

func GetConfig(env string, confFiles map[string]string) (*viper.Viper, error) {
	conf := viper.New()

	// Defaults: App
	conf.SetDefault("app.name", "default")
	conf.SetDefault("app.port", "8080")

	// Defaults: Monitoring
	conf.SetDefault("log.level", "debug")
	conf.SetDefault("health_check.route.group", "/health")
	conf.SetDefault("health_check.route.live", "/live")
	conf.SetDefault("health_check.route.ready", "/ready")

	// Defaults: Prometheus
	conf.SetDefault("prometheus.route", "/metrics")

	// Defaults: Sentry
	conf.SetDefault("sentry.dsn", "")

	// Defaults: DataBase
	conf.SetDefault("database.driver", "")
	conf.SetDefault("database.dsn", "")

	// Defaults: Machinery
	conf.SetDefault("machinery.broker_dsn", "")
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
	//conf.SetConfigType("yaml") 					// We're using yaml
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

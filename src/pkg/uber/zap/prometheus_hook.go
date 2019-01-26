package zap

// https://github.com/weaveworks/promrus

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap/zapcore"
)

// PrometheusHook exposes Prometheus counters for each of logrus' log levels.
type PrometheusHook struct {
	counterVec *prometheus.CounterVec
}

func NewPrometheusHook(levels []zapcore.Level) (func(zapcore.Entry) error, error) {
	counterVec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "log_messages",
		Help: "Total number of log messages.",
	}, []string{"level"})

	// Initialise counters for all supported levels:
	for _, level := range levels {
		counterVec.WithLabelValues(level.String())
	}
	// Try to unregister the counter vector, in case already registered for some reason,
	// e.g. double initialisation/configuration done by mistake by the end-user.
	prometheus.Unregister(counterVec)
	// Try to register the counter vector:
	err := prometheus.Register(counterVec)
	if err != nil {
		return nil, err
	}

	return func(entry zapcore.Entry) error {
		counterVec.WithLabelValues(entry.Level.String()).Inc()
		return nil
	}, nil
}

// MustNewPrometheusHook creates a new instance of PrometheusHook which exposes Prometheus counters for various log levels.
// Contrarily to NewPrometheusHook, it does not return any error to the caller, but panics instead.
// Use MustNewPrometheusHook if you want a less verbose hook creation. Use NewPrometheusHook if you want more control.
func MustNewPrometheusHook(levels []zapcore.Level) func(zapcore.Entry) error {
	hook, err := NewPrometheusHook(levels)
	if err != nil {
		panic(err)
	}
	return hook
}

package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

// https://github.com/opentracing-contrib/go-gorilla/blob/master/gorilla/example_test.go

var httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_duration_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func Prometheus(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
}

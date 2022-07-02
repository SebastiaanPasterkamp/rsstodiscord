package middleware

// From https://www.robustperception.io/prometheus-middleware-for-gorilla-mux

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusMiddleware returns a function that implements mux.MiddlewareFunc.
func PrometheusMiddleware(name string) mux.MiddlewareFunc {
	summary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       sanitize(name) + "_http_duration_seconds",
			Help:       "Duration of HTTP requests.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"path"},
	)
	prometheus.Unregister(summary)
	prometheus.Register(summary)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route := mux.CurrentRoute(r)
			path, _ := route.GetPathTemplate()
			timer := prometheus.NewTimer(summary.WithLabelValues(path))
			next.ServeHTTP(w, r)
			timer.ObserveDuration()
		})
	}
}

func sanitize(name string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9 ]+")
	return reg.ReplaceAllString(strings.ToLower(name), "")
}

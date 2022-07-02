package handlers

import (
	mw "github.com/SebastiaanPasterkamp/gobernate/middleware"
	"github.com/SebastiaanPasterkamp/gobernate/version"

	"sync/atomic"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Router register necessary routes and returns an instance of a router.
func Router(info version.Info, isReady *atomic.Value, shutdown chan bool) *mux.Router {
	r := mux.NewRouter()

	r.Use(mw.PrometheusMiddleware(info.Name))
	r.Use(mw.LoggingMiddleware)

	r.HandleFunc("/version", versionHandler(info)).
		Methods("GET")
	r.HandleFunc("/health", healthHandler).
		Methods("GET")
	r.HandleFunc("/readiness", readinessHandler(isReady, shutdown)).
		Methods("GET")
	r.Handle("/metrics", promhttp.Handler()).
		Methods("GET")
	return r
}

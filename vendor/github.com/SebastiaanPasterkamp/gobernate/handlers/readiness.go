package handlers

import (
	"net/http"
	"sync/atomic"
)

// readinessHandler is a readiness probe. Will only return http.StatusOK when
// the isReady value is true, and the shutdown channel is still open.
func readinessHandler(isReady *atomic.Value, shutdown chan bool) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		select {
		case <-shutdown:
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		default:
			w.WriteHeader(http.StatusOK)
		}
	}
}

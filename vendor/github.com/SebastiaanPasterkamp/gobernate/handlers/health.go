package handlers

import (
	"net/http"
)

// health is a liveness probe.
func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

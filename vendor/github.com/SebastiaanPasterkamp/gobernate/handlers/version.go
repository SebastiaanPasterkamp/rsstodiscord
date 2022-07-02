package handlers

import (
	"github.com/SebastiaanPasterkamp/gobernate/version"

	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// versionHandler returns a simple HTTP handler function which writes a json
// response detailing the version of the service.
func versionHandler(info version.Info) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		body, err := json.Marshal(info)
		if err != nil {
			log.Printf("Could not encode info data: %v", err)
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}

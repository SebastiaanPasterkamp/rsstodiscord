package middleware

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// LoggingMiddleware adds a Debug log line for every incoing request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Debug(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

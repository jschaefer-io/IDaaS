package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Composite of a response writer
// fo track the response status code
type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

// Allows the WriteHeader method to set
// the status code
func (w *statusResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Basic logging middleware
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		writer := statusResponseWriter{w, 200}
		next.ServeHTTP(&writer, r)

		duration := fmt.Sprintf("%dms", time.Now().Sub(start).Milliseconds())
		log.Printf("%s -> %d: (%s)\t %s\n", r.Method, writer.status, duration, r.URL)
	})
}

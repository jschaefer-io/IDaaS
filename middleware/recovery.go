package middleware

import (
	"fmt"
	"log"
	"net/http"
)

// Middleware to catch runtime panics
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				// Log recovery
				log.Printf("Recovered %s @ %s -> %v\n", r.Method, r.URL, err)

				// Write HTTP Error
				http.Error(w, fmt.Sprintf("%v", err), 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"time"
)

func (s *Server) init() {
	r := chi.NewRouter()

	// Default middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprint(writer, "hello")
	})

	s.router = r
}

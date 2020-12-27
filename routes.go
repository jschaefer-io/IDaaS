package main

import (
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/jschaefer-io/IDaaS/middleware"
	"time"
)

func (s *Server) init() {
	r := chi.NewRouter()

	// Default middleware
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(30 * time.Second))

	// API Routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.TypeJson())

		r.Get("/users", s.UsersGetAll())
		r.Post("/users", s.UsersCreate())
		r.Get("/users/{userId}", s.UsersGetSingle())
		r.Patch("/users/{userId}", s.UserUpdateSingle())
	})

	r.Get("/", s.AuthRequest())

	s.router = r
}

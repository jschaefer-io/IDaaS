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
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.ContentTypeJson)

		// jwt authentication endpoints
		r.Post("/login", s.UserAuthenticate())
		// @todo token refresh

		// management endpoints
		r.Group(func(r chi.Router) {
			r.Use(middleware.TokenAuthenticated(s.db))

			// user management
			r.Get("/users", s.UsersGetAll())
			r.Post("/users", s.UsersCreate())
			r.Get("/users/{userId}", s.UsersGetSingle())
			r.Patch("/users/{userId}", s.UserUpdateSingle())
			r.Delete("/users/{userId}", s.UserDeleteSingle())
		})
	})

	s.router = r
}

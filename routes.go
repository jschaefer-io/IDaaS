package main

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jschaefer-io/IDaaS/handler"
	"github.com/jschaefer-io/IDaaS/server"
)

func ServerRoutes(s *server.Server) {
	// TopLevelMiddleware
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(handler.Logger(s.Logger))
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Timeout(30 * time.Second))

	// Web Routes
	s.Router.Group(func(r chi.Router) {
		controller := handler.NewWebController(s.Components, s.Settings)

		// Web Middleware
		r.Use(handler.Html)

		// Login
		r.Get("/login", controller.LoginForm())
		r.Post("/login", controller.HandleLoginForm())

		// PW Reset
		r.Get("/reset", controller.PasswordResetStart())
		r.Post("/reset", controller.HandlePasswordResetStart())
		r.Get("/reset/confirm", controller.PasswordResetForm())
		r.Post("/reset/confirm", controller.HandlePasswordResetForm())

		// Registration
		r.Get("/user/confirm", controller.HandleUserConfirm())
	})

	// API Routes

	s.Router.Route("/api", func(r chi.Router) {
		controller := handler.NewApiController(s.Components, s.Settings)

		// API Middleware
		r.Use(handler.Json)

		// Create User Endpoint
		r.Post("/user", controller.AddUser())

		// Current User Endpoints
		r.Route("/me", func(r chi.Router) {
			r.Use(handler.TokenAuth(s.TokenManager, s.Components.Repositories.UserRepository))
			r.Get("/", controller.GetUser())
			r.Put("/", controller.UpdateUser()) // @todo add missing fields
		})

		// Token Endpoints
		r.Get("/token", controller.TokenCheck())
		// @todo add second token endpoint which checks sessions and grants a new token pair
		r.Post("/token/grant", controller.GrantTokens())
		r.Post("/token/refresh", controller.TokenRefresh())

		// Errors
		r.NotFound(controller.ErrorNotFound())
		r.MethodNotAllowed(controller.ErrorMethodNotAllowed())
	})

	// @todo add 404 handler
}

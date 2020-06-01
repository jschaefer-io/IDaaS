package main

import (
	"github.com/gorilla/mux"
	"github.com/jschaefer-io/IDaaS/action"
	"github.com/jschaefer-io/IDaaS/db"
	"github.com/jschaefer-io/IDaaS/middleware"
	"github.com/jschaefer-io/IDaaS/model"
	"github.com/jschaefer-io/IDaaS/resource"
	"net/http"
)

type Test struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

func main() {

	// Model Migrations
	db.Get().AutoMigrate(model.Identity{})

	r := mux.NewRouter()

	// Middleware
	r.Use(middleware.ContentJson)
	r.Use(middleware.Recovery)
	r.Use(middleware.Logger)

	// Add Resource Routes
	resource.NewResource("/identities", new(action.Identity)).Apply(r)

	// Plain Routes
	r.HandleFunc("/auth/login", action.AuthLogin).Methods("POST")
	r.Handle("/me", middleware.Auth(http.HandlerFunc(action.AuthMe))).Methods("GET")

	// Error Routes
	r.NotFoundHandler = middleware.ContentJson(http.HandlerFunc(action.Error404))
	r.MethodNotAllowedHandler = middleware.ContentJson(http.HandlerFunc(action.Error405))

	// Start Webservice
	err := http.ListenAndServe(":8080", r)

	if err != nil {
		panic(err)
	}
}

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jschaefer-io/IDaaS/action"
	"github.com/jschaefer-io/IDaaS/db"
	"github.com/jschaefer-io/IDaaS/middleware"
	"github.com/jschaefer-io/IDaaS/model"
	"github.com/jschaefer-io/IDaaS/resource"
)

func main() {

	// Model Migrations
	db.Get().AutoMigrate(model.Identity{})

	r := gin.Default()

	// Middlewares
	r.Use(middleware.Noop())

	// Add Resource Routes
	resource.NewResource("/test", new(action.BaseActionSet)).Apply(r)
	resource.NewResource("/identities", new(action.Identity)).Apply(r)

	// Plain Routes
	r.POST("/auth/login", action.AuthLogin)
	r.GET("/me", middleware.Auth(), action.AuthMe)

	// Error Routes
	r.NoRoute(action.Error404)
	r.NoMethod(action.Error405)

	// Start webserver
	r.Run()
}

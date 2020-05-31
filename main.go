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

	r := gin.New()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Add Resource Routes
	resource.NewResource("/identities", new(action.Identity)).Apply(r)

	// Plain Routes
	r.POST("/auth/login", action.AuthLogin)
	r.GET("/me", middleware.Auth(), action.AuthMe)

	// Error Routes
	r.NoRoute(action.Error404)
	r.NoMethod(action.Error405)

	// Start Webservice
	err := r.Run()
	if err != nil {
		panic(err)
	}
}

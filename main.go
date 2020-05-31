package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jschaefer-io/IDaaS/action"
	"github.com/jschaefer-io/IDaaS/resource"
)

func main() {

	// https://github.com/dgrijalva/jwt-go
	// https://golang.org/pkg/net/smtp/

	r := gin.Default()

	// Add Resource Routes
	resource.NewResource("/test", new(action.BaseActionSet)).Apply(r)
	resource.NewResource("/identities", new(action.Identity)).Apply(r)

	// Error Routes
	r.NoRoute(action.Error404)
	r.NoMethod(action.Error405)

	// Start webserver
	r.Run()
}

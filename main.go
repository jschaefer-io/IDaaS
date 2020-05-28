package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jschaefer-io/IDaaS/action"
)

func main() {

	// https://github.com/dgrijalva/jwt-go
	// https://golang.org/pkg/net/smtp/



	r := gin.Default()
	r.GET("/", action.IdentityIndex)
	r.POST("/", action.IdentityCreate)
	r.GET("/:id", action.IdentityShow)
	r.PUT("/:id", action.IdentityUpdate)
	r.DELETE("/:id", action.IdentityDelete)
	r.Run()
}
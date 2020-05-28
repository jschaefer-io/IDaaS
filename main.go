package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jschaefer-io/IDaaS/action"
)

func main() {

	// https://github.com/dgrijalva/jwt-go
	// https://golang.org/pkg/net/smtp/

	r := gin.Default()
	r.GET("/", action.GetIdentity)
	r.POST("/", action.AddIdentity)
	r.Run()
}
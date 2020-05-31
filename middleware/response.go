package middleware

import (
	"github.com/gin-gonic/gin"
)

// Test Middleware
func Noop() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
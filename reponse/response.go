package reponse

import "github.com/gin-gonic/gin"

// Default JSON Response Interface
type Response interface {
	Apply(c *gin.Context)
}

// Base Response Struct
type Base struct {
	code int
	data interface{}
}

// Applies the response to the given gin context
func (r Base) Apply(c *gin.Context) {
	c.JSON(r.code, r.data)
}

// Creates a new base response object
func NewResponse(code int, data interface{}) Base {
	return Base{code, data}
}

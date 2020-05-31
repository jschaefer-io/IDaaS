package reponse

import "github.com/gin-gonic/gin"

type Response interface {
	Apply(c *gin.Context)
}

type Base struct {
	code int
	data interface{}
}

func (r Base) Apply(c *gin.Context) {
	c.JSON(r.code, r.data)
}

func NewResponse(code int, data interface{}) Base {
	return Base{code, data}
}

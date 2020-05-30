package action

import (
	"github.com/gin-gonic/gin"
)

type Set interface {
	Index(*gin.Context)
	Create(*gin.Context)
	Show(*gin.Context, int)
	Delete(*gin.Context, int)
	Update(*gin.Context, int)
}

type BaseActionSet struct{}

func (b BaseActionSet) Index(c *gin.Context) {
	Error404(c)
}

func (b BaseActionSet) Create(c *gin.Context) {
	Error404(c)
}

func (b BaseActionSet) Show(c *gin.Context, id int) {
	Error404(c)
}

func (b BaseActionSet) Delete(c *gin.Context, id int) {
	Error404(c)
}

func (b BaseActionSet) Update(c *gin.Context, id int) {
	Error404(c)
}
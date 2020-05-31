package action

import (
	"github.com/gin-gonic/gin"
)

// Basic Resource Action Set
type Set interface {
	Index(*gin.Context)
	Create(*gin.Context)
	Show(*gin.Context, int)
	Delete(*gin.Context, int)
	Update(*gin.Context, int)
}

// Base Action set which can be
// embedded to fulfill the Set
// interface with only partial
// action support
type BaseActionSet struct{}

// Default index Route
// reults in a 404 error
func (b BaseActionSet) Index(c *gin.Context) {
	Error404(c)
}

// Default index Route
// reults in a 404 error
func (b BaseActionSet) Create(c *gin.Context) {
	Error404(c)
}

// Default show Route
// reults in a 404 error
func (b BaseActionSet) Show(c *gin.Context, id int) {
	Error404(c)
}

// Default delete Route
// reults in a 404 error
func (b BaseActionSet) Delete(c *gin.Context, id int) {
	Error404(c)
}

// Default update Route
// reults in a 404 error
func (b BaseActionSet) Update(c *gin.Context, id int) {
	Error404(c)
}

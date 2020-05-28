package resource

import (
	"github.com/gin-gonic/gin"
	"github.com/jschaefer-io/IDaaS/reponse"
	"net/http"
)

type Resource interface {
	Index(c *gin.Context)
	Create(c *gin.Context)
	Show(c *gin.Context, model interface{})
	Delete(c *gin.Context, model interface{})
	Update(c *gin.Context, model interface{})
}

type BaseResource struct{}

func (b BaseResource) Index(c *gin.Context) {
	reponse.NewError(http.StatusNotFound, "Resource not found").Apply(c)
}

func (b BaseResource) Create(c *gin.Context) {
	reponse.NewError(http.StatusNotFound, "Resource not found").Apply(c)
}

func (b BaseResource) Show(c *gin.Context, model interface{}) {
	reponse.NewError(http.StatusNotFound, "Resource not found").Apply(c)
}

func (b BaseResource) Delete(c *gin.Context, model interface{}) {
	reponse.NewError(http.StatusNotFound, "Resource not found").Apply(c)
}

func (b BaseResource) Update(c *gin.Context, model interface{}) {
	reponse.NewError(http.StatusNotFound, "Resource not found").Apply(c)
}

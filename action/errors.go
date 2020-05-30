package action

import (
	"github.com/gin-gonic/gin"
	"github.com/jschaefer-io/IDaaS/reponse"
	"net/http"
)

func Error404(c *gin.Context) {
	reponse.NewError(http.StatusNotFound, "Resource not found").Apply(c)
}

func Error405(c *gin.Context) {
	reponse.NewError(http.StatusMethodNotAllowed, "Method not allowed on this resource").Apply(c)
}

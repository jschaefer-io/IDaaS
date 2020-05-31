package action

import (
	"github.com/gin-gonic/gin"
	"github.com/jschaefer-io/IDaaS/reponse"
	"net/http"
)

// Returns the base 404 Error Response
func Error404(c *gin.Context) {
	reponse.NewError(http.StatusNotFound, "Resource not found").Apply(c)
}

// Returns the base 405 Error Response
func Error405(c *gin.Context) {
	reponse.NewError(http.StatusMethodNotAllowed, "Method not allowed on this resource").Apply(c)
}

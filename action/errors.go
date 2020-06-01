package action

import (
	"github.com/jschaefer-io/IDaaS/reponse"
	"net/http"
)

// Returns the base 404 Error Response
func Error404(w http.ResponseWriter, r *http.Request) {
	reponse.NewError(http.StatusNotFound, "Resource not found").Apply(w)
}

// Returns the base 405 Error Response
func Error405(w http.ResponseWriter, r *http.Request) {
	reponse.NewError(http.StatusMethodNotAllowed, "Method not allowed on this resource").Apply(w)
}

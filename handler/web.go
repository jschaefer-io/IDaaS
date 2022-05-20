package handler

import (
	"errors"
	"net/http"

	"github.com/jschaefer-io/IDaaS/server"
)

type WebController struct {
	baseController
}

func NewWebController(components *server.Components, settings *server.Settings) *WebController {
	return &WebController{
		baseController: newBaseController(components, settings),
	}
}

func (c *WebController) Error(writer http.ResponseWriter, statusCode int) {
	writer.WriteHeader(statusCode)
	_, _ = c.components.Templates.ExecuteToString("error", "server", nil)
}

func (c *WebController) getRedirectFromHeader(request *http.Request) (string, error) {
	redirect := request.URL.Query()["redirect"]
	if len(redirect) == 0 {
		return "", errors.New("no redirect parameter")
	}
	return redirect[0], nil
}

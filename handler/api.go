package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jschaefer-io/IDaaS/server"
)

type ApiController struct {
	baseController
}

func NewApiController(components *server.Components, settings *server.Settings) *ApiController {
	return &ApiController{
		baseController: newBaseController(components, settings),
	}
}

type errorResponse struct {
	Code  int `json:"code"`
	Error any `json:"error"`
}

func ApiError(writer http.ResponseWriter, statusCode int, data any) {
	writer.WriteHeader(statusCode)
	response := errorResponse{
		Code:  statusCode,
		Error: data,
	}
	jsonResponse, _ := json.Marshal(response)
	_, _ = writer.Write(jsonResponse)
}

func (c *ApiController) ErrorNotFound() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ApiError(writer, http.StatusNotFound, "not found")
	}
}

func (c *ApiController) ErrorMethodNotAllowed() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ApiError(writer, http.StatusMethodNotAllowed, "method not allowed")
	}
}

package reponse

import (
	"encoding/json"
	"net/http"
)

// Default JSON Response Interface
type Response interface {
	Apply(w http.ResponseWriter)
}

// Base Response Struct
type Base struct {
	code int
	data interface{}
}

// Applies the response to the given http response writer
func (r Base) Apply(w http.ResponseWriter) {
	w.WriteHeader(r.code)
	data, err := json.Marshal(r.data)
	if err == nil {
		_, _ = w.Write(data)
	}
}

// Creates a new base response object
func NewResponse(code int, data interface{}) Base {
	return Base{code, data}
}

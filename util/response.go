package util

import (
	"encoding/json"
	"net/http"
)

func MarshalToResponseWriter(data interface{}, writer http.ResponseWriter, statusCode int) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		jsonData, _ = json.Marshal(NewErrorObject(err.Error()))
	}
	writer.WriteHeader(statusCode)
	_, _ = writer.Write(jsonData)
}

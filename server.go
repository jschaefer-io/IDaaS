package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type Server struct {
	db     *gorm.DB
	router *chi.Mux
}

func NewServer(db *gorm.DB) Server {
	srv := Server{
		db: db,
	}
	srv.init()
	return srv
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(writer, r)
}

func (s *Server) newErrorObject(message string) map[string]interface{} {
	msg := make(map[string]interface{})
	msg["error"] = true
	msg["message"] = message
	return msg
}

// Generates a validation Error response
func (s *Server) validationError(err error) map[string]interface{} {
	validationError, ok := err.(validator.ValidationErrors)
	if !ok {
		return s.newErrorObject("invalid post data provided")
	}
	messages := make(map[string]interface{})
	messages["error"] = true
	errorList := make(map[string]string)
	for _, fieldErr := range validationError {
		errorList[strings.ToLower(fieldErr.Field())] = fieldError{fieldErr}.String()
	}

	messages["messages"] = errorList
	return messages
}

func (s *Server) marshalToResponseWriter(data interface{}, writer http.ResponseWriter, statusCode int) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		jsonData, _ = json.Marshal(s.newErrorObject(err.Error()))
	}
	writer.WriteHeader(statusCode)
	_, _ = writer.Write(jsonData)
}

type fieldError struct {
	err validator.FieldError
}

func (q fieldError) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Validation failed on field '%s'", strings.ToLower(q.err.Field())))

	// Print failed condition
	sb.WriteString(fmt.Sprintf(", condition: %s", q.err.ActualTag()))

	// Handle Parameter Validation
	if q.err.Param() != "" {
		sb.WriteString(fmt.Sprintf(" { %s }", q.err.Param()))
	}

	// Print actual value
	if q.err.Value() != nil && q.err.Value() != "" {
		sb.WriteString(fmt.Sprintf(", given value: %v", q.err.Value()))
	}
	return sb.String()
}

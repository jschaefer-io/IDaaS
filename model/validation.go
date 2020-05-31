package model

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jschaefer-io/IDaaS/reponse"
	"net/http"
	"strings"
)

// Generates a validation Error response
func ValidationError(err error) reponse.Response {
	code := http.StatusUnprocessableEntity
	validationError, ok := err.(validator.ValidationErrors)
	if !ok {
		return reponse.NewError(code, err.Error())
	}
	var messages []string
	for _, fieldErr := range validationError {
		messages = append(messages, fieldError{fieldErr}.String())
	}
	return reponse.NewError(code, messages)
}

// Base wrapper for the
// field validation
type fieldError struct {
	err validator.FieldError
}

// Implement stringer interface
// to cleanup the fieldError.err
func (q fieldError) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Validation failed on field '%s'", q.err.Field()))

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

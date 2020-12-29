package util

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func NewErrorObject(message string) map[string]interface{} {
	msg := make(map[string]interface{})
	msg["error"] = true
	msg["message"] = message
	return msg
}

// Generates a validation Error response
func ValidationError(err error) map[string]interface{} {
	validationError, ok := err.(validator.ValidationErrors)
	if !ok {
		return NewErrorObject("invalid post data provided")
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

package utils

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"strings"
)

type Rule func(string) error

type Validator struct {
	ErrorBag *ErrorBag
}

func NewValidator() *Validator {
	return &Validator{
		ErrorBag: NewErrorBag(),
	}
}

func (v *Validator) Validate(field string, value string, rules ...Rule) {
	for _, rule := range rules {
		err := rule(value)
		if err != nil {
			v.ErrorBag.AddError(field, err.Error())
			return
		}
	}
}

func RuleRequired() Rule {
	return func(value string) error {
		if len(value) == 0 {
			return errors.New("field is required")
		}
		return nil
	}
}

func RuleEmail() Rule {
	return func(value string) error {
		if _, ok := mail.ParseAddress(value); ok != nil {
			return errors.New("field must be valid email")
		}
		return nil
	}
}

func RulePassword() Rule {
	return func(value string) error {
		if len(value) < 6 {
			return errors.New("password must be at least 6 characters")
		}
		return nil
	}
}

func RuleEqual(compareValue string) Rule {
	return func(value string) error {
		if value != compareValue {
			return errors.New("fields must be equal")
		}
		return nil
	}
}

func RuleIsIn(values ...fmt.Stringer) Rule {
	return func(value string) error {
		for _, v := range values {
			if v.String() == value {
				return nil
			}
		}
		list := make([]string, len(values))
		for i, v := range values {
			list[i] = v.String()
		}
		return errors.New("field must be one of the following values: " + strings.Join(list, ", "))
	}
}

func RuleUrl() Rule {
	return func(value string) error {
		if _, err := url.ParseRequestURI(value); err != nil {
			return errors.New("field must be valid url")
		}
		return nil
	}
}

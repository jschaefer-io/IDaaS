package resource

import (
	"errors"
	"regexp"
)

// Returns an url parameter matching the given regex string
func GetParam(value string, regex string) (string, error) {
	if !regexp.MustCompile(regex).Match([]byte(value)) {
		return value, errors.New("parameter does not match")
	}
	return value, nil
}
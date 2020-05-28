package crypto

import (
	"fmt"
	"math/rand"
)

// Basic Token Type
type Token string

// To String type conversion
func (t *Token) String() string {
	return string(*t)
}

// Generates a new 128 char long random string
// which should be reasonably unique
func NewToken() Token {
	b := make([]byte, 64)
	rand.Read(b)
	return Token(fmt.Sprintf("%x", b))
}

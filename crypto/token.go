package crypto

import (
	"crypto/rand"
	"fmt"
)

type Token string

func (t *Token) String() string {
	return string(*t)
}

func NewToken() Token {
	b := make([]byte, 64)
	_, _ = rand.Read(b)
	return Token(fmt.Sprintf("%x", b))
}

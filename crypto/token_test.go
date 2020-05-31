package crypto

import (
	"testing"
)

// Tests if a new token is generated
// and is of proper length
func TestNewToken(t *testing.T) {
	token := NewToken()
	if len(token) != 128 {
		t.Errorf("unexpected token length. expected %d but got %d", 128, len(token))
	}
}

// Tests if subsequent calls to the
// new token method yield unique tokens
func TestTokenUnique(t *testing.T) {
	a := NewToken()
	b := NewToken()
	if a == b {
		t.Error("tokens shouldn't be unique")
	}
}

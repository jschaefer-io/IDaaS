package model

import (
	"github.com/jschaefer-io/IDaaS/crypto"
)

// Basic Identity instance
// which the full IDaaS is
// based on
type Identity struct {
	Model
	Email    string          `json:"email"`
	Password crypto.Password `json:"password"`
	Token    crypto.Token    `json:"-"`
}

// Creates and prepares the new Identity Instance
func NewIdentity(email string, password string) Identity {
	return Identity{
		Email:    email,
		Password: crypto.NewPassword(password),
		Token:    crypto.NewToken(),
	}
}

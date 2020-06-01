package model

import (
	"github.com/jschaefer-io/IDaaS/crypto"
	"github.com/jschaefer-io/IDaaS/db"
)

// Basic Identity instance
// which the full IDaaS is
// based on
type Identity struct {
	Model
	PasswordForm
	Email string       `json:"email" validate:"required,email"`
	Token crypto.Token `json:"-"`
}

// Finds an identity by its id from the database
func (i Identity) Find(id int) (Identity, error) {
	if err := db.Get().Find(&i, id).Error; err != nil {
		return i, err
	}
	return i, nil
}

// Creates and prepares the new Identity Instance
func NewIdentity(email string, password string) Identity {
	pwd := crypto.NewPassword(password)
	return Identity{
		Email: email,
		Token: crypto.NewToken(),
		PasswordForm: PasswordForm{
			Password: pwd.String(),
		},
	}
}

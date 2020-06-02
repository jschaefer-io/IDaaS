package model

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jschaefer-io/IDaaS/crypto"
	"github.com/jschaefer-io/IDaaS/db"
)

// Basic Identity instance
// which the full IDaaS is
// based on
type Identity struct {
	Model
	PasswordForm
	Email string       `json:"email" validate:"required,email,dbunique=identities"`
	Token crypto.Token `json:"-"`
}

type IdentityLogin struct {
	PasswordForm
	Email string `json:"email" validate:"required,email"`
}

// Finds an identity by its id from the database
func (i Identity) Find(id int) (Identity, error) {
	if err := db.Get().Find(&i, id).Error; err != nil {
		return i, err
	}
	return i, nil
}

// Fetches the user associated with the given
// jwt and returns the jwt validation
func (i Identity) CheckJwt(token string) (Identity, jwt.MapClaims, *jwt.Token, error) {
	claims, t, err := crypto.CheckJWT(token, "user-id", func(id int) (crypto.Token, error) {
		var err error
		i, err = i.Find(id)
		if err != nil {
			return "", err
		}
		return i.Token, nil
	})
	return i, claims, t, err
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

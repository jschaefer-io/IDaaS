package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// Basic Password type
type Password string

// Hashes the current password value
func (p *Password) rehash() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*p), bcrypt.MinCost)
	if err != nil {
		return err
	}
	newHash := Password(bytes)
	*p = newHash
	return nil
}

// Updates and rehashes the password value
func (p *Password) Update(pwd string) {
	*p = Password(pwd)
	_ = p.rehash()
}

// Checks if the given password string
// is matches the current password hash
func (p *Password) Compare(pwd string) bool {
	check := []byte(*p)
	return bcrypt.CompareHashAndPassword(check, []byte(pwd)) == nil
}

// To String conversion
func (p *Password) String() string {
	return string(*p)
}

// Generates a hashed Password instance
func NewPassword(pwd string) Password {
	p := Password(pwd)
	_ = p.rehash()
	return p
}

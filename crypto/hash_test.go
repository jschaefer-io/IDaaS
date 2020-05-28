package crypto

import (
	"testing"
)

// Tests if password can easily
// be casted to type string and back
func TestPassword(t *testing.T) {
	pw := Password("test")
	str := string(pw)
	if str != pw.String() || Password(str) != pw {
		t.Error("value of password changed during type conversion")
	}
}

// Test if a password gets hashed correctly
func TestNewPassword(t *testing.T) {
	pwd := "test"
	pw := NewPassword(pwd)
	if pw.String() == pwd {
		t.Error("password has not been hashed")
	}
	if len(pw) != 60 {
		t.Errorf("unexpected hash length, expected %d but got %d", 60, len(pw))
	}
}

// Tests, if we can compare the password hash
func TestCompare(t *testing.T) {
	pwd := "test"
	p := NewPassword(pwd)
	if !p.Compare(pwd) {
		t.Error("error comparing matching passwords")
	}
	if p.Compare(pwd + "a") {
		t.Error("error comparing passwords which do not match")
	}
}

// Checks if a password can be updated properly
func TestPassword_Update(t *testing.T) {
	pwd := "test"
	p := NewPassword(pwd)
	p.Update(pwd + "a")
	if !p.Compare(pwd + "a") {
		t.Error("error updating a password")
	}
}

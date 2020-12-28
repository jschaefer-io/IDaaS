package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	BaseModel
	Name     string `gorm:"size:255" json:"name"`
	Email    string `gorm:"size:255;unique" json:"email"`
	Password string `gorm:"size:255" json:"-"`
}

func (u *User) SetPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u.Password = string(hash)
	return u.Password
}
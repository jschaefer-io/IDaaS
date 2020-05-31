package model

import (
	"time"
)

// Base Model, which is contained in all
// other models
type Model struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


// Basic Password Form
type PasswordForm struct {
	Password string `json:"password" binding:"required"`
}
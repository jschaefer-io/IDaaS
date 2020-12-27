package models

type User struct {
	BaseModel
	Name     string `gorm:"size:255" json:"name"`
	Email    string `gorm:"size:255;unique" json:"email"`
	Password string `gorm:"size:255" json:"password"`
}

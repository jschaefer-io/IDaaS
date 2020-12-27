package models

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
	"strings"
	"time"
)

type BaseModel struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func BindJson(form interface{}, reader io.Reader) error {
	var err error

	// Get json from reader
	jsonString := new(strings.Builder)

	_, err = io.Copy(jsonString, reader)
	if err != nil {
		return err
	}

	// bind json to form struct
	err = json.Unmarshal([]byte(jsonString.String()), form)
	if err != nil {
		return err
	}

	// validates the form struct
	return validator.New().Struct(form)
}

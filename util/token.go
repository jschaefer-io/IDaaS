package util

import (
	"errors"
	"github.com/jschaefer-io/IDaaS/model"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func ExtractJWT(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	split := strings.Split(token, " ")

	// Check if token is present and properly formed
	if len(split) != 2 || split[0] != "Bearer" {
		return "", errors.New("invalid bearer token")
	}
	return split[1], nil
}

func ExtractSession(r *http.Request, db *gorm.DB) (*model.Session, error) {
	cookie, err := r.Cookie("id-session")
	if err != nil {
		return nil, err
	}

	session := new(model.Session)
	res := db.Where("token = ?", cookie.Value).Preload("User").First(session)
	if res.Error != nil {
		return nil, res.Error
	}

	return session, nil
}

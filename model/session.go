package model

import (
	"github.com/jschaefer-io/IDaaS/crypto"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Session struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint
	Token     crypto.Token `gorm:"size:128;unique"`
	ExpiresAt time.Time
	User      User
}

func (s *Session) Cookie() *http.Cookie {
	return &http.Cookie{
		Name:     "id-session",
		Value:    s.Token.String(),
		Path:     "/",
		MaxAge:   int(s.ExpiresAt.Sub(time.Now()).Seconds()),
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func NewSession(user User, db *gorm.DB) (*Session, error) {
	session := &Session{
		Token:     crypto.NewToken(),
		User:      user,
		ExpiresAt: time.Now().Add(time.Minute * 5),
	}

	res := db.Save(session)
	if res.Error != nil {
		return nil, res.Error
	}
	return session, nil
}

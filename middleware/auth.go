package middleware

import (
	"context"
	"crypto/rsa"
	"errors"
	"github.com/jschaefer-io/IDaaS/crypto"
	"github.com/jschaefer-io/IDaaS/models"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"sync"
)

func Authenticated(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		var init sync.Once
		var rsaSecret *rsa.PublicKey
		var rsaError error

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			init.Do(func() {
				rsaSecret, rsaError = crypto.ReadPublicRsaKey()
			})
			if rsaError != nil {
				w.WriteHeader(http.StatusInternalServerError)
				panic(rsaError)
			}

			tokenString, _ := extractJWT(r)
			id, err := crypto.ParseJwt(tokenString, rsaSecret)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			usr := new(models.User)
			res := db.Find(usr, id)
			if res.Error != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user", usr)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractJWT(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	split := strings.Split(token, " ")

	// Check if token is present and properly formed
	if len(split) != 2 || split[0] != "Bearer" {
		return "", errors.New("invalid bearer token")
	}
	return split[1], nil
}

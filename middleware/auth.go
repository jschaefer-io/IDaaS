package middleware

import (
	"context"
	"crypto/rsa"
	"github.com/jschaefer-io/IDaaS/crypto"
	"github.com/jschaefer-io/IDaaS/model"
	"github.com/jschaefer-io/IDaaS/util"
	"gorm.io/gorm"
	"net/http"
	"sync"
)

func TokenAuthenticated(db *gorm.DB) func(next http.Handler) http.Handler {
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

			tokenString, _ := util.ExtractJWT(r)
			id, err := crypto.ParseJwt(tokenString, rsaSecret)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			usr := new(model.User)
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

func SessionAuthenticated(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			session, err := util.ExtractSession(r, db)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				// @todo redirect to login page
				return
			}

			ctx := context.WithValue(r.Context(), "user", session.User)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

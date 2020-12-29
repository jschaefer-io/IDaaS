package middleware

import (
	"context"
	"github.com/jschaefer-io/IDaaS/model"
	"gorm.io/gorm"
	"net/http"
)

func SessionAuthenticated(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			cookie, err := r.Cookie("id-session")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				// @todo redirect to login page
				return
			}

			session := new(model.Session)
			res := db.Where("token = ?", cookie.Value).Preload("User").First(session)
			if res.Error != nil {
				w.WriteHeader(http.StatusUnauthorized)
				// @todo redirect to login page
				return
			}

			ctx := context.WithValue(r.Context(), "user", session.User)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

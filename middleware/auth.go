package middleware

import (
	"github.com/gorilla/context"
	"github.com/jschaefer-io/IDaaS/crypto"
	"github.com/jschaefer-io/IDaaS/model"
	"github.com/jschaefer-io/IDaaS/reponse"
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		message := "Permission denied"
		token, err := crypto.ExtractJWT(r)
		if err != nil {
			reponse.NewError(http.StatusUnauthorized, message).Apply(w)
			return
		}

		// Check if the jwt is valid
		identity, _, t, err := model.Identity{}.CheckJwt(token)

		if err != nil || !t.Valid {
			reponse.NewError(http.StatusUnauthorized, message).Apply(w)
			return
		}

		// Attach identity to the context
		context.Set(r, "identity", identity)

		// If Token is valid, pass to next middleware/ the action
		next.ServeHTTP(w, r)
	})
}

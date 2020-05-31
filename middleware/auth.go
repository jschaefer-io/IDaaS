package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jschaefer-io/IDaaS/crypto"
	"github.com/jschaefer-io/IDaaS/model"
	"github.com/jschaefer-io/IDaaS/reponse"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		message := "Permission denied"
		token, err := crypto.ExtractJWT(c)
		if err != nil {
			reponse.NewError(http.StatusUnauthorized, message).Apply(c)
			c.Abort()
			return
		}

		// Check if the jwt is valid
		var identity model.Identity
		t, err := crypto.CheckJWT(token, "user-id", func(id int) (crypto.Token, error) {
			var err error
			identity, err = model.Identity{}.Find(id)
			if err != nil {
				return "", err
			}
			return identity.Token, nil
		})
		if err != nil || !t.Valid {
			reponse.NewError(http.StatusUnauthorized, message).Apply(c)
			c.Abort()
			return
		}

		// If the token is valid
		// set the logged in identity
		c.Set("identity", identity)

		// If Token is valid, pass to next middleware/ the action
		c.Next()
	}
}

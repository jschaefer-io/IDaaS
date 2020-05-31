package crypto

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// Generates a new jwt with the given secret token
// mas the secret key
func NewJWT(secret Token, data map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * 10).Unix(),
	}

	// Merge custom data in the claims map
	for k, v := range data {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Checks if the given token is valid
func CheckJWT(token string, field string, check func(id int) (Token, error)) (*jwt.Token, error) {
	claims := jwt.MapClaims{}
	return jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (i interface{}, err error) {
		id, ok := claims[field].(float64)
		if !ok {
			return nil, errors.New("required field in claims not set")
		}
		secret, err := check(int(id))
		if err != nil {
			return nil, err
		}
		return []byte(secret), nil
	})
}

// Extracts the JWT from the gin.Context header
// as an Authorization Bearer-Token
func ExtractJWT(c *gin.Context) (string, error) {
	token := c.GetHeader("Authorization")
	split := strings.Split(token, " ")

	// Check if token is present and properly formed
	if len(split) != 2 || split[0] != "Bearer" {
		return "", errors.New("invalid bearer token")
	}
	return split[1], errors.New("")
}

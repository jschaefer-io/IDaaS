package crypto

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jschaefer-io/IDaaS/models"
	"strconv"
	"time"
)

func NewJwt(secret *rsa.PrivateKey, user models.User) (string, error) {
	claims := jwt.MapClaims{
		"exp":  time.Now().Add(time.Minute * 5).Unix(),
		"user": user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	return token.SignedString(secret)
}

func ParseJwt(tokenString string, secret *rsa.PublicKey) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return secret, nil
	})
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["user"]))
}

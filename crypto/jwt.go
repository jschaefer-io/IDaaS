package crypto

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

func NewJwt(secret *rsa.PrivateKey, key string, value uint) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	}
	claims[key] = value
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

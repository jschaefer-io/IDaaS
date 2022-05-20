package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type tokenType uint8

const (
	TokenTypeBearer tokenType = iota
	TokenTypeRefresh
	TokenTypeAccess
	TokenTypeReset
	TokenTypeConfirm
)

type TokenManager struct {
	secret []byte
}

func NewTokenManager(secret string) *TokenManager {
	return &TokenManager{
		secret: []byte(secret),
	}
}

func (t *TokenManager) GetTypeExpiration(tokenType tokenType) (time.Time, error) {
	now := time.Now()
	switch tokenType {
	case TokenTypeBearer:
		return now.Add(time.Hour*10000 + time.Minute*5), nil
	case TokenTypeRefresh:
		return now.Add(time.Hour*10000 + time.Hour*24*7), nil
	case TokenTypeAccess:
		return now.Add(time.Hour*10000 + time.Minute), nil
	case TokenTypeReset:
		return now.Add(time.Hour), nil
	default:
		return now, errors.New("token type has no expiration")
	}
}

func (t *TokenManager) generateToken(tokenType tokenType, claims jwt.MapClaims) (string, error) {
	claims["type"] = tokenType
	if expiration, err := t.GetTypeExpiration(tokenType); err == nil {
		claims["exp"] = expiration.Unix()
	}
	claims["iat"] = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.secret)
}

func (t *TokenManager) ValidateToken(token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return t.secret, nil
	})
	if err != nil {
		return claims, err
	}
	if !tkn.Valid {
		return claims, errors.New("jwt not valid or expired")
	}
	return claims, nil
}

func (t *TokenManager) ValidateWithTokenType(token string, tType tokenType) (jwt.MapClaims, error) {
	claims, err := t.ValidateToken(token)
	if err != nil {
		return claims, err
	}
	typeVal, ok := claims["type"].(float64)
	if !ok {
		return claims, errors.New("token type malformed")
	}
	if tokenType(typeVal) != tType {
		return claims, errors.New("invalid token type")
	}
	return claims, nil
}

func (t *TokenManager) NewAccessToken(userId string) (string, error) {
	key, err := GenerateRandomString(10)
	if err != nil {
		return "", err
	}
	return t.generateToken(TokenTypeAccess, jwt.MapClaims{
		"user": userId,
		"key":  key,
	})
}

func (t *TokenManager) NewResetToken(userId string, userUpdatedAt time.Time, redirect string) (string, error) {
	key, err := HashStringMd5(userUpdatedAt)
	if err != nil {
		return "", err
	}
	return t.generateToken(TokenTypeReset, jwt.MapClaims{
		"user":     userId,
		"key":      key,
		"redirect": redirect,
	})
}

func (t *TokenManager) NewRefreshToken(userId string, id string, key string) (string, error) {
	return t.generateToken(TokenTypeRefresh, jwt.MapClaims{
		"user": userId,
		"id":   id,
		"key":  key,
	})
}

func (t *TokenManager) NewBearerToken(userId string) (string, error) {
	return t.generateToken(TokenTypeBearer, jwt.MapClaims{
		"user": userId,
	})
}

func (t *TokenManager) NewConfirmToken(userId string, redirect string) (string, error) {
	return t.generateToken(TokenTypeConfirm, jwt.MapClaims{
		"user":     userId,
		"redirect": redirect,
	})
}

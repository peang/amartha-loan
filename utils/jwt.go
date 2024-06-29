package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/peang/amartha-loan-service/models"
)

var salt = "lpm7kxcwofm0u5za4xasao"

type TokenPayload struct {
	jwt.StandardClaims
	Payload Payload
}

type Payload struct {
	ID   uint            `json:"id"`
	Name string          `json:"name"`
	Role models.UserRole `json:"role"`
	Exp  int64           `json:"exp"`
}

func CreateJWTToken(user *models.User, rememberMe bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	exp := time.Hour * 1
	if rememberMe {
		exp = time.Hour * 1
	}

	claims["payload"] = Payload{
		ID:   user.ID,
		Name: user.Name,
		Role: user.Role,
		Exp:  time.Now().Add(exp).Unix(),
	}

	signingKey := []byte(salt)

	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (tokenInfo TokenPayload, err error) {
	payload := &TokenPayload{}
	token, err := jwt.ParseWithClaims(tokenString, payload, func(token *jwt.Token) (interface{}, error) {
		return []byte(salt), nil
	})
	if err != nil || !token.Valid {
		return *payload, err
	}

	claim := token.Claims.(*TokenPayload)
	if claim.Payload.Exp < time.Now().Unix() {
		return *payload, errors.New("TOKEN EXPIRED")
	}

	return *payload, nil
}

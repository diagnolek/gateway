package helpers

import (
	"github.com/golang-jwt/jwt"
)

func GenerateJWTToken(userID uint, secret string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ValidateJWTToken(tokenString string, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	return token, err
}

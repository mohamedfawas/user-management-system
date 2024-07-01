package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret_key string = os.Getenv("SECRET")
var JwtKey = []byte(secret_key)

func GenerateToken(userID uint, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["is_admin"] = isAdmin
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
}

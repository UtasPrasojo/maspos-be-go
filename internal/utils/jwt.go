package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// Gunakan secret key yang aman, idealnya dari .env
var secretKey = []byte("maspos-secret-key-12345")

func GenerateToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
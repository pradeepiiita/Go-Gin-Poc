package middleware

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateToken generates a new JWT token for a given user
func GenerateToken(username string) (string, error) {
	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 240).Unix(),
	})

	// Sign and return the token
	return token.SignedString(secretKey)
}

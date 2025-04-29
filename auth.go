package main

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"go-api/models"
	"log"
)

var jwtKey = []byte("your_secret_key") // Secret key for signing JWT tokens

// Struct for JWT claims
type Claims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

// GenerateJWT generates a new JWT token for the user
func GenerateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours

	claims := &Claims{
		ID:    user.ID,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "go-api",
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}

	return tokenString, nil
}

// ParseJWT parses and validates a JWT token
func ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

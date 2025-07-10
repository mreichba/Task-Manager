package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mreichba/task-manager-backend/config"
	"github.com/mreichba/task-manager-backend/models"
)

// GenerateJWT generates a JWT token for a given user ID
func GenerateJWT(userID int) (string, error) {
	// Create the token payload (claims)
	claims := &models.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // expires in 1 day
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create the token using the HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// jwtKey is the signing secret used to sign and verify tokens
	jwtKey := []byte(config.AppConfig.JWTSecret)

	// Sign the token using your secret key
	return token.SignedString(jwtKey)
}

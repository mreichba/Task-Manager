package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mreichba/task-manager-backend/config"
	"github.com/mreichba/task-manager-backend/models"
)

// contextKey is a custom type to avoid collisions in context
type contextKey string

const userIDKey contextKey = "userID"

// UserIDKey returns the key for looking up userID in request context
func UserIDKey() contextKey {
	return userIDKey
}

// JWTMiddleware ensures the request has a valid JWT before proceeding
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// Expect: Authorization: Bearer <token>
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &models.Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, jwtKeyLookup)
		if err != nil {
			log.Printf("JWT parsing error: %v\n", err)
		}
		log.Printf("Token valid: %v\n", token.Valid)
		log.Printf("Token claims: %+v\n", claims)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Inject userID into context
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func jwtKeyLookup(token *jwt.Token) (interface{}, error) {
	// Optional: check the signing algorithm for extra safety
	if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	// Return the secret used to sign the token
	secret := config.AppConfig.JWTSecret
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET not set in environment")
	}

	return []byte(secret), nil
}

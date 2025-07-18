package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mreichba/task-manager-backend/config"
	"github.com/mreichba/task-manager-backend/logger"
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

		logger.Log.Debugf("Retrieved authHeader: %v", authHeader)

		// Expect: Authorization: Bearer <token>
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Log.Error("Missing or invalid Auth header")
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		logger.Log.Debugf("AuthHeader converted: %v", tokenStr)

		claims := &models.Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, jwtKeyLookup)
		if err != nil || !token.Valid {
			logger.Log.WithError(err).
				WithField("token.Valid", token.Valid).
				Error("Token invalid or failed to parse")
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		logger.Log.Debug("Token validated successfully")
		logger.Log.WithField("claims", claims).Debug("Parsed token claims")

		// Inject userID into context
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func jwtKeyLookup(token *jwt.Token) (interface{}, error) {
	// Optional: check the signing algorithm for extra safety
	if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
		logger.Log.Error("Incorrect token signing method")
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return []byte(config.AppConfig.JWTSecret), nil
}

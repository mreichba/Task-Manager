package handlers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/mreichba/task-manager-backend/auth"
	"github.com/mreichba/task-manager-backend/db"
	"github.com/mreichba/task-manager-backend/logger"
	"github.com/mreichba/task-manager-backend/middleware"
	"github.com/mreichba/task-manager-backend/models"
	"github.com/sirupsen/logrus"
)

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	logger.Log.Debug("RegisterUserHandler called")

	// Decode JSON from request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	logger.Log.WithField("email", user.Email).Debug("User input decoded")

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)
	logger.Log.Debug("Password hashed successfully")

	// Insert into DB
	query := `INSERT INTO users (username, email, password) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at`
	err = db.DB.QueryRow(query, user.Username, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		logger.Log.WithError(err).WithField("email", user.Email).Error("Failed to insert new user")
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Don't return password in response
	user.Password = ""

	logger.Log.WithFields(logrus.Fields{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
	}).Info("New user registered successfully")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.User

	logger.Log.Debug("LoginHandler called")

	// Decode JSON request body into creds
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		logger.Log.WithError(err).Warn("Failed to decode login input")
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	logger.Log.WithField("email", creds.Email).Debug("Login input decoded")

	// Find user in database by email
	var dbUser models.User

	query := `SELECT id, username, email, password, created_at FROM users WHERE email = $1`
	err := db.DB.QueryRow(query, creds.Email).Scan(
		&dbUser.ID,
		&dbUser.Username,
		&dbUser.Email,
		&dbUser.Password,
		&dbUser.CreatedAt,
	)

	if err != nil {
		logger.Log.WithError(err).WithField("email", creds.Email).Warn("User not found or DB error")
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	logger.Log.WithField("user_id", dbUser.ID).Debug("User record retrieved from DB")

	// Compare hashed password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(creds.Password))
	if err != nil {
		logger.Log.WithError(err).WithField("user_id", dbUser.ID).Warn("Password mismatch")
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	logger.Log.Debug("Password matched successfully")

	// Generate JWT
	token, err := auth.GenerateJWT(dbUser.ID)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to generate JWT")
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	logger.Log.WithField("user_id", dbUser.ID).Info("Login successful, token generated")

	// Return 200 OK, Token, and basic user data (no password) if match
	dbUser.Password = ""

	response := models.LoginResponse{
		Token: token,
		User: models.UserResponse{
			ID:        dbUser.ID,
			Username:  dbUser.Username,
			Email:     dbUser.Email,
			CreatedAt: dbUser.CreatedAt,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value(middleware.UserIDKey())

	logger.Log.Debugf("Raw userID value from context: %v", userIDVal)

	userID, ok := userIDVal.(int)
	if !ok {
		logger.Log.Warn("User ID missing or invalid in context")
		http.Error(w, "User ID missing or invalid", http.StatusUnauthorized)
		return
	}

	logger.Log.WithField("userID", userID).Debug("Looking up user in database")

	// Query the user from the database
	var user models.User
	query := `
		SELECT id, username, email, created_at 
		FROM users 
		WHERE id = $1
		`
	err := db.DB.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		logger.Log.WithError(err).WithField("user_id", userID).Error("Failed to fetch current user from database")
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	// Log successful lookup
	logger.Log.WithFields(logrus.Fields{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
	}).Info("Fetched current user")

	// Respond with user data
	response := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	logger.Log.WithField("user", response).Debug("User data to be returned")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

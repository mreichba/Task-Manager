package handlers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/mreichba/task-manager-backend/auth"
	"github.com/mreichba/task-manager-backend/db"
	"github.com/mreichba/task-manager-backend/models"
)

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Decode JSON from request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Insert into DB
	query := `INSERT INTO users (username, email, password) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at`
	err = db.DB.QueryRow(query, user.Username, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Don't return password in response
	user.Password = ""

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds models.User

	// Decode JSON request body into creds
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

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
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Compare hashed password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	token, err := auth.GenerateJWT(dbUser.ID)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

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

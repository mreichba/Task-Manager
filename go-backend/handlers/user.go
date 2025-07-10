package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mreichba/task-manager-backend/middleware"
)

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey()).(int)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Authenticated user",
		"user_id": userID,
	})
}

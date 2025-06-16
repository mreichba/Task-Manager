package handlers

import (
	"net/http"
)

// HealthCheckHandler responds with 200 OK when backend is running
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Backend is running!"))
}

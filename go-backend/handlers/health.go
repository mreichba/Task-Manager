package handlers

import (
	"log"
	"net/http"

	"github.com/mreichba/task-manager-backend/db"
)

// HealthCheckHandler responds with 200 OK when backend is running
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Backend is running!"))
}

// DBHealthCheck pings the database to check if it's reachable
func DBHealthCheck(w http.ResponseWriter, r *http.Request) {
	err := db.DB.Ping()
	if err != nil {
		log.Println("DBHealthCheck failed:", err)
		http.Error(w, "Database not reachable", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Database is healthy!"))
}

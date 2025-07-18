package handlers

import (
	"net/http"

	"github.com/mreichba/task-manager-backend/db"
	"github.com/mreichba/task-manager-backend/logger"
)

// HealthCheckHandler responds with 200 OK when backend is running
func HealthCheck(w http.ResponseWriter, r *http.Request) {

	logger.Log.Debug("Entering HealthCheck")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Backend is running!"))
}

// DBHealthCheck pings the database to check if it's reachable
func DBHealthCheck(w http.ResponseWriter, r *http.Request) {

	logger.Log.Debug("Entering DBHealthCheck")

	err := db.DB.Ping()
	if err != nil {
		logger.Log.WithError(err).Warn("DBHealthCheck failed")
		http.Error(w, "Database not reachable", http.StatusInternalServerError)
		return
	}
	logger.Log.Info("DBHealthCheck successful")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Database is healthy!"))
}

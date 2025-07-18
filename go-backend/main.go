package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mreichba/task-manager-backend/config"
	"github.com/mreichba/task-manager-backend/db"
	"github.com/mreichba/task-manager-backend/handlers"
	"github.com/mreichba/task-manager-backend/logger"
	"github.com/mreichba/task-manager-backend/middleware"
	"github.com/sirupsen/logrus"
)

func main() {

	// Load App Config
	config.LoadConfig()

	// Initiate Database Connection
	db.Init()

	// Initiate Router
	router := mux.NewRouter()

	// Server Health Check
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	// Database Health Check
	router.HandleFunc("/db-health", handlers.DBHealthCheck).Methods("GET")

	// Register New User endpoint
	router.HandleFunc("/register", handlers.RegisterUserHandler).Methods("POST")

	// Login User endpoint
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	// Get Current User
	router.Handle("/me", middleware.JWTMiddleware(http.HandlerFunc(handlers.GetCurrentUser)))

	// Start the server
	logger.Info("Server started", logrus.Fields{
		"port": config.AppConfig.ServerPort,
	})

	logger.Info("Environment initialized", logrus.Fields{
		"env": config.AppConfig.Environment,
	})

	err := http.ListenAndServe(":"+config.AppConfig.ServerPort, router)
	if err != nil {
		logger.Fatal("Server failed to start", logrus.Fields{
			"error": err,
		})
	}
}

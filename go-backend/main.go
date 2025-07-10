package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mreichba/task-manager-backend/db"
	"github.com/mreichba/task-manager-backend/handlers"
	"github.com/mreichba/task-manager-backend/middleware"
)

func main() {
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

	router.Handle("/me", middleware.JWTMiddleware(http.HandlerFunc(handlers.GetCurrentUser)))

	// Start the server
	port := "8000"
	fmt.Printf("Server is running on port %v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

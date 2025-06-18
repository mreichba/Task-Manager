package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mreichba/task-manager-backend/db"
	"github.com/mreichba/task-manager-backend/handlers"
)

func main() {
	// Initiate Database Connection
	db.Init()

	// Initiate Router
	router := mux.NewRouter()

	// Server Health Check
	router.HandleFunc("/health", handlers.HealthCheck)

	// Database Health Check
	router.HandleFunc("/db-health", handlers.DBHealthCheck)

	// Start the server
	port := "8080"
	fmt.Printf("Server is running on port %v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mreichba/task-manager-backend/handlers"
)

func main() {

	//initiate router
	router := mux.NewRouter()

	//Health Check
	router.HandleFunc("/health", handlers.HealthCheckHandler)

	// Start the server
	port := "8080"
	fmt.Printf("Server is running on port %v\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

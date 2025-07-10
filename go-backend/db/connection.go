package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/mreichba/task-manager-backend/config"
)

var DB *sql.DB

func Init() {
	// Get DB connection string
	dbURL := config.AppConfig.DatabaseURL
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set in environment")
	}

	// Connect to PostgreSQL
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to open DB connection:", err)
	}
	DB = conn

	// Ping to test connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Connected to the database")
}

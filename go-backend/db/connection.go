package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	// Get DB connection string
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set in environment")
	}

	// Connect to PostgreSQL
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to open DB connection:", err)
	}

	// Ping to test connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Connected to the database")
}

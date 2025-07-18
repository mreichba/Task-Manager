package db

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/mreichba/task-manager-backend/config"
	"github.com/mreichba/task-manager-backend/logger"
	"github.com/sirupsen/logrus"
)

var DB *sql.DB

func Init() {
	// Get DB connection string
	dbURL := config.AppConfig.DatabaseURL
	if dbURL == "" {
		logger.Fatal("DATABASE_URL not set in environment", nil)
	}

	// Connect to PostgreSQL
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		logger.Fatal("Failed to open DB connection:", logrus.Fields{
			"error": err,
			"dbURL": dbURL,
		})
	}
	DB = conn

	// Retry pinging up to 3 times with backoff
	if err := retryPing(DB); err != nil {
		logger.Fatal("Failed to connect to database after retries", logrus.Fields{
			"error": err,
		})
	}

	logger.Info("Connected to the database", nil)
}

func retryPing(db *sql.DB) error {
	var err error
	for i := 1; i <= 3; i++ {
		err = db.Ping()
		if err == nil {
			return nil
		}
		logger.Warn("Ping failed, retrying...", logrus.Fields{
			"attempt": i,
			"error":   err,
		})
		time.Sleep(time.Duration(i) * time.Second)
	}
	return err
}

package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	ServerPort  string
	TokenTTL    time.Duration
}

var AppConfig *Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	AppConfig = &Config{
		DatabaseURL: mustGetEnv("DATABASE_URL"),
		JWTSecret:   mustGetEnv("JWT_SECRET"),
		ServerPort:  getEnv("PORT", "8000"),
		TokenTTL:    24 * time.Hour,
	}
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Required environment variable %s not set", key)
	}
	return value
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

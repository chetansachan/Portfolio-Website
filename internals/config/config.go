package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load() string {
	// 1. Check if the environment variable is already set globally in the system (e.g. Railway / Production)
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		return dbURL
	}

	// 2. Fallback: Load .env file for local development / testing
	if err := godotenv.Load(".env"); err != nil {
		if err2 := godotenv.Load("../.env"); err2 != nil {
			log.Println("Notice: No .env file found; checking system environment variables")
		}
	}

	// 3. Read after loading .env file
	dbURL = os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("Warning: DATABASE_URL is not set in system environment variables or .env file")
	}

	return dbURL
}

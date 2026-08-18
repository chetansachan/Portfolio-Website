package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Load() string {
	// 1. Load .env file for local development (godotenv does not overwrite existing system env vars)
	if err := godotenv.Load(".env"); err != nil {
		if err2 := godotenv.Load("../.env"); err2 != nil {
			log.Println("Notice: No .env file found; checking system environment variables")
		}
	}

	// 2. Check if local system database flag is enabled
	isLocal, _ := strconv.ParseBool(os.Getenv("LOCAL_SYSTEM_DB"))

	var dbURL string
	if isLocal {
		dbURL = os.Getenv("DATABASE_URL_LOCAL")
	}

	// 3. Fallback to DATABASE_URL if LOCAL_SYSTEM_DB is false or DATABASE_URL_LOCAL is empty
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}

	if dbURL == "" {
		log.Println("Warning: DATABASE_URL is not set in system environment variables or .env file")
	}

	return dbURL
}

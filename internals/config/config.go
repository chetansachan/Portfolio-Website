package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load() string {
	// Try multiple paths so the app works from both the project root and cmd/
	if err := godotenv.Load(".env"); err != nil {
		if err2 := godotenv.Load("../.env"); err2 != nil {
			log.Println("Error Loading .env file")
		}
	}

	config := os.Getenv("DATABASE_URL")

	return config
}

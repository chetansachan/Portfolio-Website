package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load() string {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Error Loading .env file")
	}

	config := os.Getenv("DATABASE_URL")

	return config
}

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database_URL string
}

func Load()(*Config) {
	if err := godotenv.Load(); err != nil {
		log.Println("Error Loading Config")
	}

	config := &Config{
		Database_URL: os.Getenv("DATABASE_URL"),
	}

	return config
}

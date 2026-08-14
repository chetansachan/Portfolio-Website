package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database_URL string
}

func Load() *Config {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Error Loading .env file")
	}

	config := &Config{
		Database_URL: os.Getenv("DATABASE_URL"),
	}

	return config
}

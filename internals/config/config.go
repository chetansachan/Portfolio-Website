package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database_URL string
}

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error Loading Config")
		return
	}

	config := os.Getenv("DATABASE_URL")
	
	return config()
}

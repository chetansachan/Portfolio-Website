package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func Connect() (config *pgxpool.Config) {
	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Error Connecting to database")
	}

	if err := conn.Ping(context.Background()); err != nil {
		log.Println("Could not ping the db!", err)
		return
	}

	return
}

package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, config string) (*pgxpool.Config, error) {
	configstg, err := pgxpool.ParseConfig(config)
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.NewWithConfig(ctx, configstg)
	if err != nil {
		log.Fatal("Error Connecting to database")
	}

	if err := conn.Ping(context.Background()); err != nil {
		log.Println("Could not ping the db!", err)
		return nil, err
	}

	return configstg, err
}

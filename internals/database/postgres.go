package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect(ctx context.Context, config string) error {
	configstg, err := pgxpool.ParseConfig(config)
	if err != nil {
		return err
	}

	Pool, err = pgxpool.NewWithConfig(ctx, configstg)
	if err != nil {
		log.Fatalf("Error Connecting to database: %v", err)
	}

	if err := Pool.Ping(context.Background()); err != nil {
		log.Println("Could not ping the db!", err)
		return err
	}

	return err
}

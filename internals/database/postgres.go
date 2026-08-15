package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, config string) error {
	configstg, err := pgxpool.ParseConfig(config)
	if err != nil {
		return err
	}

	Pool, err := pgxpool.NewWithConfig(ctx, configstg)
	if err != nil {
		log.Fatal("Error Connecting to database")
	}

	if err := Pool.Ping(context.Background()); err != nil {
		log.Println("Could not ping the db!", err)
		return err
	}

	return err
}

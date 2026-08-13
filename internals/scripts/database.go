package scripts

import (
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

var migrationFiles embed.FS

func RunMigration() {
	d, err := iofs.New(migrationFiles, "migration")
	if err != nil {
		log.Fatalf("failed to create iofs driver: %v", err)
	}

	dbURL := "postgres://postgres:chetan156@localhost:8080/portfolio_website?sslmode=disable"
	m, err := migrate.NewWithSourceInstance("iofs", d, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrations: %v", err)
	}

	log.Println("Migrations executed successfully!")
}

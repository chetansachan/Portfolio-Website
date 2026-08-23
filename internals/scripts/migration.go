package scripts

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // registers the file:// driver
)

func RunMigration(path string, dbURL string) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	instance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(path, "postgres", instance)
	if err != nil {
		log.Fatal(err)
	}

	version, dirty, err := m.Version()
	if err == nil && dirty {
		targetVersion := int(version) - 1
		if targetVersion < 0 {
			targetVersion = 0
		}
		log.Printf("Detected dirty database schema at version %d in database. Forcing schema version to %d to clear dirty status...", version, targetVersion)
		if err := m.Force(targetVersion); err != nil {
			log.Fatalf("Failed to force migration version to clear dirty status: %v", err)
		}
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

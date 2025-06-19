package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationPath, migrationTable string

	flag.StringVar(&storagePath, "storage", os.Getenv("DB_URL"), "PostgreSQL DSN (or via env DB_URL)")
	flag.StringVar(&migrationPath, "migration", "migrations", "Path to the migration files")
	flag.StringVar(&migrationTable, "table", "migrations", "Name of the migration table")
	flag.Parse()

	if storagePath == "" || migrationPath == "" {
		panic("Missing required flags: storage or migration")
	}

	migrator, err := migrate.New(
		"file://"+migrationPath,
		storagePath,
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create migrator: %v", err))
	}

	if err := migrator.Up(); err != nil && err.Error() != "no change" {
		panic(fmt.Sprintf("failed to apply migrations: %v", err))
	}
}

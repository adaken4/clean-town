package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Needed to read migrations from local file system
)

// RunMigrations applies all pending migrations from the specified directory
// against the provided PostgreSQL database connection.
func RunMigrations(db *sql.DB, migrationsDir string) error {
	// Wrap the existing *sql.DB connection in a migration driver
	// that the migrate library can use.
	migrationDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	// Create a new migrate instance, pointing to the local filesystem
	// for migration files and using the Postgres driver.
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsDir), // e.g., "./migrations"
		"postgres",                              // database name (used internally by migrate)
		migrationDriver,
	)

	// Apply all up migrations in order. If no migrations need to be run,
	// migrate.ErrNoChange will be returned, which we treat as non-fatal.
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Success — all migrations are up to date.
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}
	return nil
}

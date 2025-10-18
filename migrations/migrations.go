package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var migrationsFs embed.FS

type Migrate struct {
	db        *sql.DB
	dbName    string
	migration *migrate.Migrate
}

func (m *Migrate) UP() error {
	if err := m.migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	fmt.Println("Migrations applied successfully")
	return nil
}

func (m *Migrate) Rollback() error {
	if err := m.migration.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	fmt.Println("Last migration rolled back successfully")
	return nil
}

func (m *Migrate) Close() error {
	if m.migration != nil {
		source, driver := m.migration.Close()
		if source != nil {
			return fmt.Errorf("failed to close source: %w", source)
		}
		if driver != nil {
			return fmt.Errorf("failed to close driver: %w", driver)
		}
	}
	return nil
}

func NewMigrate(db *sql.DB, dbName string) (*Migrate, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create migration driver: %w", err)
	}

	source, err := iofs.New(migrationsFs, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to load migration files: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, dbName, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migrate: %w", err)
	}

	return &Migrate{
		db:        db,
		dbName:    dbName,
		migration: m,
	}, nil
}

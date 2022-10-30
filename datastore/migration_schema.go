package datastore

import (
	"database/sql"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var fs embed.FS

func validateSchema(db *sql.DB) error {
	sourceDriver, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}
	targetInstance, err := postgres.WithInstance(db, new(postgres.Config))
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", targetInstance)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

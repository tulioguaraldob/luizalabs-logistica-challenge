package postgres

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(dataSourceName string) error {
	m, err := migrate.New("file://./config/database/migrations", dataSourceName)
	if err != nil {
		return err
	}

	m.Up()
	return nil
}

func downMigrations(dataSourceName string) error {
	m, err := migrate.New("file://./config/database/migrations", dataSourceName)
	if err != nil {
		return err
	}

	m.Down()
	return nil
}

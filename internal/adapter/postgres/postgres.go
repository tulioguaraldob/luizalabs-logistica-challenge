package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// New opens a Postgres connection using the native lib
func New(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Printf("Failed to open Postgres connection. Details: %s\n", err.Error())
		return nil, err
	}

	// if err := downMigrations(dataSourceName); err != nil {
	// 	log.Printf("Failed to apply Postgres down migrations. Details: %s\n", err.Error())
	// 	return nil, err
	// }

	if err := runMigrations(dataSourceName); err != nil {
		log.Printf("Failed to apply Postgres migrations. Details: %s\n", err.Error())
		return nil, err
	}

	return db, nil
}

// Close will close an open Postgres connection if the connection is valid and not nil
//
// Must be used in the main entry of the application
func Close(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("must have a valid Postgres connection to be closed")
	}

	return db.Close()
}

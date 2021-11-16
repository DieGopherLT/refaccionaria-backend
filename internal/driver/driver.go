package driver

import (
	"database/sql"
)

// CreateDatabaseConnection by applying polymorphism, creates and returns a database.
//
// It's something like a builder director.
func CreateDatabaseConnection(db DatabasePoolConnectionBuilder, dsn string) (DatabasePoolConnectionBuilder, error) {
	err := db.CreatePool(dsn)
	if err != nil {
		return nil, err
	}
	db.ConfigurePool()
	return db, nil
}

// TestDatabaseConnection tests database connection by pinging to it
func TestDatabaseConnection(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}

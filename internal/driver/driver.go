package driver

import (
	"database/sql"
)

// ConnectSQL by applying polymorphism, creates and returns a database.
// It's something like a builder director
func ConnectSQL(db SQLDatabase, dsn string) (SQLDatabase, error) {
	err := db.CreatePool(dsn)
	if err != nil {
		return nil, err
	}
	db.ConfigurePool()
	return db, nil
}

func TestSQL(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}

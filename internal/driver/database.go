package driver

import (
	"database/sql"
)

// SQLDatabase A database builder
type SQLDatabase interface {
	CreatePool(dsn string) error
	ConfigurePool()
	GetPool() *sql.DB
}

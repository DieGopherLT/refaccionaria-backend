package driver

import (
	"database/sql"
)

// DatabasePoolConnectionBuilder A database builder
type DatabasePoolConnectionBuilder interface {
	CreatePool(connectionURL string) error
	ConfigurePool()
	GetPool() *sql.DB
}

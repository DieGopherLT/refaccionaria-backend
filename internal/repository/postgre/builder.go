package postgre

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	maxOpenDbConn = 10
	maxIdleDbConn = 5
	maxDbLifeTime = 5 * time.Minute
)

// Builder a builder for a PostgreSQL database
type Builder struct {
	db *sql.DB
}

func NewBuilder() *Builder {
	return &Builder{}
}

// CreatePool creates database pool connection for postgres
func (r *Builder) CreatePool(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	r.db = db
	return nil
}

// ConfigurePool configures the pool
func (r *Builder) ConfigurePool() {
	r.db.SetMaxOpenConns(maxOpenDbConn)
	r.db.SetMaxIdleConns(maxIdleDbConn)
	r.db.SetConnMaxLifetime(maxDbLifeTime)
}

// GetPool returns the pool
func (r *Builder) GetPool() *sql.DB {
	return r.db
}

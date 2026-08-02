package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/light-group/light-ex-backend/internal/platform/config"
)

const (
	driverName     = "pgx"
	connectTimeout = 5 * time.Second
)

// Database represents the application's PostgreSQL connection.
type Database struct {
	db *sql.DB
}

// New creates and verifies a PostgreSQL connection pool.
func New(cfg config.DatabaseConfig) (*Database, error) {

	db, err := sql.Open(driverName, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		connectTimeout,
	)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &Database{
		db: db,
	}, nil
}

// SQL returns the underlying *sql.DB.
//
// This is intended for repositories and migration tooling.
func (d *Database) SQL() *sql.DB {
	return d.db
}

// Health verifies database connectivity.
func (d *Database) Health(ctx context.Context) error {

	if d == nil || d.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	return d.db.PingContext(ctx)
}

// Close gracefully closes the connection pool.
func (d *Database) Close() error {

	if d == nil || d.db == nil {
		return nil
	}

	return d.db.Close()
}

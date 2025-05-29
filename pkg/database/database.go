package database

import (
	"context"
	"database/sql"
	"fmt"
	"maxwellzp/golang-db-api-context/pkg/config"
	"time"
)

type DB struct {
	*sql.DB
}

func New(ctx context.Context, cfg config.DbConfig) (*DB, error) {
	// The context timeout ensures that if the database connection takes too long to
	// establish, the operation will fail fast rather than hang indefinitely
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	//if _, err := db.ExecContext(ctx, "SELECT SLEEP(1000)"); err != nil {
	//	return nil, fmt.Errorf("sleeping: %w", err)
	//}

	// Actually tests the connection with context timeout
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &DB{db}, nil
}

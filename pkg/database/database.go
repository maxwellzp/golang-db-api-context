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
	// Verify DSN with context
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

	return &DB{db}, nil
}

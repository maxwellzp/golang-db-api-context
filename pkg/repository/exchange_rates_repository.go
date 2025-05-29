package repository

import (
	"context"
	"fmt"

	"maxwellzp/golang-db-api-context/pkg/database"
	"maxwellzp/golang-db-api-context/pkg/models"
)

type ExchangeRatesRepository struct {
	db *database.DB
}

func NewExchangeRatesRepository(db *database.DB) *ExchangeRatesRepository {
	return &ExchangeRatesRepository{db: db}
}

func (r *ExchangeRatesRepository) StoreRates(ctx context.Context, rates []models.ExchangeRate) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO exchange_rates 
			(currency_code, base_currency_code, rate, date_updated) 
		VALUES (?, ?, ?, ?)
	`)

	if err != nil {
		return fmt.Errorf("preparing statement: %w", err)
	}
	defer stmt.Close()

	for _, rate := range rates {
		if _, err := stmt.ExecContext(ctx,
			rate.CurrencyCode,
			rate.BaseCurrencyCode,
			rate.Rate,
			rate.DateUpdated,
		); err != nil {
			return fmt.Errorf("inserting rate: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

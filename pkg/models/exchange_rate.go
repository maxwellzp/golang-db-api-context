package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type ExchangeRate struct {
	CurrencyCode     string          `json:"currency_code"`
	BaseCurrencyCode string          `json:"base_currency_code"`
	Rate             decimal.Decimal `json:"rate"`
	DateUpdated      time.Time       `json:"date_updated"`
}

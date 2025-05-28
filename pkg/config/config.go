package config

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
	"time"
)

type DbConfig struct {
	DSN     string
	Timeout time.Duration
}

type ApiConfig struct {
	URL     string
	ApiKey  string
	Timeout time.Duration
}

type Config struct {
	API ApiConfig
	DB  DbConfig
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("loading .env file: %w", err)
	}

	apiTimeout, err := time.ParseDuration(os.Getenv("API_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("parsing API timeout: %w", err)
	}

	dbTimeout, err := time.ParseDuration(os.Getenv("MYSQL_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("parsing DB timeout: %w", err)
	}

	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("MYSQL_USER")
	cfg.Passwd = os.Getenv("MYSQL_PASSWORD")
	cfg.Net = "tcp"
	cfg.Addr = fmt.Sprintf("127.0.0.1:%s", os.Getenv("MYSQL_PORT"))
	cfg.DBName = os.Getenv("MYSQL_DATABASE")

	return &Config{
		API: ApiConfig{
			URL:     os.Getenv("API_URL"),
			ApiKey:  os.Getenv("API_KEY"),
			Timeout: apiTimeout,
		},
		DB: DbConfig{
			DSN:     cfg.FormatDSN(),
			Timeout: dbTimeout,
		},
	}, nil
}

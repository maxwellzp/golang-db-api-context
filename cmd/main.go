package main

import (
	"context"
	"log"
	"maxwellzp/golang-db-api-context/pkg/config"
	"maxwellzp/golang-db-api-context/pkg/database"
	"maxwellzp/golang-db-api-context/pkg/exchangerate"
	"maxwellzp/golang-db-api-context/pkg/repository"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.New(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Setup graceful shutdown
	shutdownCtx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,  // Triggered when user presses Ctrl+C in terminal
		syscall.SIGTERM, // Default signal sent by kill command and container orchestrators
	)
	defer stop()

	db.Ping()

	// Initialize services
	apiClient := exchangerate.NewClient(cfg.API)

	repo := repository.NewExchangeRatesRepository(db)

	// Context timeout (API logic)
	opCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	rates, err := apiClient.GetExchangeRates(opCtx, "USD")
	if err != nil {
		log.Fatalf("Failed to get exchange rates: %v", err)
	}

	if err = repo.StoreRates(shutdownCtx, rates); err != nil {
		log.Printf("Failed to store rates: %v", err)
		os.Exit(1)
	}
}

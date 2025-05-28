package main

import (
	"context"
	"fmt"
	"log"
	"maxwellzp/golang-db-api-context/pkg/config"
	"maxwellzp/golang-db-api-context/pkg/database"
	"maxwellzp/golang-db-api-context/pkg/exchangerate"
	"time"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("ApiKey: ", cfg.API.ApiKey)
	log.Println("URL: ", cfg.API.URL)
	log.Println("Timeout: ", cfg.API.Timeout)

	log.Println("DSN: ", cfg.DB.DSN)
	log.Println("DB timeout: ", cfg.DB.Timeout)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.New(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	db.Ping()

	// Initialize services
	apiClient := exchangerate.NewClient(cfg.API)

	opCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	rates, err := apiClient.GetExchangeRates(opCtx, "USD")
	if err != nil {
		log.Fatalf("Failed to get exchange rates: %v", err)
	}
	fmt.Println(rates)
}

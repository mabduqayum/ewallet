package main

import (
	"context"
	"ewallet/internal/config"
	"ewallet/internal/database"
	"ewallet/scripts"
	"log"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("development")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to the database
	ctx := context.Background()
	db, err := database.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Seed clients
	clients := scripts.SeedClients(ctx, db)

	// Seed wallets
	wallets := scripts.SeedWallets(ctx, db, clients)

	// Seed transactions
	scripts.SeedTransactions(ctx, db, wallets)

	log.Println("Seeding completed successfully")
}

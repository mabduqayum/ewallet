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

	// Run migrations if needed
	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Seed data
	if err := scripts.SeedData(ctx, db.GetPool()); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	log.Println("Seeding completed successfully")
}

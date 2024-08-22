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
	cfg, err := config.LoadConfig("development")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx := context.Background()
	db, err := database.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	if err := scripts.SeedData(ctx, db.GetPool()); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	log.Println("Seeding completed successfully")
}

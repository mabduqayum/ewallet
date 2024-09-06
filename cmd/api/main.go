package main

import (
	"context"
	"log"
	"os"

	"github.com/mabduqayum/ewallet/internal/config"
	"github.com/mabduqayum/ewallet/internal/database"
	"github.com/mabduqayum/ewallet/internal/server"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	cfg, err := config.LoadConfig(env)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	ctx := context.Background()
	db, err := database.New(ctx, &cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	s := server.New(&cfg.Server, db)
	s.RegisterFiberRoutes()

	log.Printf("Starting server on %s", cfg.Server.Address())
	if err := s.Listen(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

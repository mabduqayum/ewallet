package main

import (
	"ewallet/internal/config"
	"ewallet/internal/database"
	"ewallet/internal/server"
	"log"
	"os"

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

	db := database.New(cfg)
	defer db.Close()

	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	s := server.New(cfg, db)
	s.RegisterFiberRoutes()

	log.Printf("Starting server on %s", cfg.Server.Address())
	if err := s.Listen(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

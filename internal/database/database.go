package database

import (
	"context"
	"errors"
	"ewallet/internal/config"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close()

	RunMigrations() error
	GetPool() *pgxpool.Pool
}

type service struct {
	db     *pgxpool.Pool
	config *config.Config
}

var dbInstance *service

func New(ctx context.Context, cfg *config.Config) (Service, error) {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance, nil
	}

	connStr := cfg.Database.ConnectionString()

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}

	// Set additional pool configuration if needed
	// poolConfig.MaxConns = 10

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	dbInstance = &service{
		db:     db,
		config: cfg,
	}
	return dbInstance, nil
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Printf("Database health check failed: %v", err)
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get pool stats
	poolStats := s.db.Stat()
	stats["total_connections"] = strconv.Itoa(int(poolStats.TotalConns()))
	stats["acquired_connections"] = strconv.Itoa(int(poolStats.AcquiredConns()))
	stats["idle_connections"] = strconv.Itoa(int(poolStats.IdleConns()))

	// Evaluate stats to provide a health message
	if poolStats.TotalConns() > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if poolStats.AcquiredConns() > poolStats.TotalConns()/2 {
		stats["message"] = "More than half of the connections are in use, indicating high load."
	}

	return stats
}

// Close closes the database connection pool.
func (s *service) Close() {
	log.Printf("Closing connection pool to database: %s", s.config.Database.DBName)
	s.db.Close()
}

func (s *service) RunMigrations() error {
	migrationPath := "file://migrations"
	m, err := migrate.New(migrationPath, s.config.Database.ConnectionString())
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

// GetPool a method to get the database pool if needed in other parts of your application
func (s *service) GetPool() *pgxpool.Pool {
	return s.db
}

package database

import (
	"context"
	"fmt"
	"log"
	"time"
	"time_capsule_memories/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// Connect initializes the database connection pool using config values.
func Connect() error {
	cfg := config.GetConfig()

	log.Println("Connecting to database:", cfg.PostgresURL)

	var err error
	DB, err = pgxpool.New(context.Background(), cfg.PostgresURL)
	if err != nil {
		return fmt.Errorf("failed to create DB pool: %w", err)
	}

	// Test DB connection with a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if DB == nil {
		return fmt.Errorf("database pool is nil")
	}

	if err := DB.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")
	return nil
}

// Close shuts down the database connection pool.
func Close() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}

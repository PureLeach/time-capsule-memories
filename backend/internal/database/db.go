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

func Connect(config config.Config) error {
	log.Println("Подключение к базе данных: " + config.PostgresURL)
	var err error
	DB, err = pgxpool.New(context.Background(), config.PostgresURL)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных: %w", err)
	}

	// Проверяем подключение к базе данных
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := DB.Ping(ctx); err != nil {
		return fmt.Errorf("не удалось выполнить ping к базе данных: %w", err)
	}

	return nil
}

// Close закрывает соединение с базой данных
func Close() {
	if DB != nil {
		DB.Close()
	}
}

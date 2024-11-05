package database

import (
	"context"
	"fmt"
	"reflect"
	"time"
	"time_capsule_memories/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// Connect устанавливает соединение с базой данных PostgreSQL
func Connect(config config.Config) error {
	// Формируем строку подключения
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.User, config.Password, config.Host, config.Port, config.DBName)
	fmt.Printf("databaseURL: %#v Type: %v\n", databaseURL, reflect.TypeOf(databaseURL))

	// Подключаемся к базе данных
	var err error
	DB, err = pgxpool.New(context.Background(), databaseURL)
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

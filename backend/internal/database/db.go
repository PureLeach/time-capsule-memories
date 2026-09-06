// Package database owns the Postgres connection pool.
package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"time_capsule_memories/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

const pingTimeout = 5 * time.Second

func New(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("create db pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	// Logged from the parsed config, not the URL, which carries the password.
	slog.Info("database connection established",
		"host", pool.Config().ConnConfig.Host,
		"database", pool.Config().ConnConfig.Database,
	)
	return pool, nil
}

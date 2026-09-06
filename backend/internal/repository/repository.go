// Package repository holds the SQL data access layer.
package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const dbTimeout = 5 * time.Second

// dbPool is the subset of *pgxpool.Pool the repositories use. pgxmock satisfies
// it, which is what makes them testable without a database.
type dbPool interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

package repository

import (
	"context"
	"log/slog"

	"time_capsule_memories/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Feedback struct {
	pool *pgxpool.Pool
}

func NewFeedback(pool *pgxpool.Pool) *Feedback {
	return &Feedback{pool: pool}
}

func (r *Feedback) Create(ctx context.Context, feedback *models.CreateFeedbackRequest) (*models.FeedbackResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := `
	INSERT INTO users_feedback (message)
	VALUES ($1)
	RETURNING id, created_at, message;
    `

	created := &models.FeedbackResponse{}

	err := r.pool.QueryRow(
		ctx,
		query,
		feedback.Message,
	).Scan(
		&created.ID,
		&created.CreatedAt,
		&created.Message,
	)
	if err != nil {
		slog.Error("failed to create feedback", "error", err)
		return nil, err
	}

	return created, nil
}

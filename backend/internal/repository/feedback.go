package repository

import (
	"context"
	"fmt"

	"time_capsule_memories/internal/models"
)

type Feedback struct {
	pool dbPool
}

func NewFeedback(pool dbPool) *Feedback {
	return &Feedback{pool: pool}
}

func (r *Feedback) Create(ctx context.Context, feedback *models.CreateFeedbackRequest) (*models.FeedbackResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	const query = `
	INSERT INTO users_feedback (message)
	VALUES ($1)
	RETURNING id, created_at, message;
	`

	created := &models.FeedbackResponse{}
	err := r.pool.QueryRow(ctx, query, feedback.Message).Scan(
		&created.ID,
		&created.CreatedAt,
		&created.Message,
	)
	if err != nil {
		return nil, fmt.Errorf("insert feedback: %w", err)
	}

	return created, nil
}

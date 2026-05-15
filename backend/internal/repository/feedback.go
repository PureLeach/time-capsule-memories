package repository

import (
	"context"
	"log/slog"

	"time_capsule_memories/internal/database"
	"time_capsule_memories/internal/models"
)

// CreateUserFeedback creates a new feedback from a user and returns the created feedback data.
func CreateUserFeedback(ctx context.Context, feedback *models.CreateFeedbackRequest) (createdFeedback *models.FeedbackResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	query := `
	INSERT INTO users_feedback (message)
	VALUES ($1)
	RETURNING id, created_at, message;
    `

	createdFeedback = &models.FeedbackResponse{}

	// Execute the query to insert the feedback into the database
	err = database.DB.QueryRow(
		ctx,
		query,
		feedback.Message,
	).Scan(
		&createdFeedback.ID,
		&createdFeedback.CreatedAt,
		&createdFeedback.Message,
	)

	if err != nil {
		slog.Error("failed to create feedback", "error", err)
		return nil, err
	}

	return createdFeedback, nil
}

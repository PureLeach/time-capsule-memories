package repository

import (
	"context"
	"log"

	"time_capsule_memories/internal/database"
	"time_capsule_memories/internal/models"
)

// CreateUserFeedback creates a new feedback from a user and returns the created feedback data.
func CreateUserFeedback(feedback *models.CreateFeedbackRequest) (createdFeedback *models.FeedbackResponse, err error) {
	query := `
	INSERT INTO users_feedback (message)
	VALUES ($1)
	RETURNING id, created_at, message;
    `

	createdFeedback = &models.FeedbackResponse{}

	// Execute the query to insert the feedback into the database
	err = database.DB.QueryRow(
		context.Background(),
		query,
		feedback.Message,
	).Scan(
		&createdFeedback.ID,
		&createdFeedback.CreatedAt,
		&createdFeedback.Message,
	)

	if err != nil {
		log.Printf("Error creating feedback in the database: %v", err)
		return nil, err
	}

	return createdFeedback, nil
}

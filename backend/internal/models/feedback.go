package models

import "time"

type CreateFeedbackRequest struct {
	Message string `json:"message" example:"Test Message" validate:"required,max=4096"`
}

type FeedbackResponse struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Message   string    `json:"message"`
}

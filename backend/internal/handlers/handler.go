// Package handlers implements the HTTP transport layer.
package handlers

import (
	"context"

	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/services"
)

type pinger interface {
	Ping(ctx context.Context) error
}

type feedbackCreator interface {
	Create(ctx context.Context, feedback *models.CreateFeedbackRequest) (*models.FeedbackResponse, error)
}

type Handler struct {
	db           pinger
	store        services.ObjectStore
	capsuleRepo  services.CapsuleRepository
	feedbackRepo feedbackCreator
	mailer       services.Mailer
}

func New(
	db pinger,
	store services.ObjectStore,
	capsuleRepo services.CapsuleRepository,
	feedbackRepo feedbackCreator,
	mailer services.Mailer,
) *Handler {
	return &Handler{
		db:           db,
		store:        store,
		capsuleRepo:  capsuleRepo,
		feedbackRepo: feedbackRepo,
		mailer:       mailer,
	}
}

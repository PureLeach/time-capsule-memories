package handlers

import (
	"time_capsule_memories/internal/repository"
	"time_capsule_memories/internal/services"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool         *pgxpool.Pool
	store        services.ObjectStore
	capsuleRepo  services.CapsuleRepository
	feedbackRepo *repository.Feedback
	mailer       services.Mailer
}

func New(
	pool *pgxpool.Pool,
	store services.ObjectStore,
	capsuleRepo services.CapsuleRepository,
	feedbackRepo *repository.Feedback,
	mailer services.Mailer,
) *Handler {
	return &Handler{
		pool:         pool,
		store:        store,
		capsuleRepo:  capsuleRepo,
		feedbackRepo: feedbackRepo,
		mailer:       mailer,
	}
}

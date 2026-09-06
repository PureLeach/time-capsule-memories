package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"time_capsule_memories/internal/models"
)

const revertTimeout = 5 * time.Second

type CapsuleRepository interface {
	Create(ctx context.Context, capsule *models.CreateCapsuleRequest) (*models.CapsuleResponse, error)
	ClaimDue(ctx context.Context, limit int) ([]*models.CapsuleResponse, error)
	SetStatus(ctx context.Context, capsuleID int, status models.CapsuleStatus) error
}

type ObjectStore interface {
	GetFilesInDirectory(ctx context.Context, directoryUUID string) ([]models.FileObject, error)
	PresignUpload(ctx context.Context, directoryUUID, contentType string, expiration time.Duration) (*models.PresignedUpload, error)
	Ping(ctx context.Context) error
}

type Mailer interface {
	Send(ctx context.Context, subject, body, to string, attachments []models.FileObject) error
}

type CapsuleService struct {
	repo   CapsuleRepository
	store  ObjectStore
	mailer Mailer
}

func NewCapsuleService(repo CapsuleRepository, store ObjectStore, mailer Mailer) *CapsuleService {
	return &CapsuleService{repo: repo, store: store, mailer: mailer}
}

// Process delivers one claimed capsule. A failure before the mail leaves returns
// the row to 'waiting' for the next tick; after it leaves the row is never
// reverted, since a duplicate delivery is worse than a stuck row.
func (s *CapsuleService) Process(ctx context.Context, capsule *models.CapsuleResponse) error {
	slog.Info("processing capsule", "capsule_id", capsule.ID)

	var attachments []models.FileObject
	if folder, ok := capsule.AttachmentsFolder(); ok {
		files, err := s.store.GetFilesInDirectory(ctx, folder)
		if err != nil {
			slog.Error("failed to retrieve capsule attachments",
				"capsule_id", capsule.ID,
				"folder_uuid", folder,
				"error", err,
			)
			s.revertToWaiting(ctx, capsule.ID, "object_store_fetch")
			return fmt.Errorf("retrieve attachments: %w", err)
		}
		attachments = files
	}

	subject := fmt.Sprintf("You've received a time capsule from %s", capsule.SenderName)

	if err := s.mailer.Send(ctx, subject, capsule.Message, capsule.RecipientEmail, attachments); err != nil {
		s.revertToWaiting(ctx, capsule.ID, "smtp_send")
		return fmt.Errorf("send capsule email: %w", err)
	}

	if err := s.repo.SetStatus(ctx, capsule.ID, models.StatusDone); err != nil {
		slog.Error("capsule delivered but status update failed; row left in progress",
			"capsule_id", capsule.ID, "error", err,
		)
		return fmt.Errorf("mark capsule done: %w", err)
	}

	slog.Info("capsule processing complete", "capsule_id", capsule.ID)
	return nil
}

func (s *CapsuleService) revertToWaiting(ctx context.Context, capsuleID int, after string) {
	// The parent may already be canceled by the dispatch timeout, which would
	// leave the capsule stuck in progress.
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), revertTimeout)
	defer cancel()

	if err := s.repo.SetStatus(ctx, capsuleID, models.StatusWaiting); err != nil {
		slog.Error("failed to revert capsule to waiting; row stuck in progress",
			"capsule_id", capsuleID,
			"after", after,
			"error", err,
		)
	}
}

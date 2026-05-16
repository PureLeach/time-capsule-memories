package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"time_capsule_memories/internal/models"
)

type CapsuleRepository interface {
	Create(ctx context.Context, capsule *models.CreateCapsuleRequest) (*models.CapsuleResponse, error)
	ClaimDue(ctx context.Context, limit int) ([]*models.CapsuleResponse, error)
	SetStatus(ctx context.Context, capsuleID int, status string) error
}

type ObjectStore interface {
	GetFilesInDirectory(ctx context.Context, directoryUUID string) ([]models.FileObject, error)
	GeneratePresignedUploadURL(ctx context.Context, objectName string, expiration time.Duration) (string, error)
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

func (s *CapsuleService) Process(ctx context.Context, capsule *models.CapsuleResponse) error {
	slog.Info("processing capsule", "capsule_id", capsule.ID)

	if *capsule.FilesFolderUUID != "" {
		slog.Debug("capsule has attachments folder",
			"capsule_id", capsule.ID,
			"folder_uuid", *capsule.FilesFolderUUID,
		)
	}

	attachments, err := s.store.GetFilesInDirectory(ctx, *capsule.FilesFolderUUID)
	if err != nil {
		slog.Error("failed to retrieve capsule attachments",
			"capsule_id", capsule.ID,
			"folder_uuid", *capsule.FilesFolderUUID,
			"error", err,
		)
		s.revertToWaiting(ctx, capsule.ID, "minio_fetch")
		return fmt.Errorf("error retrieving files from MinIO: %v", err)
	}

	subject := fmt.Sprintf("You've received a time capsule from %s", capsule.SenderName)

	if err := s.mailer.Send(ctx, subject, capsule.Message, capsule.RecipientEmail, attachments); err != nil {
		s.revertToWaiting(ctx, capsule.ID, "smtp_send")
		return fmt.Errorf("error sending email: %v", err)
	}

	if err := s.repo.SetStatus(ctx, capsule.ID, "done"); err != nil {
		slog.Error("capsule delivered but status update failed; row left in progress",
			"capsule_id", capsule.ID, "error", err,
		)
		return fmt.Errorf("error updating capsule status: %v", err)
	}

	slog.Info("capsule processing complete", "capsule_id", capsule.ID)
	return nil
}

func (s *CapsuleService) revertToWaiting(ctx context.Context, capsuleID int, after string) {
	if err := s.repo.SetStatus(ctx, capsuleID, "waiting"); err != nil {
		slog.Error("failed to revert capsule to waiting; row stuck in progress",
			"capsule_id", capsuleID,
			"after", after,
			"error", err,
		)
	}
}

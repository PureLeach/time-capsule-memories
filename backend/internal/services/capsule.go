package services

import (
	"context"
	"fmt"
	"log/slog"
	"time_capsule_memories/internal/minio_client"
	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/repository"
)

func ProcessCapsule(ctx context.Context, capsule *models.CapsuleResponse) error {
	slog.Info("processing capsule", "capsule_id", capsule.ID)

	if *capsule.FilesFolderUUID != "" {
		slog.Debug("capsule has attachments folder",
			"capsule_id", capsule.ID,
			"folder_uuid", *capsule.FilesFolderUUID,
		)
	}

	attachments, err := minio_client.GetFilesInDirectory(ctx, *capsule.FilesFolderUUID)
	if err != nil {
		slog.Error("failed to retrieve capsule attachments",
			"capsule_id", capsule.ID,
			"folder_uuid", *capsule.FilesFolderUUID,
			"error", err,
		)
		revertToWaiting(ctx, capsule.ID, "minio_fetch")
		return fmt.Errorf("error retrieving files from MinIO: %v", err)
	}

	subject := fmt.Sprintf("You've received a time capsule from %s", capsule.SenderName)

	if err := SendEmail(ctx, subject, capsule.Message, capsule.RecipientEmail, attachments); err != nil {
		revertToWaiting(ctx, capsule.ID, "smtp_send")
		return fmt.Errorf("error sending email: %v", err)
	}

	if err := repository.SetCapsuleStatus(ctx, capsule.ID, "done"); err != nil {
		slog.Error("capsule delivered but status update failed; row left in progress",
			"capsule_id", capsule.ID, "error", err,
		)
		return fmt.Errorf("error updating capsule status: %v", err)
	}

	slog.Info("capsule processing complete", "capsule_id", capsule.ID)
	return nil
}

func revertToWaiting(ctx context.Context, capsuleID int, after string) {
	if err := repository.SetCapsuleStatus(ctx, capsuleID, "waiting"); err != nil {
		slog.Error("failed to revert capsule to waiting; row stuck in progress",
			"capsule_id", capsuleID,
			"after", after,
			"error", err,
		)
	}
}

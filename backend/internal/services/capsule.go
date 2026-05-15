package services

import (
	"fmt"
	"log/slog"
	"time_capsule_memories/internal/minio_client"
	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/repository"
)

// ProcessCapsule processes the capsule, retrieves files from MinIO, sends an email with the files,
// and updates the capsule's status to "done" in the database.
func ProcessCapsule(capsule *models.CapsuleResponse) error {
	slog.Info("processing capsule", "capsule_id", capsule.ID)

	if *capsule.FilesFolderUUID != "" {
		slog.Debug("capsule has attachments folder",
			"capsule_id", capsule.ID,
			"folder_uuid", *capsule.FilesFolderUUID,
		)
	}

	// Retrieve files from MinIO
	attachments, err := minio_client.GetFilesInDirectory(*capsule.FilesFolderUUID)
	if err != nil {
		slog.Error("failed to retrieve capsule attachments",
			"capsule_id", capsule.ID,
			"folder_uuid", *capsule.FilesFolderUUID,
			"error", err,
		)
		return fmt.Errorf("error retrieving files from MinIO: %v", err)
	}

	// Prepare the email subject and send it with the message and attachments
	subject := fmt.Sprintf("You've received a time capsule from %s", capsule.SenderName)

	if err := SendEmail(subject, capsule.Message, capsule.RecipientEmail, attachments); err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	// Update the capsule's status to "done" in the database
	if err := repository.UpdateCapsuleStatusByID(capsule.ID, "done"); err != nil {
		return fmt.Errorf("error updating capsule status: %v", err)
	}

	slog.Info("capsule processing complete", "capsule_id", capsule.ID)
	return nil
}

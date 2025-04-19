package services

import (
	"fmt"
	"log"
	"time_capsule_memories/internal/minio_client"
	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/repository"
)

// ProcessCapsule processes the capsule, retrieves files from MinIO, sends an email with the files,
// and updates the capsule's status to "done" in the database.
func ProcessCapsule(capsule *models.CapsuleResponse) error {
	log.Printf("Processing capsule with ID: %d", capsule.ID)

	// If the FilesFolderUUID is provided, log it.
	if *capsule.FilesFolderUUID != "" {
		log.Printf("Files Folder UUID: %s", *capsule.FilesFolderUUID)
	}

	// Retrieve files from MinIO
	attachments, err := minio_client.GetFilesInDirectory(*capsule.FilesFolderUUID)
	if err != nil {
		// Fatal log is too strong here, use a normal error log and return the error instead.
		log.Printf("Error retrieving files from directory %s: %v", *capsule.FilesFolderUUID, err)
		return fmt.Errorf("error retrieving files from MinIO: %v", err)
	}

	// Prepare the email subject and send it with the message and attachments
	subject := fmt.Sprintf("You've received a time capsule from %s", capsule.SenderName)

	if err := SendEmail(subject, capsule.Message, capsule.RecipientEmail, attachments); err != nil {
		log.Printf("Error details: %v", err)
		return fmt.Errorf("error sending email: %v", err)
	}

	// Update the capsule's status to "done" in the database
	if err := repository.UpdateCapsuleStatusByID(capsule.ID, "done"); err != nil {
		return fmt.Errorf("error updating capsule status: %v", err)
	}

	log.Println("Capsule processing completed")
	return nil
}

package repository

import (
	"context"
	"log"

	"time_capsule_memories/internal/database"
	"time_capsule_memories/internal/models"
)

func CreateCapsule(capsule *models.CreateCapsuleRequest) (*models.CapsuleResponse, error) {
	query := `
	INSERT INTO capsules (sender_name, send_at, message, recipient_email, recipient_tg_username, files_folder_UUID)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, sender_name, created_at, send_at, message, recipient_email, recipient_tg_username, files_folder_UUID;
    `

	createdCapsule := &models.CapsuleResponse{}

	err := database.DB.QueryRow(
		context.Background(),
		query,
		capsule.SenderName,
		capsule.SendAt,
		capsule.Message,
		capsule.RecipientEmail,
		capsule.RecipientTgUsername,
		capsule.FilesFolderUUID,
	).Scan(
		&createdCapsule.ID,
		&createdCapsule.SenderName,
		&createdCapsule.CreatedAt,
		&createdCapsule.SendAt,
		&createdCapsule.Message,
		&createdCapsule.RecipientEmail,
		&createdCapsule.RecipientTgUsername,
		&createdCapsule.FilesFolderUUID,
	)

	if err != nil {
		log.Printf("Error creating capsule in the database: %v", err)
		return nil, err
	}

	return createdCapsule, nil
}

package repository

import (
	"context"
	"log"
	"time"

	"time_capsule_memories/internal/database"
	"time_capsule_memories/internal/models"
)

// CreateCapsule creates a new capsule in the database and returns the created capsule data.
func CreateCapsule(capsule *models.CreateCapsuleRequest) (createdCapsule *models.CapsuleResponse, err error) {
	query := `
	INSERT INTO capsules (sender_name, send_at, message, recipient_email, files_folder_UUID)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, sender_name, created_at, send_at, message, recipient_email, files_folder_UUID, status;
    `

	createdCapsule = &models.CapsuleResponse{}

	err = database.DB.QueryRow(
		context.Background(),
		query,
		capsule.SenderName,
		capsule.SendAt,
		capsule.Message,
		capsule.RecipientEmail,
		capsule.FilesFolderUUID,
	).Scan(
		&createdCapsule.ID,
		&createdCapsule.SenderName,
		&createdCapsule.CreatedAt,
		&createdCapsule.SendAt,
		&createdCapsule.Message,
		&createdCapsule.RecipientEmail,
		&createdCapsule.FilesFolderUUID,
		&createdCapsule.Status,
	)

	if err != nil {
		log.Printf("Error creating capsule in the database: %v", err)
		return nil, err
	}

	return createdCapsule, nil
}

// GetCapsulesByToday retrieves all capsules scheduled for today with a "waiting" status.
func GetCapsulesByToday() (capsules []*models.CapsuleResponse, err error) {
	currentDate := time.Now().Format("2006-01-02")

	query := `
	SELECT id, sender_name, created_at, send_at, message, recipient_email, files_folder_UUID, status
	FROM capsules
	WHERE send_at::date = $1 AND status = 'waiting';
	`

	rows, err := database.DB.Query(
		context.Background(),
		query,
		currentDate,
	)
	if err != nil {
		log.Printf("Error querying capsules from the database: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Mapping rows to capsule objects
	for rows.Next() {
		capsule := &models.CapsuleResponse{}
		if err := rows.Scan(
			&capsule.ID,
			&capsule.SenderName,
			&capsule.CreatedAt,
			&capsule.SendAt,
			&capsule.Message,
			&capsule.RecipientEmail,
			&capsule.FilesFolderUUID,
			&capsule.Status,
		); err != nil {
			log.Printf("Error scanning row into CapsuleResponse: %v", err)
			return nil, err
		}
		capsules = append(capsules, capsule)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	return capsules, nil
}

// UpdateCapsuleStatusByID updates the status of a capsule by its ID.
func UpdateCapsuleStatusByID(capsuleID int, newStatus string) error {
	query := `
	UPDATE capsules
	SET status = $1
	WHERE id = $2;
	`

	// Execute the status update
	_, err := database.DB.Exec(
		context.Background(),
		query,
		newStatus,
		capsuleID,
	)
	if err != nil {
		log.Printf("Error updating capsule status for capsule ID %d: %v", capsuleID, err)
		return err
	}

	return nil
}

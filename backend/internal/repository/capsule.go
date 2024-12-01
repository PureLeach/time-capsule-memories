package repository

import (
	"context"
	"log"
	"time"

	"time_capsule_memories/internal/database"
	"time_capsule_memories/internal/models"
)

func CreateCapsule(capsule *models.CreateCapsuleRequest) (*models.CapsuleResponse, error) {
	query := `
	INSERT INTO capsules (sender_name, send_at, message, recipient_email, recipient_tg_username, files_folder_UUID)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, sender_name, created_at, send_at, message, recipient_email, recipient_tg_username, files_folder_UUID, status;
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
		&createdCapsule.Status,
	)

	if err != nil {
		log.Printf("Error creating capsule in the database: %v", err)
		return nil, err
	}

	return createdCapsule, nil
}

func GetCapsulesByToday() ([]*models.CapsuleResponse, error) {
	currentDate := time.Now().Format("2006-01-02")

	query := `
	SELECT id, sender_name, created_at, send_at, message, recipient_email, recipient_tg_username, files_folder_UUID, status
	FROM capsules
	WHERE send_at::date = $1 AND status = 'waiting';
	`

	var capsules []*models.CapsuleResponse

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

	// Маппинг данных из строки результата в структуру
	for rows.Next() {
		capsule := &models.CapsuleResponse{}
		err := rows.Scan(
			&capsule.ID,
			&capsule.SenderName,
			&capsule.CreatedAt,
			&capsule.SendAt,
			&capsule.Message,
			&capsule.RecipientEmail,
			&capsule.RecipientTgUsername,
			&capsule.FilesFolderUUID,
		)
		if err != nil {
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

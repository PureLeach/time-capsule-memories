package repository

import (
	"context"
	"log"
	"time"

	"time_capsule_memories/internal/database"
	"time_capsule_memories/internal/models"
)

// CheckDatabaseConnection проверяет, доступна ли база данных
func CheckDatabaseConnection() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := database.DB.Ping(ctx); err != nil {
		log.Printf("Ошибка подключения к базе данных!!!: %v", err)
		return false
	}
	return true
}

// CreateCapsule сохраняет событие в базе данных и возвращает созданное событие с ID.
func CreateCapsule(capsule *models.CreateCapsule) (*models.CapsuleResponse, error) {
	query := `
	INSERT INTO capsules (sender_name, send_at, message, recipient_email, recipient_tg_username, files_folder_UUID)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, sender_name, created_at, send_at, message, recipient_email, recipient_tg_username, files_folder_UUID;
    `

	createdCapsule := &models.CapsuleResponse{}

	// Выполняем запрос и заполняем созданное событие
	err := database.DB.QueryRow(
		context.Background(),
		query,
		capsule.SenderName,
		// Устанавливаем значение времени отправки
		capsule.SendAt,
		capsule.Message,
		capsule.RecipientEmail,
		capsule.RecipientTgUsername,
		capsule.FilesFolderUUID, // Изменено на FilesFolderUUID
	).Scan(
		&createdCapsule.ID,
		&createdCapsule.SenderName,
		&createdCapsule.CreatedAt, // Заменено на CreatedAt
		&createdCapsule.SendAt,    // Заменено на SendAt
		&createdCapsule.Message,
		&createdCapsule.RecipientEmail,
		&createdCapsule.RecipientTgUsername,
		&createdCapsule.FilesFolderUUID, // Изменено на FilesFolderUUID
	)

	if err != nil {
		// Записываем лог ошибки с детальной информацией
		log.Printf("Error creating capsule: %v", err)
		return nil, err
	}

	return createdCapsule, nil // Возвращаем созданное событие с ID
}

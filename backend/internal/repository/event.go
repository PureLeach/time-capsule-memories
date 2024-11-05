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

// CreateEvent сохраняет событие в базе данных и возвращает созданное событие с ID.
func CreateEvent(event *models.CreateEvent) (*models.EventResponse, error) {
	query := `
        INSERT INTO events (sender_name, message_date, message, open_date, recipient_email, telegram_nick, photos)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, sender_name, message_date, message, open_date, recipient_email, telegram_nick, photos;
    `

	createdEvent := &models.EventResponse{}

	// Выполняем запрос и заполняем созданное событие
	err := database.DB.QueryRow(
		context.Background(),
		query,
		event.SenderName,
		event.MessageDate,
		event.Message,
		event.OpenDate,
		event.RecipientEmail,
		event.TelegramNick,
		event.Photos,
	).Scan(
		&createdEvent.ID,
		&createdEvent.SenderName,
		&createdEvent.MessageDate,
		&createdEvent.Message,
		&createdEvent.OpenDate,
		&createdEvent.RecipientEmail,
		&createdEvent.TelegramNick,
		&createdEvent.Photos,
	)

	if err != nil {
		// Записываем лог ошибки с детальной информацией
		log.Printf("Error creating event: %v", err)
		return nil, err
	}

	return createdEvent, nil // Возвращаем созданное событие с ID
}

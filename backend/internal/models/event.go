package models

import "time"

// CreateEvent представляет событие, которое будет создано.
type CreateEvent struct {
	SenderName     string    `json:"sender_name"`
	MessageDate    time.Time `json:"message_date" swaggertype:"string" example:"2024-11-05T15:04:05Z"` // Убедитесь, что формат времени соответствует
	Message        string    `json:"message"`
	OpenDate       time.Time `json:"open_date" swaggertype:"string" example:"2024-11-05T15:04:05Z"` // Убедитесь, что формат времени соответствует
	RecipientEmail string    `json:"recipient_email"`
	TelegramNick   string    `json:"telegram_nick"`
	Photos         string    `json:"photos"`
}

// Event представляет событие с ID.
type EventResponse struct {
	ID             int       `json:"id"` // Поле для ID события
	SenderName     string    `json:"sender_name"`
	MessageDate    time.Time `json:"message_date"`
	Message        string    `json:"message"`
	OpenDate       time.Time `json:"open_date"`
	RecipientEmail string    `json:"recipient_email"`
	TelegramNick   string    `json:"telegram_nick"`
	Photos         string    `json:"photos"`
}

// ErrorResponse представляет структуру для ответа с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

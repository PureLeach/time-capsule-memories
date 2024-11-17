package models

import "time"

// CreateCapsule представляет данные для создания события.
type CreateCapsule struct {
	SenderName          string  `json:"sender_name" example:"John Doe" validate:"required"`
	SendAt              string  `json:"send_at" swaggertype:"string" example:"2024-11-18" validate:"required,send_at_date_format,future_date"`
	Message             string  `json:"message" example:"Test Message" validate:"required"`
	RecipientEmail      string  `json:"recipient_email" example:"test@example.com" validate:"required,email"`
	RecipientTgUsername string  `json:"recipient_tg_username" example:"testuser" validate:"required"`
	FilesFolderUUID     *string `json:"files_folder_uuid,omitempty" example:"07023417-5079-429d-a113-cbef2ef164d7" validate:"omitempty,uuid4"`
}

// CapsuleResponse представляет событие с ID.
type CapsuleResponse struct {
	ID                  int       `json:"id"`
	SenderName          string    `json:"sender_name"`
	CreatedAt           time.Time `json:"created_at"`
	SendAt              time.Time `json:"send_at"`
	Message             string    `json:"message"`
	RecipientEmail      string    `json:"recipient_email"`
	RecipientTgUsername string    `json:"recipient_tg_username"`
	FilesFolderUUID     *string   `json:"files_folder_uuid"`
}

// ErrorResponse представляет структуру ответа с ошибкой.
type ErrorResponse struct {
	Error string `json:"error"`
}

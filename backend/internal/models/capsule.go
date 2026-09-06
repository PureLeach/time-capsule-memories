package models

import "time"

// CapsuleStatus mirrors the status_type enum in the capsules migration.
type CapsuleStatus string

const (
	StatusWaiting    CapsuleStatus = "waiting"
	StatusInProgress CapsuleStatus = "in progress"
	StatusDone       CapsuleStatus = "done"
)

type CreateCapsuleRequest struct {
	SenderName      string  `json:"sender_name" example:"John Doe" validate:"required,max=100"`
	SendAt          string  `json:"send_at" swaggertype:"string" example:"2099-11-18" validate:"required,send_at_date_format,future_date"`
	Message         string  `json:"message" example:"Test Message" validate:"required,max=4096"`
	RecipientEmail  string  `json:"recipient_email" example:"test@example.com" validate:"required,email,max=255"`
	FilesFolderUUID *string `json:"files_folder_uuid,omitempty" example:"07023417-5079-429d-a113-cbef2ef164d7" validate:"omitempty,uuid4"`
}

type CapsuleResponse struct {
	ID              int           `json:"id"`
	SenderName      string        `json:"sender_name"`
	CreatedAt       time.Time     `json:"created_at"`
	SendAt          time.Time     `json:"send_at"`
	Message         string        `json:"message"`
	RecipientEmail  string        `json:"recipient_email"`
	FilesFolderUUID *string       `json:"files_folder_uuid"`
	Status          CapsuleStatus `json:"status"`
}

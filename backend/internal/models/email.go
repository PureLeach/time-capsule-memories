package models

type EmailDataRequest struct {
	Subject         string  `json:"subject" example:"Test subject" validate:"required,max=255"`
	Body            string  `json:"body" example:"Test body" validate:"required,max=4096"`
	RecipientEmail  string  `json:"recipient_email" example:"test@example.com" validate:"required,email,max=255"`
	FilesFolderUUID *string `json:"files_folder_uuid,omitempty" example:"07023417-5079-429d-a113-cbef2ef164d7" validate:"omitempty,uuid4"`
}

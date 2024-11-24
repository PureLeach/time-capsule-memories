package models

type GeneratePresignedURLRequest struct {
	Directory string `json:"directory" query:"directory" validate:"required,uuid4"`
}

type PresignedURLResponse struct {
	PresignedURL string `json:"presigned_url"`
}

package models

type GeneratePresignedURLRequest struct {
	Directory string `json:"directory" query:"directory" validate:"required,uuid4"`
}

type PresignedURLResponse struct {
	PresignedURL string `json:"presigned_url"`
}

// Структура для хранения информации о файле
type FileObject struct {
	FileName    string
	Content     []byte
	ContentType string
}

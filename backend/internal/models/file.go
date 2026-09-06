package models

type GeneratePresignedURLRequest struct {
	Directory   string `json:"directory" query:"directory" validate:"required,uuid4"`
	ContentType string `json:"content_type" query:"content_type" validate:"required,image_content_type"`
}

// PresignedUpload is a signed multipart/form-data target. Fields are covered by
// the signature and must precede the file part.
type PresignedUpload struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}

type FileObject struct {
	FileName    string
	Content     []byte
	ContentType string
}

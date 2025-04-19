package validators

import (
	"fmt"
	"time_capsule_memories/internal/models"

	"github.com/google/uuid"
)

// ValidateGeneratePresignedURLRequest validates the request for generating a presigned URL.
func ValidateGeneratePresignedURLRequest(req *models.GeneratePresignedURLRequest) error {
	// Check if the directory parameter is provided
	if req.Directory == "" {
		return fmt.Errorf("directory parameter is required")
	}

	// Validate if the directory field contains a valid UUID
	if _, err := uuid.Parse(req.Directory); err != nil {
		return fmt.Errorf("invalid UUID format for directory")
	}

	return nil
}

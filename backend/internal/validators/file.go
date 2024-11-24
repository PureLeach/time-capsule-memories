package validators

import (
	"fmt"

	"time_capsule_memories/internal/models"

	"github.com/google/uuid"
)

func ValidateGeneratePresignedURLRequest(req *models.GeneratePresignedURLRequest) error {
	if req.Directory == "" {
		return fmt.Errorf("directory parameter is required")
	}

	if _, err := uuid.Parse(req.Directory); err != nil {
		return fmt.Errorf("invalid UUID format for directory")
	}

	return nil
}

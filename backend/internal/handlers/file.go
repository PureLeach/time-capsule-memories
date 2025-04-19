package handlers

import (
	"fmt"
	"net/http"
	"time"

	"time_capsule_memories/internal/minio_client"
	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/validators"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// @Summary Generate a presigned URL for file upload
// @Description Generates a presigned URL for uploading a file to a specific directory (UUID) in MinIO.
// @Tags file
// @Accept json
// @Produce json
// @Param directory query string true "Target directory UUID"
// @Success 200 {object} models.PresignedURLResponse "Presigned URL generated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 500 {object} models.ErrorResponse "Failed to generate presigned URL"
// @Router /generate-presigned-url [get]
func GeneratePresignedURLHandler(c echo.Context) error {
	directory := c.QueryParam("directory")
	req := models.GeneratePresignedURLRequest{
		Directory: directory,
	}

	// Validate input
	if err := validators.ValidateGeneratePresignedURLRequest(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Validation error: " + err.Error(),
		})
	}

	// Generate a unique file name
	fileName := uuid.New().String()

	// Generate a presigned upload URL valid for 1 hour
	presignedURL, err := minio_client.GeneratePresignedUploadURL(fmt.Sprintf("%s/%s", directory, fileName), time.Hour)
	if err != nil {
		c.Logger().Errorf("Failed to generate presigned URL: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Could not generate presigned URL",
		})
	}

	return c.JSON(http.StatusOK, models.PresignedURLResponse{
		PresignedURL: presignedURL,
	})
}

package handlers

import (
	"fmt"
	"net/http"
	"time"

	"time_capsule_memories/internal/logging"
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
func (h *Handler) GeneratePresignedURL(c echo.Context) error {
	directory := c.QueryParam("directory")
	req := models.GeneratePresignedURLRequest{
		Directory: directory,
	}

	if err := validators.ValidateGeneratePresignedURLRequest(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Validation error: " + err.Error(),
		})
	}

	fileName := uuid.New().String()

	ctx := c.Request().Context()
	url, err := h.store.GeneratePresignedUploadURL(ctx, fmt.Sprintf("%s/%s", directory, fileName), time.Hour)
	if err != nil {
		logging.FromContext(ctx).Error("failed to generate presigned URL",
			"directory", directory,
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Could not generate presigned URL",
		})
	}

	return c.JSON(http.StatusOK, models.PresignedURLResponse{
		PresignedURL: url,
	})
}

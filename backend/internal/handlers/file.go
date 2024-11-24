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
// @Description Generates a presigned URL for uploading a file to MinIO in a specific directory (UUID).
// @Tags file
// @Accept json
// @Produce json
// @Param directory query string true "UUID4 directory for file upload"
// @Success 200 {object} models.PresignedURLResponse "Presigned URL for file upload"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /generate-presigned-url [get]
func GeneratePresignedURLHandler(c echo.Context) error {
	// Заполняем модель параметрами запроса
	directory := c.QueryParam("directory")
	req := models.GeneratePresignedURLRequest{
		Directory: directory,
	}

	// Валидируем запрос через функцию валидации
	if err := validators.ValidateGeneratePresignedURLRequest(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	// Генерация случайного имени файла (UUID)
	fileName := (uuid.New()).String()

	// Генерация presigned URL для загрузки файла
	presignedURL, err := minio_client.GeneratePresignedUploadURL(directory+"/"+fileName, time.Hour)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: fmt.Sprintf("Failed to generate presigned URL: %v", err)})
	}

	// Возвращаем presigned URL
	return c.JSON(http.StatusOK, models.PresignedURLResponse{
		PresignedURL: presignedURL,
	})
}

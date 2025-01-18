package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time_capsule_memories/internal/minio_client"
	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/services"
	"time_capsule_memories/internal/validators"

	"github.com/labstack/echo/v4"
)

// @Summary Отправить тестовый email
// @Description Генерируем и отправляем тестовый email
// @Tags email
// @Accept json
// @Produce json
// @Param email body models.EmailDataRequest true "Данные для отправки по почте"
// @Success 204 {object} nil "Письмо успешно отправлено"
// @Failure 400 {object} models.ErrorResponse "Некорректные данные"
// @Failure 500 {object} models.ErrorResponse "Не удалось отправить письмо"
// @Router /send-test-email [post]
func SendTestEmail(c echo.Context) error {
	var emailData models.EmailDataRequest

	// Привязка данных из запроса
	if err := c.Bind(&emailData); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Некорректные данные: " + err.Error()})
	}

	// Валидация данных
	if err := validators.ValidateStruct(emailData); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})

	}

	// Получаем файлы из MinIO
	attachments, err := minio_client.GetFilesInDirectory(*emailData.FilesFolderUUID)
	if err != nil {
		log.Fatalf("Ошибка при получении файлов из каталога %s: %v", *emailData.FilesFolderUUID, err)
	}

	err = services.SendEmail(emailData.Subject, emailData.Body, emailData.RecipientEmail, attachments)
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	return c.JSON(http.StatusNoContent, nil)
}

package handlers

import (
	"net/http"
	"time_capsule_memories/internal/minio_client"
	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/services"
	"time_capsule_memories/internal/validators"

	"github.com/labstack/echo/v4"
)

// @Summary Send a test email
// @Description Generates and sends a test email
// @Tags email
// @Accept json
// @Produce json
// @Param email body models.EmailDataRequest true "Email payload"
// @Success 204 "Email sent successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid input data"
// @Failure 500 {object} models.ErrorResponse "Failed to send email"
// @Router /send-test-email [post]
func SendTestEmail(c echo.Context) error {
	var emailData models.EmailDataRequest

	// Bind request body to struct
	if err := c.Bind(&emailData); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request payload: " + err.Error(),
		})
	}

	// Validate input data
	if err := validators.ValidateStruct(emailData); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
	}

	// Fetch attachments from MinIO
	attachments, err := minio_client.GetFilesInDirectory(*emailData.FilesFolderUUID)
	if err != nil {
		c.Logger().Errorf("Failed to get files from directory %s: %v", *emailData.FilesFolderUUID, err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Could not retrieve attachments",
		})
	}

	// Send the email
	if err := services.SendEmail(emailData.Subject, emailData.Body, emailData.RecipientEmail, attachments); err != nil {
		c.Logger().Errorf("Failed to send email: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Could not send email",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

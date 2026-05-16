package handlers

import (
	"net/http"

	"time_capsule_memories/internal/logging"
	"time_capsule_memories/internal/models"
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
func (h *Handler) SendTestEmail(c echo.Context) error {
	var emailData models.EmailDataRequest

	if err := c.Bind(&emailData); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request payload: " + err.Error(),
		})
	}

	if err := validators.ValidateStruct(emailData); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
	}

	ctx := c.Request().Context()
	log := logging.FromContext(ctx)

	attachments, err := h.store.GetFilesInDirectory(ctx, *emailData.FilesFolderUUID)
	if err != nil {
		log.Error("failed to get files from directory",
			"folder_uuid", *emailData.FilesFolderUUID,
			"error", err,
		)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Could not retrieve attachments",
		})
	}

	if err := h.mailer.Send(ctx, emailData.Subject, emailData.Body, emailData.RecipientEmail, attachments); err != nil {
		log.Error("failed to send email", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Could not send email",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

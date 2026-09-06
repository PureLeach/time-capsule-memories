package handlers

import (
	"net/http"

	"time_capsule_memories/internal/logging"
	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/validators"

	"github.com/labstack/echo/v4"
)

// @Summary Send a test email
// @Description Development-only helper for verifying SMTP delivery end to end. Registered only when ENABLE_TEST_EMAIL_ENDPOINT is true, because it is unauthenticated and would otherwise let anyone send mail through the configured relay.
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
		return badRequest(c, "Invalid request payload")
	}

	if err := validators.ValidateStruct(emailData); err != nil {
		return badRequest(c, err.Error())
	}

	ctx := c.Request().Context()
	log := logging.FromContext(ctx)

	var attachments []models.FileObject
	if emailData.FilesFolderUUID != nil && *emailData.FilesFolderUUID != "" {
		files, err := h.store.GetFilesInDirectory(ctx, *emailData.FilesFolderUUID)
		if err != nil {
			log.Error("failed to get files from directory",
				"folder_uuid", *emailData.FilesFolderUUID,
				"error", err,
			)
			return internalError(c, "Could not retrieve attachments")
		}
		attachments = files
	}

	if err := h.mailer.Send(ctx, emailData.Subject, emailData.Body, emailData.RecipientEmail, attachments); err != nil {
		log.Error("failed to send test email", "error", err)
		return internalError(c, "Could not send email")
	}

	return c.NoContent(http.StatusNoContent)
}

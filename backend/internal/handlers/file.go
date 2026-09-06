package handlers

import (
	"net/http"
	"time"

	"time_capsule_memories/internal/logging"
	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/validators"

	"github.com/labstack/echo/v4"
)

// Short on purpose: the browser uses the form immediately, so a longer life
// only widens the window for replaying a leaked URL.
const uploadURLTTL = 5 * time.Minute

// @Summary Generate a presigned upload target for an attachment
// @Description Returns a signed multipart/form-data POST target scoped to one directory. The signed policy pins the content type and caps the upload at 5 MB, so the limits hold even if the client ignores them. Submit the returned fields verbatim, then the file part.
// @Tags file
// @Produce json
// @Param directory query string true "Target directory UUID"
// @Param content_type query string true "Image MIME type" Enums(image/jpeg, image/png, image/webp, image/gif)
// @Success 200 {object} models.PresignedUpload "Upload target generated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request"
// @Failure 500 {object} models.ErrorResponse "Failed to generate upload target"
// @Router /generate-presigned-url [get]
func (h *Handler) GeneratePresignedURL(c echo.Context) error {
	req := models.GeneratePresignedURLRequest{
		Directory:   c.QueryParam("directory"),
		ContentType: c.QueryParam("content_type"),
	}

	if err := validators.ValidateStruct(req); err != nil {
		return badRequest(c, err.Error())
	}

	ctx := c.Request().Context()
	upload, err := h.store.PresignUpload(ctx, req.Directory, req.ContentType, uploadURLTTL)
	if err != nil {
		logging.FromContext(ctx).Error("failed to presign upload",
			"directory", req.Directory,
			"error", err,
		)
		return internalError(c, "Could not generate upload target")
	}

	return c.JSON(http.StatusOK, upload)
}

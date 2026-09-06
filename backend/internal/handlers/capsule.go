package handlers

import (
	"net/http"

	"time_capsule_memories/internal/logging"
	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/validators"

	"github.com/labstack/echo/v4"
)

// @Summary Create a new capsule
// @Description Schedules a message, with optional attachments, for delivery on a future date.
// @Tags capsules
// @Accept json
// @Produce json
// @Param capsule body models.CreateCapsuleRequest true "Capsule creation payload"
// @Success 201 {object} models.CapsuleResponse "Capsule created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request data"
// @Failure 500 {object} models.ErrorResponse "Failed to create capsule"
// @Router /capsules [post]
func (h *Handler) CreateCapsule(c echo.Context) error {
	var capsule models.CreateCapsuleRequest

	if err := c.Bind(&capsule); err != nil {
		return badRequest(c, "Invalid request payload")
	}

	if err := validators.ValidateStruct(capsule); err != nil {
		return badRequest(c, err.Error())
	}

	ctx := c.Request().Context()
	created, err := h.capsuleRepo.Create(ctx, &capsule)
	if err != nil {
		logging.FromContext(ctx).Error("failed to create capsule", "error", err)
		return internalError(c, "Could not create capsule")
	}

	return c.JSON(http.StatusCreated, created)
}

func badRequest(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: message})
}

func internalError(c echo.Context, message string) error {
	return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: message})
}

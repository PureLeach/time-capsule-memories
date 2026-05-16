package handlers

import (
	"net/http"

	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/validators"

	"github.com/labstack/echo/v4"
)

// @Summary Create a new capsule
// @Description Creates a new time capsule with the given parameters
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
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request payload: " + err.Error(),
		})
	}

	if err := validators.ValidateCapsule(capsule); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
	}

	created, err := h.capsuleRepo.Create(c.Request().Context(), &capsule)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Could not create capsule",
		})
	}

	return c.JSON(http.StatusCreated, created)
}

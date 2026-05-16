package handlers

import (
	"net/http"

	"time_capsule_memories/internal/logging"
	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/validators"

	"github.com/labstack/echo/v4"
)

// @Summary Submit user feedback
// @Description Stores user feedback in the database
// @Tags feedback
// @Accept json
// @Produce json
// @Param feedback body models.CreateFeedbackRequest true "Feedback data"
// @Success 201 {object} models.FeedbackResponse "Feedback successfully created"
// @Failure 400 {object} models.ErrorResponse "Invalid input data"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /feedback [post]
func (h *Handler) CreateFeedback(c echo.Context) error {
	var feedback models.CreateFeedbackRequest

	if err := c.Bind(&feedback); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request payload: " + err.Error(),
		})
	}

	if err := validators.ValidateStruct(feedback); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Validation error: " + err.Error(),
		})
	}

	ctx := c.Request().Context()
	created, err := h.feedbackRepo.Create(ctx, &feedback)
	if err != nil {
		logging.FromContext(ctx).Error("failed to create feedback", "error", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create feedback",
		})
	}

	return c.JSON(http.StatusCreated, created)
}

package handlers

import (
	"net/http"

	"time_capsule_memories/internal/logging"
	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/validators"

	"github.com/labstack/echo/v4"
)

// @Summary Submit user feedback
// @Description Stores a free-form feedback message.
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
		return badRequest(c, "Invalid request payload")
	}

	if err := validators.ValidateStruct(feedback); err != nil {
		return badRequest(c, err.Error())
	}

	ctx := c.Request().Context()
	created, err := h.feedbackRepo.Create(ctx, &feedback)
	if err != nil {
		logging.FromContext(ctx).Error("failed to create feedback", "error", err)
		return internalError(c, "Could not save feedback")
	}

	return c.JSON(http.StatusCreated, created)
}

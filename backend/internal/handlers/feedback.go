package handlers

import (
	"net/http"
	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/repository"

	"github.com/go-playground/validator/v10"
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
func CreateFeedback(c echo.Context) error {
	var feedback models.CreateFeedbackRequest

	// Bind request payload
	if err := c.Bind(&feedback); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request payload: " + err.Error(),
		})
	}

	// Validate the input
	validate := validator.New()
	if err := validate.Struct(feedback); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Validation error: " + err.Error(),
		})
	}

	// Store feedback in the database
	createdFeedback, err := repository.CreateUserFeedback(&feedback)
	if err != nil {
		c.Logger().Errorf("Failed to create feedback: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to create feedback",
		})
	}

	return c.JSON(http.StatusCreated, createdFeedback)
}

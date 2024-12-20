package handlers

import (
	"net/http"
	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// @Summary Отправить отзыв
// @Description Создаёт запись отзыва пользователя
// @Tags feedback
// @Accept json
// @Produce json
// @Param feedback body models.CreateFeedbackRequest true "Данные для создания отзыва"
// @Success 201 {object} models.FeedbackResponse "Успешно создан отзыв"
// @Failure 400 {object} models.ErrorResponse "Некорректные данные"
// @Failure 500 {object} models.ErrorResponse "Не удалось создать событие"
// @Router /feedback [post]
func CreateFeedback(c echo.Context) error {
	var feedback models.CreateFeedbackRequest

	// Привязка данных из запроса
	if err := c.Bind(&feedback); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Некорректные данные: " + err.Error()})
	}

	// Валидация данных
	var validate = validator.New()
	if err := validate.Struct(feedback); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Ошибка валидации: " + err.Error()})
	}

	// Сохранение события в базе данных
	createdFeedback, err := repository.CreateFeedback(&feedback)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Не удалось создать событие"})
	}

	return c.JSON(http.StatusCreated, createdFeedback)
}

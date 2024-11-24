package handlers

import (
	"net/http"
	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/repository"
	"time_capsule_memories/internal/validators"

	"github.com/labstack/echo/v4"
)

// @Summary Создать новое событие
// @Description Создает новое событие с заданными параметрами
// @Tags capsules
// @Accept json
// @Produce json
// @Param capsule body models.CreateCapsuleRequest true "Данные для создания события"
// @Success 201 {object} models.CapsuleResponse "Успешно создано событие"
// @Failure 400 {object} models.ErrorResponse "Некорректные данные"
// @Failure 500 {object} models.ErrorResponse "Не удалось создать событие"
// @Router /capsules [post]
func CreateCapsule(c echo.Context) error {
	var capsule models.CreateCapsuleRequest

	// Привязка данных из запроса
	if err := c.Bind(&capsule); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Некорректные данные: " + err.Error()})
	}

	// Валидация данных
	if err := validators.ValidateCapsule(capsule); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	// Сохранение события в базе данных
	createdCapsule, err := repository.CreateCapsule(&capsule)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Не удалось создать событие"})
	}

	return c.JSON(http.StatusCreated, createdCapsule)
}

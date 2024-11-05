package handlers

import (
	"net/http"
	"time_capsule_memories/internal/models"
	"time_capsule_memories/internal/repository"

	"github.com/labstack/echo/v4"
)

// // @Summary Get events
// // @Description Get a list of events
// // @Tags events
// // @Accept json
// // @Produce json
// // @Success 200 {object} models.EventResponse
// // @Router /events [get]
// func GetEvents(c echo.Context) error {
// 	response := models.EventResponse{Hello: "world"}
// 	return c.JSON(http.StatusOK, response)
// }

// CreateEvent обрабатывает создание нового события
// @Summary Создать новое событие
// @Description Создает новое событие с заданными параметрами
// @Tags events
// @Accept json
// @Produce json
// @Param event body models.CreateEvent true "Данные для создания события"
// @Success 201 {object} models.EventResponse "Успешно создано событие"
// @Failure 400 {object} models.ErrorResponse "Некорректные данные"
// @Failure 500 {object} models.ErrorResponse "Не удалось создать событие"
// @Router /events [post]
func CreateEvent(c echo.Context) error {
	var event models.CreateEvent
	if err := c.Bind(&event); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	// Сохранение события в базе данных
	createdEvent, err := repository.CreateEvent(&event) // Получаем созданное событие
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Не удалось создать событие"})
	}

	return c.JSON(http.StatusCreated, createdEvent) // Используем ID созданного события
}

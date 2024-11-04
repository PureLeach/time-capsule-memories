package handlers

import (
    "net/http"
    "time_capsule_memories/internal/models"
    "github.com/labstack/echo/v4"
)


// @Summary Get events
// @Description Get a list of events
// @Tags events
// @Accept json
// @Produce json
// @Success 200 {object} models.EventResponse
// @Router /events [get]
func GetEvents(c echo.Context) error {
    response := models.EventResponse{Hello: "world"}
    return c.JSON(http.StatusOK, response)
}

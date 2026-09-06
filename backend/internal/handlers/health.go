package handlers

import (
	"context"
	"net/http"
	"time"

	"time_capsule_memories/internal/logging"

	"github.com/labstack/echo/v4"
)

const readinessTimeout = 2 * time.Second

// @Summary Liveness probe
// @Description Reports that the process is up. Does not touch dependencies.
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string "Service is alive"
// @Router /healthz [get]
func (h *Handler) Healthz(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// @Summary Readiness probe
// @Description Reports per-dependency health. Returns 503 when any check fails.
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string "All dependencies reachable"
// @Failure 503 {object} map[string]string "At least one dependency is down"
// @Router /readyz [get]
func (h *Handler) Readyz(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), readinessTimeout)
	defer cancel()
	log := logging.FromContext(c.Request().Context())

	checks := map[string]string{"database": "ok", "object_store": "ok"}
	ready := true

	if err := h.db.Ping(ctx); err != nil {
		log.Warn("readiness: database ping failed", "error", err)
		checks["database"] = "down"
		ready = false
	}

	if err := h.store.Ping(ctx); err != nil {
		log.Warn("readiness: object store ping failed", "error", err)
		checks["object_store"] = "down"
		ready = false
	}

	status := http.StatusOK
	if !ready {
		status = http.StatusServiceUnavailable
	}
	return c.JSON(status, checks)
}

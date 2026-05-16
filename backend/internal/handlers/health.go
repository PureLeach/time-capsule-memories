package handlers

import (
	"context"
	"net/http"
	"time"

	"time_capsule_memories/internal/logging"

	"github.com/labstack/echo/v4"
)

const readinessTimeout = 2 * time.Second

func (h *Handler) Healthz(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) Readyz(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), readinessTimeout)
	defer cancel()
	log := logging.FromContext(c.Request().Context())

	checks := map[string]string{"database": "ok", "minio": "ok"}
	ready := true

	if err := h.pool.Ping(ctx); err != nil {
		log.Warn("readiness: database ping failed", "error", err)
		checks["database"] = "down"
		ready = false
	}

	if err := h.store.Ping(ctx); err != nil {
		log.Warn("readiness: minio list buckets failed", "error", err)
		checks["minio"] = "down"
		ready = false
	}

	status := http.StatusOK
	if !ready {
		status = http.StatusServiceUnavailable
	}
	return c.JSON(status, checks)
}

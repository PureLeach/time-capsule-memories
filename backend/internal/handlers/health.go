package handlers

import (
	"context"
	"net/http"
	"time"

	"time_capsule_memories/internal/database"
	"time_capsule_memories/internal/logging"
	"time_capsule_memories/internal/minio_client"

	"github.com/labstack/echo/v4"
)

// readinessTimeout caps each backing check so an unhealthy dependency cannot
// hang the probe and, by extension, the orchestrator's polling loop.
const readinessTimeout = 2 * time.Second

// Healthz signals that the process is alive. Cheap on purpose — used by
// liveness probes that just want to know the binary is still running.
func Healthz(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// Readyz checks the downstream dependencies the app needs to handle traffic:
// the Postgres pool and MinIO. A 503 means the orchestrator should keep this
// instance out of rotation until it recovers.
func Readyz(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), readinessTimeout)
	defer cancel()
	log := logging.FromContext(c.Request().Context())

	checks := map[string]string{"database": "ok", "minio": "ok"}
	ready := true

	if database.DB == nil {
		checks["database"] = "uninitialised"
		ready = false
	} else if err := database.DB.Ping(ctx); err != nil {
		log.Warn("readiness: database ping failed", "error", err)
		checks["database"] = "down"
		ready = false
	}

	minioClient, err := minio_client.GetMinioClient()
	if err != nil {
		log.Warn("readiness: minio client unavailable", "error", err)
		checks["minio"] = "down"
		ready = false
	} else if _, err := minioClient.ListBuckets(ctx); err != nil {
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

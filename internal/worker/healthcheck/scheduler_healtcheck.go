package workerhealthcheck

import (
	"context"
	domainhealthcheck "go-bootstrap/internal/domain/healthcheck"
	"log/slog"
	"time"
)

type SchedulerHealthCheck struct {
	healthcheckService domainhealthcheck.HealthCheckService
}

func NewSchedulerHealthCheck(
	healthcheckService domainhealthcheck.HealthCheckService,
) *SchedulerHealthCheck {
	return &SchedulerHealthCheck{
		healthcheckService: healthcheckService,
	}
}

func (w *SchedulerHealthCheck) CheckDependencies() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("Starting health check...")

	outputHealthcheck := w.healthcheckService.CheckDependencies(ctx)

	// Log overall health status
	switch outputHealthcheck.Status {
	case domainhealthcheck.StatusHealthCheckDegraded:
		slog.Warn("Health check status: DEGRADED",
			"status", outputHealthcheck.Status,
			"timestamp", outputHealthcheck.Timestamp,
		)
	case domainhealthcheck.StatusHealthCheckUnhealthy:
		slog.Error("Health check status: UNHEALTHY",
			"status", outputHealthcheck.Status,
			"timestamp", outputHealthcheck.Timestamp,
		)
	case domainhealthcheck.StatusHealthCheckHealthy:
		slog.Info("Health check status: HEALTHY",
			"status", outputHealthcheck.Status,
			"timestamp", outputHealthcheck.Timestamp,
		)
	}

	// Log database dependency status
	switch outputHealthcheck.Database.Status {
	case domainhealthcheck.StatusDependencyError:
		slog.Error("[SCHEDULER] Database dependency check failed",
			"status", outputHealthcheck.Database.Status,
			"response_time", outputHealthcheck.Database.ResponseTime,
			"message", outputHealthcheck.Database.Message,
		)
	case domainhealthcheck.StatusDependencyOk:
		slog.Info("[SCHEDULER] Database dependency check passed",
			"status", outputHealthcheck.Database.Status,
			"response_time", outputHealthcheck.Database.ResponseTime,
			"message", outputHealthcheck.Database.Message,
		)
	}
}

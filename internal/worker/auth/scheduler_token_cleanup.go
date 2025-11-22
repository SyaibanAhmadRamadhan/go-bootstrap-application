package workerauth

import (
	"context"
	domainauth "go-bootstrap/internal/domain/auth"
	"log/slog"
	"time"
)

type SchedulerTokenCleanup struct {
	authService domainauth.AuthService
}

func NewSchedulerTokenCleanup(
	authService domainauth.AuthService,
) *SchedulerTokenCleanup {
	return &SchedulerTokenCleanup{
		authService: authService,
	}
}

func (w *SchedulerTokenCleanup) CleanupExpiredTokens() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	slog.Info("Starting expired tokens cleanup...")

	w.authService.WorkerDeleteExpiredTokens(ctx)
}

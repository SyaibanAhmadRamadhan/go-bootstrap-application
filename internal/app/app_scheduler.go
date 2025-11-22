package app

import (
	"context"
	"errors"
	"go-bootstrap/internal/config"
	"go-bootstrap/internal/infrastructure"
	healthcheckrepository "go-bootstrap/internal/module/healthcheck/repository"
	healthcheckservice "go-bootstrap/internal/module/healthcheck/service"
	workerhealthcheck "go-bootstrap/internal/worker/healthcheck"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

type schedulerApp struct {
	cron    *cron.Cron
	closeFn []func() error
}

func NewSchedulerApp() *schedulerApp {
	schedulerApp := &schedulerApp{
		cron:    cron.New(cron.WithSeconds()),
		closeFn: make([]func() error, 0),
	}

	schedulerApp.init()

	return schedulerApp
}

func (s *schedulerApp) Start() {
	s.cron.Start()
	slog.Info("Running scheduler...")
}

func (s *schedulerApp) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down scheduler...")

	stopCtx := s.cron.Stop()

	select {
	case <-stopCtx.Done():
		slog.Info("All cron jobs completed gracefully")
	case <-ctx.Done():
		slog.Warn("Shutdown timeout reached, forcing stop")
	}

	errs := make([]error, 0, len(s.closeFn))
	for _, fn := range s.closeFn {
		if err := fn(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	slog.Info("Scheduler shut down successfully")
	return nil
}

func (s *schedulerApp) init() {
	db, err := infrastructure.NewDB()
	if err != nil {
		panic(err)
	}
	s.closeFn = append(s.closeFn, db.Close)

	healthcheckRepo := healthcheckrepository.NewRepository(db)
	healthcheckService := healthcheckservice.NewService(healthcheckRepo)
	healthcheckWorker := workerhealthcheck.NewSchedulerHealthCheck(healthcheckService)

	s.registerCronJobs(healthcheckWorker)
}

func (s *schedulerApp) registerCronJobs(
	healthcheckWorker *workerhealthcheck.SchedulerHealthCheck,
) {
	schedulerConfig := config.GetAppScheduler()

	_, err := s.cron.AddFunc(schedulerConfig.HealthCheckInterval, func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Panic recovered in CheckDependencies", "panic", r)
			}
		}()
		healthcheckWorker.CheckDependencies()
	})
	if err != nil {
		slog.Error("Failed to register CheckDependencies", "error", err)
	} else {
		slog.Info("Registered CheckDependencies", "schedule", "every 5 minutes")
	}
}

// WaitForNextRun blocks until the next scheduled job runs
// Useful for testing or ensuring at least one job cycle completes
func (s *schedulerApp) WaitForNextRun(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	entries := s.cron.Entries()
	if len(entries) == 0 {
		return false
	}

	nextRun := entries[0].Next
	for _, entry := range entries {
		if entry.Next.After(nextRun) {
			nextRun = entry.Next
		}
	}

	waitDuration := time.Until(nextRun)
	if waitDuration < 0 {
		return false
	}

	select {
	case <-time.After(waitDuration + time.Second):
		return true
	case <-timer.C:
		return false
	}
}

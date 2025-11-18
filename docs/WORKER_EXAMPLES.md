# Worker Layer Examples

This document provides detailed examples for implementing workers, schedulers, and background jobs.

## Table of Contents

- [Worker Layer Examples](#worker-layer-examples)
  - [Table of Contents](#table-of-contents)
  - [Scheduler (Cron Jobs)](#scheduler-cron-jobs)
    - [Basic Example](#basic-example)
  - [How Scheduler Works](#how-scheduler-works)
    - [Scheduler App Implementation](#scheduler-app-implementation)
  - [Cron Expression Format](#cron-expression-format)
    - [Special Characters](#special-characters)
  - [Common Cron Expressions](#common-cron-expressions)
  - [Complete Example](#complete-example)
    - [Worker Implementation](#worker-implementation)
    - [Scheduler App with Multiple Jobs](#scheduler-app-with-multiple-jobs)
    - [Configuration Example](#configuration-example)
  - [Best Practices](#best-practices)

## Scheduler (Cron Jobs)

Schedulers handle time-triggered background tasks using cron expressions.

### Basic Example

```go
package workerhealthcheck

import (
    "context"
    "log/slog"
    "time"
    domainhealthcheck "project/internal/domain/healthcheck"
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

// CheckDependencies runs periodic health checks
// Called by cron scheduler based on healthcheck_interval in config
func (w *SchedulerHealthCheck) CheckDependencies() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    slog.Info("Starting health check...")

    // Call domain service
    output := w.healthcheckService.CheckDependencies(ctx)

    // Log results with structured logging
    switch output.Status {
    case domainhealthcheck.StatusHealthCheckHealthy:
        slog.Info("Health check status: HEALTHY",
            "status", output.Status,
            "timestamp", output.Timestamp,
        )
    case domainhealthcheck.StatusHealthCheckDegraded:
        slog.Warn("Health check status: DEGRADED",
            "status", output.Status,
            "timestamp", output.Timestamp,
        )
    case domainhealthcheck.StatusHealthCheckUnhealthy:
        slog.Error("Health check status: UNHEALTHY",
            "status", output.Status,
            "timestamp", output.Timestamp,
        )
    }

    // Log database dependency status
    switch output.Database.Status {
    case domainhealthcheck.StatusDependencyHealthy:
        slog.Info("Database status: HEALTHY",
            "response_time", output.Database.ResponseTime,
        )
    case domainhealthcheck.StatusDependencyUnhealthy:
        slog.Error("Database status: UNHEALTHY",
            "error", output.Database.Message,
            "response_time", output.Database.ResponseTime,
        )
    }
}
```

## How Scheduler Works

The scheduler is configured in `internal/app/scheduler.go` and manages all cron jobs.

### Scheduler App Implementation

```go
package app

import (
    "context"
    "erp-directory-service/internal/config"
    workerhealthcheck "erp-directory-service/internal/worker/healthcheck"
    "log/slog"
    "time"

    "github.com/robfig/cron/v3"
)

type schedulerApp struct {
    cron    *cron.Cron
    closeFn []func() error
}

func NewSchedulerApp() *schedulerApp {
    // Create cron with seconds support
    c := cron.New(cron.WithSeconds())
    
    // Get configuration
    appCfg := config.GetAppScheduler()
    
    // Initialize database and services (dependency injection)
    db, err := infrastructure.NewDB()
    if err != nil {
        slog.Error("Failed to initialize database", "error", err)
        panic(err)
    }
    
    // Create repositories
    healthcheckRepo := healthcheckrepository.NewRepository(db)
    
    // Create services
    healthcheckService := healthcheckservice.NewService(healthcheckRepo)
    
    // Initialize workers
    healthcheckWorker := workerhealthcheck.NewSchedulerHealthCheck(
        healthcheckService,
    )
    
    // Register cron jobs with panic recovery
    _, err := c.AddFunc(appCfg.HealthCheckInterval, func() {
        defer func() {
            if r := recover(); r != nil {
                slog.Error("Panic recovered in health check job", "panic", r)
            }
        }()
        healthcheckWorker.CheckDependencies()
    })
    if err != nil {
        slog.Error("Failed to register health check job", "error", err)
    }
    
    return &schedulerApp{
        cron:    c,
        closeFn: []func() error{},
    }
}

func (s *schedulerApp) Start() {
    slog.Info("Starting scheduler...")
    s.cron.Start()
    slog.Info("Scheduler started successfully")
}

func (s *schedulerApp) Stop(ctx context.Context) error {
    slog.Info("Stopping scheduler...")
    
    // Stop accepting new jobs
    stopCtx := s.cron.Stop()
    
    // Wait for running jobs to complete or timeout
    select {
    case <-stopCtx.Done():
        slog.Info("All jobs completed")
    case <-ctx.Done():
        slog.Warn("Scheduler shutdown timed out, some jobs may not have completed")
    }
    
    // Close other resources
    for _, fn := range s.closeFn {
        if err := fn(); err != nil {
            slog.Error("Error closing resource", "error", err)
        }
    }
    
    slog.Info("Scheduler stopped successfully")
    return nil
}
```

## Cron Expression Format

The scheduler uses `robfig/cron/v3` which supports seconds:

```text
┌─────────────── second (0 - 59)
│ ┌───────────── minute (0 - 59)
│ │ ┌─────────── hour (0 - 23)
│ │ │ ┌───────── day of month (1 - 31)
│ │ │ │ ┌─────── month (1 - 12)
│ │ │ │ │ ┌───── day of week (0 - 6) (Sunday to Saturday)
│ │ │ │ │ │
│ │ │ │ │ │
* * * * * *
```

### Special Characters

- `*` - Any value
- `,` - Value list separator (e.g., `1,3,5`)
- `-` - Range of values (e.g., `1-5`)
- `/` - Step values (e.g., `*/5` = every 5)
- `?` - No specific value (day of month/week)

## Common Cron Expressions

```text
# Every 5 minutes
0 */5 * * * *

# Daily at 2 AM
0 0 2 * * *

# Every 6 hours
0 0 */6 * * *

# 9:30 AM on weekdays (Monday-Friday)
0 30 9 * * 1-5

# Every 10 seconds
*/10 * * * * *

# First day of every month at midnight
0 0 0 1 * *

# Every Sunday at 3 AM
0 0 3 * * 0

# Every hour on the hour
0 0 * * * *

# Every 15 minutes
0 */15 * * * *

# Twice a day (9 AM and 5 PM)
0 0 9,17 * * *
```

## Complete Example

Here's a complete example showing multiple cron jobs:

### Worker Implementation

```go
package workeruser

import (
    "context"
    "log/slog"
    "time"
    domainuser "project/internal/domain/user"
)

type SchedulerUser struct {
    userService domainuser.UserService
}

func NewSchedulerUser(
    userService domainuser.UserService,
) *SchedulerUser {
    return &SchedulerUser{
        userService: userService,
    }
}

// SendDailyReport sends daily user activity reports
// Cron: 0 0 9 * * * (Daily at 9 AM)
func (w *SchedulerUser) SendDailyReport() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    slog.Info("Starting daily report generation...")

    err := w.userService.GenerateDailyReport(ctx)
    if err != nil {
        slog.Error("Failed to generate daily report", "error", err)
        return
    }

    slog.Info("Daily report generated successfully")
}

// CleanupInactiveUsers removes users inactive for 90+ days
// Cron: 0 0 2 * * 0 (Every Sunday at 2 AM)
func (w *SchedulerUser) CleanupInactiveUsers() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
    defer cancel()

    slog.Info("Starting inactive user cleanup...")

    count, err := w.userService.CleanupInactiveUsers(ctx, 90*24*time.Hour)
    if err != nil {
        slog.Error("Failed to cleanup inactive users", "error", err)
        return
    }

    slog.Info("Inactive user cleanup completed", "deleted_count", count)
}

// SyncWithExternalSystem syncs user data with external CRM
// Cron: 0 */30 * * * * (Every 30 minutes)
func (w *SchedulerUser) SyncWithExternalSystem() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()

    slog.Info("Starting user sync with external system...")

    err := w.userService.SyncWithExternalCRM(ctx)
    if err != nil {
        slog.Error("Failed to sync with external system", "error", err)
        return
    }

    slog.Info("User sync completed successfully")
}
```

### Scheduler App with Multiple Jobs

```go
func NewSchedulerApp() *schedulerApp {
    c := cron.New(cron.WithSeconds())
    appCfg := config.GetAppScheduler()
    
    // Initialize services
    healthcheckService := initHealthCheckService()
    userService := initUserService()
    
    // Initialize workers
    healthcheckWorker := workerhealthcheck.NewSchedulerHealthCheck(healthcheckService)
    userWorker := workeruser.NewSchedulerUser(userService)
    
    // Register health check job (from config)
    c.AddFunc(appCfg.HealthCheckInterval, wrapWithRecovery(
        "health_check",
        healthcheckWorker.CheckDependencies,
    ))
    
    // Register user jobs
    c.AddFunc("0 0 9 * * *", wrapWithRecovery(
        "daily_report",
        userWorker.SendDailyReport,
    ))
    
    c.AddFunc("0 0 2 * * 0", wrapWithRecovery(
        "cleanup_inactive_users",
        userWorker.CleanupInactiveUsers,
    ))
    
    c.AddFunc("0 */30 * * * *", wrapWithRecovery(
        "sync_external_system",
        userWorker.SyncWithExternalSystem,
    ))
    
    return &schedulerApp{cron: c}
}

// wrapWithRecovery wraps a job function with panic recovery
func wrapWithRecovery(jobName string, fn func()) func() {
    return func() {
        defer func() {
            if r := recover(); r != nil {
                slog.Error("Panic recovered in job",
                    "job", jobName,
                    "panic", r,
                )
            }
        }()
        
        slog.Info("Job started", "job", jobName)
        start := time.Now()
        
        fn()
        
        duration := time.Since(start)
        slog.Info("Job completed",
            "job", jobName,
            "duration", duration,
        )
    }
}
```

### Configuration Example

```json
{
    "app_scheduler": {
        "name": "my-service-scheduler",
        "env": "production",
        "debug_mode": false,
        "healthcheck_interval": "0 */5 * * * *",
        "pprof": {
            "enable": false,
            "port": 6060,
            "static_token": "secret"
        }
    }
}
```

## Best Practices

1. **Always use context with timeout** - Prevent jobs from running indefinitely
2. **Use structured logging** - Log job start, completion, and errors with context
3. **Implement panic recovery** - Wrap jobs with panic recovery to prevent crashes
4. **Make jobs idempotent** - Jobs should be safe to run multiple times
5. **Monitor job execution** - Log duration and success/failure metrics
6. **Handle graceful shutdown** - Wait for running jobs to complete before exiting
7. **Use appropriate timeouts** - Short jobs (seconds), long jobs (minutes)
8. **Avoid blocking operations** - Use goroutines for parallel work when needed

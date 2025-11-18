# GO SERVICE BOOTSTRAP

A starter template for building internal backend services.  
This repository provides a ready-to-use project layout, wiring, and tooling so you don't have to set up everything from scratch for every new service.

It's designed for a "distributed monolith" style architecture: many services, one database, with clear ownership and strict API contracts between services.

## Table of Contents

- [GO SERVICE BOOTSTRAP](#go-service-bootstrap)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Makefile Commands](#makefile-commands)
    - [Running the Application](#running-the-application)
    - [Code Generation](#code-generation)
    - [Utilities](#utilities)
  - [Configuration](#configuration)
    - [Configuration Structure](#configuration-structure)
    - [Application Configs](#application-configs)
    - [Pprof Configuration (Realtime Hot-Reload)](#pprof-configuration-realtime-hot-reload)
    - [How Configuration Works Internally](#how-configuration-works-internally)
  - [Project Structure](#project-structure)
    - [Domain Layer (`internal/domain/`)](#domain-layer-internaldomain)
    - [Module Layer (`internal/module/`)](#module-layer-internalmodule)
    - [Transport Layer (`internal/transport/`)](#transport-layer-internaltransport)
    - [Worker Layer (`internal/worker/`)](#worker-layer-internalworker)
    - [Other Important Directories](#other-important-directories)
  - [Architecture Flow](#architecture-flow)
  - [Development Workflow](#development-workflow)
  - [Naming Conventions](#naming-conventions)
    - [Folder vs Package Names](#folder-vs-package-names)
    - [Package Names Rules](#package-names-rules)
    - [Interface Names](#interface-names)
    - [Struct Names](#struct-names)
    - [File Names](#file-names)
    - [Transport Layer Naming](#transport-layer-naming)
    - [Worker Layer Naming](#worker-layer-naming)
    - [Method Names](#method-names)
    - [Constants and Enums](#constants-and-enums)
    - [Function Parameters and Return Values](#function-parameters-and-return-values)
    - [Constructor Functions](#constructor-functions)
  - [Testing](#testing)
  - [Performance \& Debugging](#performance--debugging)
    - [Memory Leak Detection](#memory-leak-detection)

## Prerequisites

Before you start, make sure you have the following installed:

- Go 1.24 or higher
- Node.js and npm (for OpenAPI tooling)
- Make

## Installation

Run the installation script to install all required dependencies:

```bash
bash scripts/install_dependency.sh
```

This script will install:

- **buf** - Protocol buffer management tool
- **oapi-codegen** - OpenAPI code generator for Go
- **mockgen** - Mock generator for Go interfaces
- **gosec** - Go security checker
- **openapi-format** - OpenAPI formatter and validator
- **@redocly/cli** - OpenAPI documentation preview tool

## Makefile Commands

The project includes a `Makefile` with several useful commands:

### Running the Application

- `make run-restapi` - Start the REST API server with stdout logging and file writer
- `make run-grpcapi` - Start the gRPC API server with stdout logging and file writer
- `make run-scheduler` - Start the Scheduler with stdout logging and file writer

### Code Generation

- `make generate` - Run all code generation (API, gRPC, and mocks)
- `make api_generate` - Generate REST API server and types from OpenAPI spec
- `make buf_generate` - Generate gRPC code from Protocol Buffers
- `make go_generate` - Run Go's built-in code generation (e.g., mockgen)

### Utilities

- `make preview_open_api` - Preview OpenAPI documentation in browser using Redocly
- `make clean` - Clean generated files, builds, and artifacts

## Configuration

### Configuration Structure

The application uses JSON-based configuration files (`env.json`) that can be hot-reloaded during development. The configuration is separated by application type for better modularity and independence.

**Config file location:** `env.json` (create from `env.example.json`)

**Important:** After cloning the repository, run `make generate` to generate all required code before starting the application.

### Application Configs

Each application (HTTP API, gRPC API, Scheduler) has its own independent configuration:

**`app_rest_api` - REST API Configuration:**

```json
{
    "app_rest_api": {
        "name": "directory-service-rest-api",   // Service name for logging/tracing
        "env": "development",                    // Environment: development/staging/production
        "debug_mode": true,                      // Enable debug logging and SQL query logging
        "port": 8080,                            // HTTP server port
        "pprof": {                               // Pprof configuration (nested)
            "enable": true,                      // Enable/disable pprof server
            "port": 8080,                        // Pprof HTTP server port
            "static_token": "your-secret-token"  // Static token for authentication
        }
    }
}
```

**`app_grpc_api` - gRPC API Configuration:**

```json
{
    "app_grpc_api": {
        "name": "directory-service-grpc-api",   // Service name for logging/tracing
        "env": "development",                    // Environment: development/staging/production
        "debug_mode": true,                      // Enable debug logging and SQL query logging
        "port": 9090,                            // gRPC server port
        "pprof": {                               // Pprof configuration (nested)
            "enable": false,                     // Enable/disable pprof server
            "port": 7070,                        // Pprof HTTP server port
            "static_token": "your-secret-token"  // Static token for authentication
        }
    }
}
```

**`app_scheduler` - Background Jobs Configuration:**

```json
{
    "app_scheduler": {
        "name": "directory-service-scheduler",     // Service name for logging/tracing
        "env": "development",                       // Environment: development/staging/production
        "debug_mode": true,                         // Enable debug logging
        "healthcheck_interval": "0 */5 * * * *",   // Cron expression (every 5 minutes)
        "pprof": {                                  // Pprof configuration (nested)
            "enable": true,                         // Enable/disable pprof server
            "port": 7070,                           // Pprof HTTP server port
            "static_token": "your-secret-token"     // Static token for authentication
        }
    }
}
```

**Why Separate Configs?**

- Each service can have different names for logging/monitoring
- Independent debug modes per service
- Different environment settings (e.g., HTTP in production, gRPC in staging)
- Each service has its own pprof configuration (can enable/disable independently)
- Better multi-service architecture support

**`database` - Database Configuration:**

```json
{
    "database": {
        "dsn": "user:password@tcp(host:port)/dbname?parseTime=true",
        "max_open_conns": 25,        // Maximum open connections
        "max_idle_conns": 25,        // Maximum idle connections
        "conn_max_lifetime": "300s", // Connection max lifetime
        "conn_max_idle_time": "60s"  // Connection max idle time
    }
}
```

### Pprof Configuration (Realtime Hot-Reload)

Each application (REST API, gRPC API, Scheduler) has its own **independent pprof configuration** nested within its config. This allows you to enable/disable profiling per service.

**Pprof is nested inside each app config:**

```json
{
    "app_rest_api": {
        "pprof": {
            "enable": true,
            "port": 8080,
            "static_token": "your-secret-token"
        }
    },
    "app_grpc_api": {
        "pprof": {
            "enable": false,
            "port": 7070,
            "static_token": "your-secret-token"
        }
    },
    "app_scheduler": {
        "pprof": {
            "enable": true,
            "port": 6060,
            "static_token": "your-secret-token"
        }
    }
}
```

**Hot-Reload Behavior:**

The pprof server monitors configuration changes and automatically starts/stops based on the `enable` flag in each app's config:

1. **Enable pprof:**
   - Set `"enable": true` in the respective app's pprof config
   - Server starts automatically on the configured port
   - No application restart required

2. **Disable pprof:**
   - Set `"enable": false` in the respective app's pprof config
   - Server stops gracefully
   - No application restart required

3. **Change pprof port:**
   - **Step 1:** Set `"enable": false` (stops current server)
   - **Step 2:** Change `"port": 8080` to desired port
   - **Step 3:** Set `"enable": true` (starts server on new port)

   **Important:** The restart trigger is based on changes to the `enable` flag. Direct port changes without toggling `enable` will not take effect.

**Accessing Pprof Endpoints:**

```bash
# REST API pprof (port 8080)
curl http://localhost:8080/debug/pprof/heap

# gRPC API pprof (port 7070)
curl http://localhost:7070/debug/pprof/goroutine

# Scheduler pprof (port 6060)
open http://localhost:6060/debug/pprof/

# Available endpoints:
# /debug/pprof/          - Index page
# /debug/pprof/heap     - Memory heap profile
# /debug/pprof/goroutine - Goroutine profile
# /debug/pprof/profile  - CPU profile
# /debug/pprof/trace    - Execution trace
# /debug/pprof/allocs   - Memory allocations
# /debug/pprof/block    - Blocking profile
# /debug/pprof/mutex    - Mutex contention
```

**Security Note:** In production, always use a strong `static_token` and restrict access to pprof endpoints via network policies or authentication middleware.

**Why Nested Pprof Config?**

- Each service can enable/disable profiling independently
- Different services can use different pprof ports
- No port conflicts when running multiple services simultaneously
- Better isolation and control per service

### How Configuration Works Internally

**For beginners:** Understanding how configuration flows through the application:

1. **Configuration Files (`env.json`):**
   - Contains all application settings
   - Watched for changes during development
   - Automatically reloaded when modified

2. **Config Loading (`internal/config/load_config.go`):**

   ```go
   // In cmd layer - Application startup (cmd.go PersistentPreRun)
   config.LoadConfig("app_rest_api.debug_mode")  // Initialize with hot-reload key
   appCfg := config.GetAppRestApi()               // Get REST API config
   ```

   **Hot-Reload Key:** The parameter `"app_rest_api.debug_mode"` tells the config loader to watch for changes to this specific key. When `debug_mode` changes, the config automatically reloads.

3. **Config Types (`internal/config/type.go`):**
   - Defines struct types for all configurations
   - Uses struct tags for JSON mapping: `env:"field_name"`

4. **Usage in Application:**

   ```go
   // In cmd layer - Application startup
   appCfg := config.GetAppRestApi()
   provider.NewLogging(filename, slogHook, zerologHook, 
                       appCfg.DebugMode, appCfg.Env, appCfg.Name)
   
   // In app layer - Feature initialization
   appCfg := config.GetAppRestApi()
   db := provider.NewDB(appCfg.DebugMode)
   ```

**Configuration Flow:**

```sh
env.json → LoadConfig() → root struct → Getter functions → Application
```

**Getter Functions:**

- `config.GetAppRestApi()` - Get REST API config (includes nested pprof)
- `config.GetAppGrpcApi()` - Get gRPC API config (includes nested pprof)
- `config.GetAppScheduler()` - Get Scheduler config (includes nested pprof)
- `config.GetPprofAppRestApi()` - Get pprof config for REST API only
- `config.GetPprofAppGrpcApi()` - Get pprof config for gRPC API only
- `config.GetPprofAppScheduler()` - Get pprof config for Scheduler only
- `config.GetDatabase()` - Get database config
- `config.UnwatchLoader()` - Stop watching config file (called on shutdown)

**Why This Design?**

- **Separation of Concerns:** Config loading separated from business logic
- **Type Safety:** Strongly typed configuration access
- **Hot-Reload:** Changes detected automatically in development
- **Testability:** Easy to mock config in tests
- **Independence:** Each app can use different configs

## Project Structure

### Domain Layer (`internal/domain/`)

The domain layer contains business logic interfaces and data structures. Each domain represents a business capability:

```text
internal/domain/<feature>/
├── dto.go              # Data Transfer Objects
├── repository.go       # Repository interface definitions
├── service.go          # Service interface definitions
└── value_object.go     # Domain value objects and enums
```

**Key principles:**

- Defines **interfaces only**, no implementation
- Contains DTOs, value objects, and business rules
- Framework-agnostic and pure Go code
- Generates mock interfaces for testing using `go:generate` directives

**Example:**

```go
type HealthCheckService interface {
    CheckDependencies(ctx context.Context) (output CheckDependenciesOutput)
}
```

### Module Layer (`internal/module/`)

The module layer contains concrete implementations of domain interfaces. This is where the actual business logic lives:

```text
internal/module/<feature>/
├── repository/
│   ├── repo_<feature>.go              # Repository interface wrapper
│   ├── repo_<feature>_datastore.go    # Database implementation
│   └── repo_<feature>_datastore_test.go
└── service/
    ├── service_<feature>.go           # Service implementation
    └── service_<feature>_test.go
```

**Key principles:**

- Implements domain interfaces
- Contains actual business logic
- Handles data persistence and external dependencies
- Separated into `repository` (data access) and `service` (business logic)
- Each implementation can have multiple variants (e.g., datastore, cache, external API)

**Example:**

```go
type service struct {
    healthcheckRepo domainhealthcheck.HealthCheckRepositoryDatastore
}

func (s *service) CheckDependencies(ctx context.Context) (output domainhealthcheck.CheckDependenciesOutput) {
    // Business logic implementation
}
```

### Transport Layer (`internal/transport/`)

The transport layer acts as the adapter between external protocols (HTTP REST, gRPC) and your domain/business logic. It handles request/response transformation and protocol-specific concerns:

```text
internal/transport/<feature>/
├── grpc_<feature>.go      # gRPC handlers
└── restapi_<feature>.go   # REST API handlers
```

**Key principles:**

- Handles protocol-specific request/response formatting
- Transforms between protocol types (protobuf, HTTP) and domain DTOs
- No business logic - only marshalling/unmarshalling and calling services
- Each transport handler depends on domain service interfaces
- Thin layer that delegates to the domain/module layer

**gRPC Transport Example:**

```go
package transporthealthcheck

import (
    "context"
    domainhealthcheck "project/internal/domain/healthcheck"
    "project/internal/gen/grpcgen/healthcheck"
)

type TransportHealthCheckGrpc struct {
    healthcheckService domainhealthcheck.HealthCheckService
    healthcheck.UnimplementedHealthCheckServiceServer
}

func NewTransportGrpc(
    healthcheckService domainhealthcheck.HealthCheckService,
) *TransportHealthCheckGrpc {
    return &TransportHealthCheckGrpc{
        healthcheckService: healthcheckService,
    }
}

func (t *TransportHealthCheckGrpc) ApiV1HealthCheck(
    ctx context.Context, 
    req *healthcheck.ApiV1HealthCheckRequest,
) (*healthcheck.ApiV1HealthCheckResponse, error) {
    // Call domain service
    output := t.healthcheckService.CheckDependencies(ctx)

    // Transform domain output to protobuf response
    return &healthcheck.ApiV1HealthCheckResponse{
        Status:    mapStatusToProto(output.Status),
        Timestamp: timestamppb.New(output.Timestamp),
        Dependencies: &healthcheck.Dependencies{
            Database: &healthcheck.Dependency{
                Status:       mapDependencyStatusToProto(output.Database.Status),
                ResponseTime: output.Database.ResponseTime.String(),
                Message:      output.Database.Message,
            },
        },
    }, nil
}
```

**REST API Transport Example:**

```go
package transporthealthcheck

import (
    "net/http"
    domainhealthcheck "project/internal/domain/healthcheck"
    "project/internal/gen/restapigen"
)

type TransportHealthCheckRestApi struct {
    healthcheckService domainhealthcheck.HealthCheckService
}

func NewTransportRestApi(
    healthcheckService domainhealthcheck.HealthCheckService,
) *TransportHealthCheckRestApi {
    return &TransportHealthCheckRestApi{
        healthcheckService: healthcheckService,
    }
}

func (t *TransportHealthCheckRestApi) ApiV1GetHealthCheck(
    w http.ResponseWriter, 
    r *http.Request,
) {
    // Call domain service
    output := t.healthcheckService.CheckDependencies(r.Context())

    // Transform domain output to REST API response
    resp := restapigen.ApiV1GetHealthCheckResponse{
        Status:    restapigen.Status(output.Status),
        Timestamp: output.Timestamp,
        Dependencies: restapigen.Dependencies{
            Database: &restapigen.Dependency{
                Status:       restapigen.DependencyStatus(output.Database.Status),
                ResponseTime: output.Database.ResponseTime.String(),
                Message:      output.Database.Message,
            },
        },
    }

    // Write JSON response
    writeJSON(w, http.StatusOK, resp)
}
```

**Transport Responsibilities:**

- **Request Parsing** - Extract and validate request parameters
- **Data Transformation** - Convert between protocol types and domain DTOs
- **Response Formatting** - Serialize domain outputs to protocol responses
- **Error Handling** - Map domain errors to appropriate HTTP/gRPC status codes
- **Protocol Concerns** - Handle HTTP headers, gRPC metadata, status codes

**Transport Does NOT:**

- Contain business logic
- Directly access databases or external services
- Perform data validation (domain layer responsibility)
- Make decisions based on business rules

### Worker Layer (`internal/worker/`)

The worker layer is responsible for **all background processing tasks**. This includes scheduled jobs (cron), message broker consumers (Kafka, RabbitMQ), and any asynchronous processing that runs independently of HTTP/gRPC requests.

```text
internal/worker/<feature>/
├── scheduler_<feature>.go  # Scheduled jobs (cron)
├── consumer_<feature>.go   # Message broker consumers (future)
└── job_<feature>.go        # Other background jobs (future)
```

**What Goes in Worker Layer:**

1. **Scheduled Jobs (Cron)** - Time-triggered tasks
   - Health checks
   - Data cleanup
   - Report generation
   - Scheduled notifications

2. **Message Consumers (Future)** - Event-driven tasks
   - Kafka consumers
   - RabbitMQ consumers
   - Redis pub/sub subscribers
   - Event processing

3. **Background Jobs** - Async processing
   - Email sending
   - File processing
   - Data synchronization
   - Batch operations

**Key principles:**

- Handles all background processing (scheduled, event-driven, async)
- No business logic - only triggering and calling services
- Each worker depends on domain service interfaces
- Thin layer that delegates to the domain/module layer
- Independent from HTTP/gRPC layers

**Scheduler (Cron) Example:**

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

**How Scheduler Works:**

The scheduler is configured in `internal/app/scheduler.go` and manages all cron jobs:

```go
package app

import (
    "github.com/robfig/cron/v3"
    workerhealthcheck "project/internal/worker/healthcheck"
)

type schedulerApp struct {
    cron *cron.Cron
}

func NewSchedulerApp() *schedulerApp {
    // Create cron with seconds support
    c := cron.New(cron.WithSeconds())
    
    // Initialize workers (dependency injection)
    healthcheckWorker := workerhealthcheck.NewSchedulerHealthCheck(
        healthcheckService,
    )
    
    // Register cron jobs
    appCfg := config.GetAppScheduler()
    c.AddFunc(appCfg.HealthCheckInterval, healthcheckWorker.CheckDependencies)
    
    return &schedulerApp{
        cron: c,
    }
}

func (s *schedulerApp) Start() {
    s.cron.Start()
}

func (s *schedulerApp) Stop() {
    s.cron.Stop()
}
```

**Cron Expression Format:**

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

**Common Cron Expressions:**

- `0 */5 * * * *` - Every 5 minutes
- `0 0 2 * * *` - Daily at 2 AM
- `0 0 */6 * * *` - Every 6 hours
- `0 30 9 * * 1-5` - 9:30 AM on weekdays
- `*/10 * * * * *` - Every 10 seconds

**Worker Responsibilities:**

- **Schedule Management** - Register and manage cron schedules
- **Event Processing** - Consume and process messages from brokers
- **Context Creation** - Create proper context with timeouts for jobs
- **Error Handling** - Handle errors and send alerts if jobs fail
- **Logging** - Log job execution for monitoring
- **Service Delegation** - Call domain services to perform actual work

**Worker Does NOT:**

- Contain business logic
- Directly access databases or external services
- Make decisions based on business rules
- Handle HTTP requests or gRPC calls (use Transport layer for that)

**Worker vs Transport Layer:**

| Aspect | Worker Layer | Transport Layer |
|--------|-------------|-----------------|
| **Trigger** | Time schedule, events, async | HTTP/gRPC requests |
| **Entry Point** | Cron, message broker, background | API endpoint |
| **Response** | No response (fire and forget) | Synchronous response required |
| **Use Cases** | Scheduled jobs, event processing | API requests, real-time operations |
| **Examples** | Daily cleanup, Kafka consumer | GET /users, gRPC GetUser |

### Other Important Directories

- `api/` - API specifications (OpenAPI YAML, Protocol Buffers)
- `cmd/` - Application entry points and CLI commands
- `internal/app/` - Application initialization (gRPC, REST API, routing)
- `internal/config/` - Configuration loading and types
- `internal/gen/` - Auto-generated code (excluded from git, must run `make generate` after cloning)
  - `grpcgen/` - Generated gRPC server stubs and protocol buffer types from `.proto` files in `api/proto`
  - `restapigen/` - Generated REST API server handlers and types from OpenAPI in files`api/openapi/api.yaml` specification
  - `mockgen/` - Generated mock implementations of domain interfaces for testing
- `internal/provider/` - Infrastructure providers (database, observability)
- `internal/transport/` - Transport layer handlers (gRPC, REST API)
- `internal/worker/` - Background job handlers (cron jobs, scheduled tasks)

## Architecture Flow

1. **API Specification** (`api/`) → Defines contracts
2. **Code Generation** → Creates server stubs and types
3. **Domain Layer** → Defines business interfaces
4. **Module Layer** → Implements business logic
5. **Entry Points** → Handles requests via:
   - Transport Layer (HTTP/gRPC)
   - Worker Layer (Cron jobs)
6. **Application Layer** → Wires everything together

## Development Workflow

1. Define your API contract in `api/openapi/api.yaml` or `api/proto/`
2. Create domain interfaces in `internal/domain/<feature>/`
3. Implement business logic in `internal/module/<feature>/`
4. Create entry point handlers:
   - Transport handlers in `internal/transport/<feature>/` for HTTP/gRPC
   - Worker handlers in `internal/worker/<feature>/` for cron jobs
5. Run `make generate` to generate code
6. Wire dependencies in `internal/app/`
7. Run the application with `make run-restapi` or `make run-grpcapi`

## Naming Conventions

### Folder vs Package Names

**IMPORTANT:** Folder names and package names follow different conventions:

**Folder Names:**

- Use simple, lowercase feature names
- No prefixes or suffixes
- Examples: `healthcheck/`, `user/`, `product/`

**Package Names:**

- Use descriptive names with prefixes indicating the layer
- Combine folder location + feature name
- Examples: `domainhealthcheck`, `healthcheckservice`, `transporthealthcheck`

**Examples:**

```text
# Domain Layer
Folder:  internal/domain/healthcheck/
Package: package domainhealthcheck

Folder:  internal/domain/user/
Package: package domainuser

# Module Layer - Service
Folder:  internal/module/healthcheck/service/
Package: package healthcheckservice

Folder:  internal/module/user/service/
Package: package userservice

# Module Layer - Repository
Folder:  internal/module/healthcheck/repository/
Package: package healthcheckrepository

Folder:  internal/module/user/repository/
Package: package userrepository

# Transport Layer
Folder:  internal/transport/healthcheck/
Package: package transporthealthcheck

Folder:  internal/transport/user/
Package: package transportuser

# Worker Layer
Folder:  internal/worker/healthcheck/
Package: package workerhealthcheck

Folder:  internal/worker/user/
Package: package workeruser
```

### Package Names Rules

- Use lowercase, single word package names when possible
- For domain packages, prefix with `domain`: `domainhealthcheck`, `domainuser`
- For service packages, use feature + `service`: `healthcheckservice`, `userservice`
- For repository packages, use feature + `repository`: `healthcheckrepository`, `userrepository`
- For transport packages, prefix with `transport`: `transporthealthcheck`, `transportuser`
- For worker packages, prefix with `worker`: `workerhealthcheck`, `workeruser`

**Example:**

```go
package domainhealthcheck    // domain layer
package healthcheckservice   // module service implementation
package transporthealthcheck // transport layer
package workerhealthcheck    // worker layer
```

### Interface Names

Interfaces should describe behavior and use clear, descriptive names:

**Service Interfaces:**

```go
type HealthCheckService interface {
    CheckDependencies(ctx context.Context) (output CheckDependenciesOutput)
}

type UserService interface {
    CreateUser(ctx context.Context, input CreateUserInput) (output CreateUserOutput, err error)
    GetUserByID(ctx context.Context, userID string) (output UserOutput, err error)
}
```

**Repository Interfaces:**

```go
type HealthCheckRepositoryDatastore interface {
    PingDatabase(ctx context.Context) (responseTime time.Duration, err error)
}

type UserRepositoryDatastore interface {
    Create(ctx context.Context, user User) (err error)
    FindByID(ctx context.Context, id string) (user User, err error)
}
```

### Struct Names

**Implementation Structs:**

Use lowercase `service` or `repository` for concrete implementations:

```go
type service struct {
    healthcheckRepo domainhealthcheck.HealthCheckRepositoryDatastore
}

type repository struct {
    db *sql.DB
}
```

**DTO/Value Object Structs:**

Use descriptive names with suffixes indicating their purpose:

```go
// Input DTOs - for incoming data
type CreateUserInput struct {
    Name  string
    Email string
}

// Output DTOs - for outgoing data
type CreateUserOutput struct {
    UserID    string
    CreatedAt time.Time
}

type UserOutput struct {
    ID        string
    Name      string
    Email     string
    CreatedAt time.Time
}

// Value Objects - domain entities
type User struct {
    ID        string
    Name      string
    Email     string
    Status    UserStatus
    CreatedAt time.Time
}
```

### File Names

Follow Go conventions with snake_case:

```text
// Domain layer
dto.go                  # Data Transfer Objects
repository.go           # Repository interfaces
service.go              # Service interfaces
value_object.go         # Domain entities and enums

// Module layer - Repository
repo_<feature>.go              # e.g., repo_healthcheck.go
repo_<feature>_datastore.go    # e.g., repo_healthcheck_datastore.go
repo_<feature>_cache.go        # e.g., repo_user_cache.go

// Module layer - Service
service_<feature>.go           # e.g., service_healthcheck.go
service_<feature>_test.go      # e.g., service_healthcheck_test.go

// Transport layer
grpc_<feature>.go              # e.g., grpc_healthcheck.go
restapi_<feature>.go           # e.g., restapi_healthcheck.go

// Worker layer
scheduler_<feature>.go         # e.g., scheduler_healthcheck.go
consumer_<feature>.go          # e.g., consumer_user.go (future)
```

### Transport Layer Naming

**Package Name:**

```go
package transporthealthcheck  // prefix 'transport' + feature name
package transportuser
package transportproduct
```

**Struct Names:**

Use descriptive names with `Transport` prefix and protocol suffix:

```go
// gRPC Transport
type TransportHealthCheckGrpc struct {
    healthcheckService domainhealthcheck.HealthCheckService
    healthcheck.UnimplementedHealthCheckServiceServer
}

type TransportUserGrpc struct {
    userService domainuser.UserService
    user.UnimplementedUserServiceServer
}

// REST API Transport
type TransportHealthCheckRestApi struct {
    healthcheckService domainhealthcheck.HealthCheckService
}

type TransportUserRestApi struct {
    userService domainuser.UserService
}
```

**Constructor Names:**

```go
// gRPC Transport Constructor
func NewTransportGrpc(
    healthcheckService domainhealthcheck.HealthCheckService,
) *TransportHealthCheckGrpc {
    return &TransportHealthCheckGrpc{
        healthcheckService: healthcheckService,
    }
}

// REST API Transport Constructor
func NewTransportRestApi(
    healthcheckService domainhealthcheck.HealthCheckService,
) *TransportHealthCheckRestApi {
    return &TransportHealthCheckRestApi{
        healthcheckService: healthcheckService,
    }
}
```

**Method Names:**

Methods should match the generated API contract:

```go
// gRPC - matches protobuf service definition
func (t *TransportHealthCheckGrpc) ApiV1HealthCheck(
    ctx context.Context, 
    req *healthcheck.ApiV1HealthCheckRequest,
) (*healthcheck.ApiV1HealthCheckResponse, error) {
    // Implementation
}

func (t *TransportUserGrpc) ApiV1CreateUser(
    ctx context.Context,
    req *user.ApiV1CreateUserRequest,
) (*user.ApiV1CreateUserResponse, error) {
    // Implementation
}

// REST API - matches OpenAPI operation IDs
func (t *TransportHealthCheckRestApi) ApiV1GetHealthCheck(
    w http.ResponseWriter, 
    r *http.Request,
) {
    // Implementation
}

func (t *TransportUserRestApi) ApiV1PostUsers(
    w http.ResponseWriter, 
    r *http.Request,
) {
    // Implementation
}
```

**File Organization:**

```text
internal/transport/healthcheck/
├── grpc_healthcheck.go      # Contains: TransportHealthCheckGrpc
└── restapi_healthcheck.go   # Contains: TransportHealthCheckRestApi

internal/transport/user/
├── grpc_user.go             # Contains: TransportUserGrpc
└── restapi_user.go          # Contains: TransportUserRestApi
```

### Worker Layer Naming

**Package Name:**

```go
package workerhealthcheck  // prefix 'worker' + feature name
package workeruser
package workerproduct
```

**Struct Names:**

Use descriptive names with `Scheduler` prefix for cron jobs:

```go
type SchedulerHealthCheck struct {
    healthcheckService domainhealthcheck.HealthCheckService
}

type SchedulerUser struct {
    userService domainuser.UserService
}
```

**Constructor Names:**

```go
func NewSchedulerHealthCheck(
    healthcheckService domainhealthcheck.HealthCheckService,
) *SchedulerHealthCheck {
    return &SchedulerHealthCheck{
        healthcheckService: healthcheckService,
    }
}
```

**Method Names:**

Method names should be descriptive of the job they perform (no prefix needed):

```go
func (w *SchedulerHealthCheck) CheckDependencies() {
    // Runs periodic health checks
}

func (w *SchedulerUser) SendDailyReport() {
    // Sends daily user activity reports
}

func (w *SchedulerUser) CleanupInactiveUsers() {
    // Removes users inactive for 90+ days
}
```

**File Names:**

```text
internal/worker/healthcheck/
└── scheduler_healthcheck.go  # Contains: SchedulerHealthCheck

internal/worker/user/
├── scheduler_user.go         # Contains: SchedulerUser
└── consumer_user.go          # Contains: ConsumerUser (future - for Kafka/RabbitMQ)
```

### Method Names

**Service Methods:**

Use descriptive verb-noun combinations:

```go
CheckDependencies(ctx context.Context) (output CheckDependenciesOutput)
CreateUser(ctx context.Context, input CreateUserInput) (output CreateUserOutput, err error)
GetUserByID(ctx context.Context, userID string) (output UserOutput, err error)
UpdateUserStatus(ctx context.Context, userID string, status UserStatus) (err error)
DeleteUser(ctx context.Context, userID string) (err error)
```

**Repository Methods:**

Use simple CRUD operations or specific queries:

```go
Create(ctx context.Context, user User) (err error)
FindByID(ctx context.Context, id string) (user User, err error)
FindByEmail(ctx context.Context, email string) (user User, err error)
Update(ctx context.Context, user User) (err error)
Delete(ctx context.Context, id string) (err error)
PingDatabase(ctx context.Context) (responseTime time.Duration, err error)
```

### Constants and Enums

Use typed constants with descriptive names:

```go
type StatusHealthCheck string

const (
    StatusHealthCheckOk       StatusHealthCheck = "OK"
    StatusHealthCheckDegraded StatusHealthCheck = "DEGRADED"
    StatusHealthCheckDown     StatusHealthCheck = "DOWN"
)

type UserStatus string

const (
    UserStatusActive   UserStatus = "ACTIVE"
    UserStatusInactive UserStatus = "INACTIVE"
    UserStatusSuspended UserStatus = "SUSPENDED"
)
```

### Function Parameters and Return Values

**Named Return Values:**

Use named returns for clarity, especially for multiple return values:

```go
func (s *service) CheckDependencies(ctx context.Context) (output CheckDependenciesOutput) {
    // Implementation
    return
}

func (s *service) CreateUser(ctx context.Context, input CreateUserInput) (output CreateUserOutput, err error) {
    // Implementation
    return
}
```

**Input/Output Pattern:**

For service methods with complex data:

```go
type CreateUserInput struct {
    Name     string
    Email    string
    Password string
}

type CreateUserOutput struct {
    UserID    string
    CreatedAt time.Time
}

func (s *service) CreateUser(ctx context.Context, input CreateUserInput) (output CreateUserOutput, err error) {
    // Implementation
}
```

### Constructor Functions

Use `New` prefix for constructors:

```go
func NewService(
    healthcheckRepo domainhealthcheck.HealthCheckRepositoryDatastore,
) *service {
    return &service{
        healthcheckRepo: healthcheckRepo,
    }
}

func NewRepository(db *sql.DB) *repository {
    return &repository{
        db: db,
    }
}
```

## Testing

The project uses mockgen for generating mocks. Mocks are automatically generated when you run `make generate` or `make go_generate`.

To generate mocks for a specific interface, add a `go:generate` directive in your domain file:

```go
//go:generate go tool mockgen -source=service.go -destination=../../gen/mockgen/feature_service_mock.gen.go -package=mockgen
```

## Performance & Debugging

### Memory Leak Detection

The application includes built-in pprof endpoints for profiling and memory leak detection. For detailed instructions on detecting and analyzing memory leaks, see:

**[Memory Leak Detection Guide](docs/MEMORY_LEAK_DETECTION.md)**

Quick access to pprof endpoints (replace `<PORT>` with your application port):

- Heap Profile: `http://localhost:<PORT>/debug/pprof/heap`
- Goroutine Profile: `http://localhost:<PORT>/debug/pprof/goroutine`
- All Profiles: `http://localhost:<PORT>/debug/pprof/`

Example workflow:

```bash
# Capture baseline
curl http://localhost:8080/debug/pprof/heap > heap_before.prof

# Run load test
hey -n 10000 -c 100 http://localhost:8080/api/v1/healthcheck

# Capture after load
curl http://localhost:8080/debug/pprof/heap > heap_after.prof

# Compare profiles
go tool pprof -base=heap_before.prof heap_after.prof
```

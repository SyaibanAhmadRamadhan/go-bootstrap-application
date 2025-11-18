# Project Structure

This document provides a detailed overview of the project structure and architectural layers.

## Table of Contents

- [Layer Overview](#layer-overview)
- [Domain Layer](#domain-layer)
- [Module Layer](#module-layer)
- [Transport Layer](#transport-layer)
- [Worker Layer](#worker-layer)
- [Other Important Directories](#other-important-directories)
- [Architecture Flow](#architecture-flow)

## Layer Overview

The project follows a clean architecture pattern with clear separation of concerns:

```text
├── api/                    # API specifications
├── cmd/                    # Application entry points
├── docs/                   # Documentation
├── internal/
│   ├── domain/            # Business interfaces & DTOs
│   ├── module/            # Business logic implementation
│   ├── transport/         # HTTP/gRPC handlers
│   ├── worker/            # Background jobs & schedulers
│   ├── app/               # Application initialization
│   ├── config/            # Configuration management
│   ├── gen/               # Auto-generated code
│   └── provider/          # Infrastructure providers
└── scripts/               # Build & deployment scripts
```

## Domain Layer

**Location:** `internal/domain/`

The domain layer contains business logic interfaces and data structures. Each domain represents a business capability.

### Structure

```text
internal/domain/<feature>/
├── dto.go              # Data Transfer Objects
├── repository.go       # Repository interface definitions
├── service.go          # Service interface definitions
└── value_object.go     # Domain value objects and enums
```

### Key Principles

- Defines **interfaces only**, no implementation
- Contains DTOs, value objects, and business rules
- Framework-agnostic and pure Go code
- Generates mock interfaces for testing using `go:generate` directives

### Example

```go
type HealthCheckService interface {
    CheckDependencies(ctx context.Context) (output CheckDependenciesOutput)
}
```

## Module Layer

**Location:** `internal/module/`

The module layer contains concrete implementations of domain interfaces. This is where the actual business logic lives.

### Structure

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

### Key Principles

- Implements domain interfaces
- Contains actual business logic
- Handles data persistence and external dependencies
- Separated into `repository` (data access) and `service` (business logic)
- Each implementation can have multiple variants (e.g., datastore, cache, external API)

### Example

```go
type service struct {
    healthcheckRepo domainhealthcheck.HealthCheckRepositoryDatastore
}

func (s *service) CheckDependencies(ctx context.Context) (output domainhealthcheck.CheckDependenciesOutput) {
    // Business logic implementation
}
```

## Transport Layer

**Location:** `internal/transport/`

The transport layer acts as the adapter between external protocols (HTTP REST, gRPC) and your domain/business logic.

### Structure

```text
internal/transport/<feature>/
├── grpc_<feature>.go      # gRPC handlers
└── restapi_<feature>.go   # REST API handlers
```

### Key Principles

- Handles protocol-specific request/response formatting
- Transforms between protocol types (protobuf, HTTP) and domain DTOs
- No business logic - only marshalling/unmarshalling and calling services
- Each transport handler depends on domain service interfaces
- Thin layer that delegates to the domain/module layer

### Responsibilities

- **Request Parsing** - Extract and validate request parameters
- **Data Transformation** - Convert between protocol types and domain DTOs
- **Response Formatting** - Serialize domain outputs to protocol responses
- **Error Handling** - Map domain errors to appropriate HTTP/gRPC status codes
- **Protocol Concerns** - Handle HTTP headers, gRPC metadata, status codes

### Does NOT Handle

- Business logic
- Direct database or external service access
- Data validation (domain layer responsibility)
- Business rule decisions

### Examples

See [TRANSPORT_EXAMPLES.md](TRANSPORT_EXAMPLES.md) for detailed gRPC and REST API transport examples.

## Worker Layer

**Location:** `internal/worker/`

The worker layer is responsible for **all background processing tasks**.

### Structure

```text
internal/worker/<feature>/
├── scheduler_<feature>.go  # Scheduled jobs (cron)
├── consumer_<feature>.go   # Message broker consumers (future)
└── job_<feature>.go        # Other background jobs (future)
```

### What Goes in Worker Layer

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

### Key Principles

- Handles all background processing (scheduled, event-driven, async)
- No business logic - only triggering and calling services
- Each worker depends on domain service interfaces
- Thin layer that delegates to the domain/module layer
- Independent from HTTP/gRPC layers

### Responsibilities

- **Schedule Management** - Register and manage cron schedules
- **Event Processing** - Consume and process messages from brokers
- **Context Creation** - Create proper context with timeouts for jobs
- **Error Handling** - Handle errors and send alerts if jobs fail
- **Logging** - Log job execution for monitoring
- **Service Delegation** - Call domain services to perform actual work

### Does NOT Handle

- Business logic
- Direct database or external service access
- Business rule decisions
- HTTP requests or gRPC calls (use Transport layer)

### Worker vs Transport

| Aspect | Worker Layer | Transport Layer |
|--------|-------------|-----------------|
| **Trigger** | Time schedule, events, async | HTTP/gRPC requests |
| **Entry Point** | Cron, message broker, background | API endpoint |
| **Response** | No response (fire and forget) | Synchronous response required |
| **Use Cases** | Scheduled jobs, event processing | API requests, real-time operations |
| **Examples** | Daily cleanup, Kafka consumer | GET /users, gRPC GetUser |

### Examples

See [WORKER_EXAMPLES.md](WORKER_EXAMPLES.md) for detailed scheduler and cron job examples.

## Other Important Directories

### API Specifications (`api/`)

Contains API contract definitions:

- `api/openapi/` - OpenAPI/Swagger specifications
- `api/proto/` - Protocol Buffer definitions

### Application Entry Points (`cmd/`)

CLI commands and application initialization:

- `cmd_rest_api.go` - REST API server command
- `cmd_grpc_api.go` - gRPC API server command
- `cmd_scheduler.go` - Scheduler command

### Application Layer (`internal/app/`)

Application wiring and initialization:

- `restapi.go` - REST API app setup
- `grpc.go` - gRPC app setup
- `scheduler.go` - Scheduler app setup
- `pprof.go` - Pprof server management

### Configuration (`internal/config/`)

Configuration management:

- `type.go` - Configuration struct definitions
- `load_config.go` - Config loading and getters

### Generated Code (`internal/gen/`)

Auto-generated code (excluded from git):

- `grpcgen/` - Generated gRPC stubs from `.proto` files
- `restapigen/` - Generated REST API handlers from OpenAPI spec
- `mockgen/` - Generated mock implementations for testing

**Important:** Run `make generate` after cloning to generate these files.

### Infrastructure Providers (`internal/provider/`)

Infrastructure and cross-cutting concerns:

- `database.go` - Database connection management
- `observability.go` - Logging and monitoring setup

### Worker Handlers (`internal/worker/`)

Background job implementations:

- Scheduled tasks (cron jobs)
- Message consumers (future)
- Async processing jobs (future)

## Architecture Flow

1. **API Specification** (`api/`) → Defines contracts
2. **Code Generation** (`make generate`) → Creates server stubs and types
3. **Domain Layer** (`internal/domain/`) → Defines business interfaces
4. **Module Layer** (`internal/module/`) → Implements business logic
5. **Entry Points** → Handles requests via:
   - **Transport Layer** (`internal/transport/`) - HTTP/gRPC
   - **Worker Layer** (`internal/worker/`) - Cron jobs
6. **Application Layer** (`internal/app/`) → Wires everything together
7. **CMD Layer** (`cmd/`) → Starts the application

### Request Flow Example

```text
HTTP Request
    ↓
Transport Layer (restapi_healthcheck.go)
    ↓
Domain Service Interface
    ↓
Module Service Implementation (service_healthcheck.go)
    ↓
Domain Repository Interface
    ↓
Module Repository Implementation (repo_healthcheck_datastore.go)
    ↓
Database
```

### Worker Flow Example

```text
Cron Schedule (config: "0 */5 * * * *")
    ↓
Scheduler App (scheduler.go)
    ↓
Worker Handler (scheduler_healthcheck.go)
    ↓
Domain Service Interface
    ↓
Module Service Implementation
    ↓
Domain Repository Interface
    ↓
Module Repository Implementation
    ↓
Database / External Service
```

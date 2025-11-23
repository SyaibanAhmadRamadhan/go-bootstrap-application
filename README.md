# GO SERVICE BOOTSTRAP

A starter template for building internal backend services with Go.  
This repository provides a ready-to-use project layout, wiring, and tooling so you don't have to set up everything from scratch for every new service.

Designed for a "distributed monolith" architecture: many services, one database, with clear ownership and strict API contracts between services.

## Table of Contents

- [GO SERVICE BOOTSTRAP](#go-service-bootstrap)
  - [Table of Contents](#table-of-contents)
  - [Quick Start](#quick-start)
  - [Prerequisites](#prerequisites)
  - [Database Setup](#database-setup)
    - [Using Docker Compose](#using-docker-compose)
    - [Manual PostgreSQL Setup](#manual-postgresql-setup)
  - [Documentation](#documentation)
  - [Makefile Commands](#makefile-commands)
    - [Running the Application](#running-the-application)
    - [Code Generation](#code-generation)
    - [Utilities](#utilities)
  - [Project Structure](#project-structure)
  - [Architecture Overview](#architecture-overview)
    - [Layer Responsibilities](#layer-responsibilities)
  - [Development Workflow](#development-workflow)
  - [Error Handling](#error-handling)
    - [Error Types](#error-types)
    - [Usage Example in Service Layer](#usage-example-in-service-layer)
    - [Transport Layer with GinHelper](#transport-layer-with-ginhelper)
    - [GinHelper Features](#ginhelper-features)
    - [Benefits](#benefits)
    - [Error Response Format](#error-response-format)
  - [Configuration](#configuration)
    - [Example Configuration](#example-configuration)
  - [Performance \& Debugging](#performance--debugging)
    - [Pprof Profiling](#pprof-profiling)
    - [Memory Leak Detection](#memory-leak-detection)
  - [Testing](#testing)
  - [Key Features](#key-features)
  - [Examples](#examples)
    - [Quick Examples](#quick-examples)

## Quick Start

```bash
# 1. Clone the repository
git clone <repository-url>

# 2. Start PostgreSQL with Docker Compose
docker-compose up -d

# 3. Install dependencies
bash scripts/install_dependency.sh

# 4. Create configuration
cp env.example.json env.json

# 5. Update database DSN in env.json
# "dsn": "postgres://postgres:postgres@localhost:5432/go_bootstrap?sslmode=disable"

# 6. Generate code
make generate

# 7. Run the application
make run-restapi    # REST API server
make run-grpcapi    # gRPC API server
make run-scheduler  # Background scheduler
```

## Prerequisites

- Go 1.24 or higher
- Docker & Docker Compose (for PostgreSQL)
- Node.js and npm (for OpenAPI tooling)
- Make

## Database Setup

### Using Docker Compose

The easiest way to set up PostgreSQL for development:

```bash
# Start PostgreSQL container
docker-compose up -d

# Check container status
docker-compose ps

# View logs
docker-compose logs -f postgres

# Stop PostgreSQL
docker-compose down

# Stop and remove data
docker-compose down -v
```

**Database Configuration:**

- **Host:** localhost
- **Port:** 5432
- **Database:** go_bootstrap
- **User:** postgres
- **Password:** postgres
- **DSN:** `postgres://postgres:postgres@localhost:5432/go_bootstrap?sslmode=disable`

**Migrations:** Automatically run on container startup via `docker-entrypoint-initdb.d`.

### Manual PostgreSQL Setup

If you prefer manual setup:

```bash
# 1. Install PostgreSQL
brew install postgresql@16  # macOS
# or use your OS package manager

# 2. Start PostgreSQL
brew services start postgresql@16

# 3. Create database and user
psql postgres
CREATE DATABASE go_bootstrap;
CREATE USER go_user WITH ENCRYPTED PASSWORD 'go_password';
GRANT ALL PRIVILEGES ON DATABASE go_bootstrap TO go_user;
\q

# 4. Run migrations manually
psql -U go_user -d go_bootstrap -f migrations/001_create_auth_and_user_tables.sql

# 5. Update DSN in env.json
# "dsn": "postgres://go_user:go_password@localhost:5432/go_bootstrap?sslmode=disable"
```

## Documentation

- **[Docker Setup Guide](docs/DOCKER_SETUP.md)** - PostgreSQL setup with Docker Compose
- **[Configuration Guide](docs/CONFIGURATION.md)** - Configuration structure, pprof hot-reload, and settings
- **[Project Structure](docs/PROJECT_STRUCTURE.md)** - Architecture layers and directory organization
- **[Naming Conventions](docs/NAMING_CONVENTIONS.md)** - Package, file, and code naming standards
- **[Transport Examples](docs/TRANSPORT_EXAMPLES.md)** - gRPC and REST API handler examples
- **[Worker Examples](docs/WORKER_EXAMPLES.md)** - Scheduler and cron job implementation guides
- **[Memory Leak Detection](docs/MEMORY_LEAK_DETECTION.md)** - Performance profiling and debugging

## Makefile Commands

### Running the Application

```bash
make run-restapi    # Start REST API server (port 8080)
make run-grpcapi    # Start gRPC API server (port 9090)
make run-scheduler  # Start background scheduler
```

### Code Generation

```bash
make generate       # Run all code generation (API, gRPC, mocks)
make api_generate   # Generate REST API server from OpenAPI spec
make buf_generate   # Generate gRPC code from protobuf
make go_generate    # Generate mocks with mockgen
```

### Utilities

```bash
make preview_open_api  # Preview OpenAPI docs in browser
make clean             # Clean generated files and artifacts
```

## Project Structure

```text
├── api/                    # API specifications (OpenAPI, protobuf)
├── cmd/                    # Application entry points
├── docs/                   # Documentation
├── internal/
│   ├── app/               # Application initialization
│   ├── config/            # Configuration management
│   ├── domain/            # Business interfaces & DTOs (clean architecture)
│   ├── gen/               # Auto-generated code (gitignored)
│   │   ├── grpcgen/      # Generated gRPC stubs
│   │   ├── mockgen/      # Generated mocks
│   │   └── restapigen/   # Generated REST API handlers
│   ├── module/            # Business logic implementation
│   │   ├── <feature>/repository/  # Data access layer
│   │   └── <feature>/service/     # Business logic layer
│   ├── infrastructure/    # Infrastructure setup (DB, logging, observability)
│   ├── transport/         # HTTP/gRPC handlers
│   └── worker/            # Background jobs & schedulers
├── scripts/               # Build & deployment scripts
└── makefile              # Build automation
```

## Architecture Overview

This project follows clean architecture principles with clear separation of concerns:

```text
┌─────────────────────────────────────────────────────┐
│                   Entry Points                       │
│  (REST API, gRPC API, Scheduler, CLI Commands)      │
└─────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────┐
│              Transport & Worker Layers               │
│  • Transport: HTTP/gRPC handlers                     │
│  • Worker: Cron jobs, message consumers             │
│  (Protocol-specific, thin delegation layer)         │
└─────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────┐
│                   Domain Layer                       │
│  • Interfaces (Service, Repository)                 │
│  • DTOs (Input/Output)                              │
│  • Value Objects & Enums                            │
│  (Framework-agnostic, pure Go)                      │
└─────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────┐
│                   Module Layer                       │
│  • Service Implementation (business logic)          │
│  • Repository Implementation (data access)          │
│  (Concrete implementations of domain interfaces)    │
└─────────────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────┐
│                  Infrastructure                      │
│  (Database, External APIs, File System, etc.)       │
└─────────────────────────────────────────────────────┘
```

### Layer Responsibilities

| Layer | Purpose | Examples |
|-------|---------|----------|
| **Transport** | Protocol handlers (HTTP/gRPC) | REST endpoints, gRPC services |
| **Worker** | Background processing | Cron jobs, Kafka consumers |
| **Domain** | Business interfaces & DTOs | Service interfaces, Input/Output DTOs |
| **Module** | Business logic implementation | Service & repository implementations |
| **Provider** | Infrastructure setup | Database connections, logging |

## Development Workflow

1. **Define API Contract** - Create OpenAPI spec (`api/openapi/`) or protobuf (`api/proto/`)
2. **Generate Code** - Run `make generate` to create server stubs and types
3. **Define Domain** - Create interfaces in `internal/domain/<feature>/`
4. **Implement Business Logic** - Write implementations in `internal/module/<feature>/`
5. **Create Handlers** - Add transport handlers (`internal/transport/`) or workers (`internal/worker/`)
6. **Wire Dependencies** - Configure in `internal/app/`
7. **Run & Test** - Use `make run-restapi` or `make run-grpcapi`

## Error Handling

This project uses standardized error handling with **AppError** from `go-foundation-kit/apperror`:

### Error Types

```go
// BadRequest - 400 (client error, validation failed, etc.)
apperror.BadRequest("email already registered")

// Unauthorized - 401 (authentication required)
apperror.Unauthorized("invalid credentials")

// Forbidden - 403 (authenticated but not authorized)
apperror.Forbidden("insufficient permissions")

// NotFound - 404 (resource not found)
apperror.NotFound("user not found")

// Conflict - 409 (resource conflict)
apperror.Conflict("email already exists")

// StdUnknown - 500 (wrap unexpected errors)
apperror.StdUnknown(err)
```

### Usage Example in Service Layer

```go
func (s *service) Register(ctx context.Context, input RegisterInput) (RegisterOutput, error) {
    // Check if user exists
    existingUser, _ := s.userRepo.GetDetailUser(ctx, filters)
    if existingUser.ID != "" {
        return RegisterOutput{}, apperror.BadRequest("email already registered")
    }
    
    // Handle unexpected errors
    passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return RegisterOutput{}, apperror.StdUnknown(err)
    }
    
    // Continue...
}
```

### Transport Layer with GinHelper

The transport layer uses **`ginx.GinHelper`** to handle requests and responses automatically:

```go
type UserRestAPIHandler struct {
    userService domainuser.UserService
    helper      *ginx.GinHelper  // Helper for request/response handling
}

func NewRestAPIHandler(userService domainuser.UserService, helper *ginx.GinHelper) *UserRestAPIHandler {
    return &UserRestAPIHandler{
        userService: userService,
        helper:      helper,
    }
}

func (h *UserRestAPIHandler) ApiV1PostUsersRegister(c *gin.Context) {
    var req restapigen.ApiV1PostUsersRegisterRequest
    
    // Automatically bind request body and validate
    if err := h.helper.BindBody(c, &req); err != nil {
        return  // Error response handled automatically
    }
    
    // Call service
    output, err := h.userService.Register(c.Request.Context(), domainuser.RegisterInput{
        Email:    req.Email,
        Password: req.Password,
        Name:     req.Name,
    })
    
    // Automatically handle error response with proper status code
    if err != nil {
        h.helper.ResponseError(c, err)
        return
    }
    
    // Automatically send success response
    h.helper.ResponseSuccess(c, http.StatusCreated, restapigen.ApiV1PostUsersRegisterResponse{
        UserId:    output.UserID,
        Email:     output.Email,
        Name:      output.Name,
        CreatedAt: output.CreatedAt,
    })
}
```

### GinHelper Features

**Initialization:**

```go
// In app initialization (internal/app/app_rest_api.go)
ginHelper := ginx.NewGinHelper("message", "errors")  // JSON keys for error responses
```

**Key Methods:**

1. **`MustShouldBind(c, &req)`** - Bind and validate request body
   - Returns validation errors automatically
   - Supports struct tags validation

2. **`ResponseError(c, err)`** - Handle error response
   - Maps `apperror` to proper HTTP status codes
   - Returns structured JSON error format
   - Example: `{"message": "email already registered"}`

### Benefits

- **Automatic Error Mapping** - AppError → HTTP status codes (BadRequest=400, NotFound=404, etc.)
- **Structured Responses** - Consistent JSON format for success and errors
- **Validation Handling** - Built-in request validation with detailed error messages
- **Less Boilerplate** - No manual error handling in transport layer
- **Type Safety** - Works seamlessly with generated OpenAPI types

### Error Response Format

```json
// Success Response
{
    "user_id": "12345",
    "email": "user@example.com",
    "name": "John Doe"
}

// Error Response (validation)
{
    "message": "Validation failed",
    "errors": [
        {
            "field": "email",
            "message": ["Email is required", "Email is invalid"]
        }
    ]
}

// Error Response (business logic)
{
    "message": "email already registered"
}
```

See implementations in `internal/transport/auth/restapi_auth.go` and `internal/transport/user/restapi_user.go`.

## Configuration

Configuration is JSON-based (`env.json`) with hot-reload support during development.

### Example Configuration

```json
{
    "app_rest_api": {
        "name": "my-service-rest-api",
        "env": "development",
        "debug_mode": true,
        "port": 8080,
        "pprof": {
            "enable": true,
            "port": 8080,
            "static_token": "secret-token"
        }
    },
    "app_grpc_api": {
        "name": "my-service-grpc-api",
        "env": "development",
        "debug_mode": true,
        "port": 9090,
        "pprof": {
            "enable": false,
            "port": 7070,
            "static_token": "secret-token"
        }
    },
    "app_scheduler": {
        "name": "my-service-scheduler",
        "env": "development",
        "debug_mode": true,
        "healthcheck_interval": "0 */5 * * * *",
        "pprof": {
            "enable": true,
            "port": 6060,
            "static_token": "secret-token"
        }
    },
    "database": {
        "dsn": "user:password@tcp(host:port)/dbname?parseTime=true",
        "max_open_conns": 25,
        "max_idle_conns": 25,
        "conn_max_lifetime": "300s",
        "conn_max_idle_time": "60s"
    }
}
```

**Key Features:**

- **Independent configs per service** - REST API, gRPC API, and Scheduler have separate configs
- **Hot-reload** - Changes automatically reloaded during development
- **Nested pprof** - Each service has its own profiling configuration
- **Type-safe** - Strongly typed configuration structs

See [Configuration Guide](docs/CONFIGURATION.md) for details.

## Performance & Debugging

### Pprof Profiling

Each application has built-in pprof support with realtime hot-reload and authentication:

```bash
# Enable pprof in env.json
{
    "app_rest_api": {
        "pprof": {
            "enable": true,
            "port": 8080,
            "static_token": "your-secret-token"
        }
    }
}

# Access pprof endpoints (authentication required)
curl -H "Authotization: your-secret-token" http://localhost:8080/debug/pprof/heap
curl -H "Authotization: your-secret-token" http://localhost:8080/debug/pprof/goroutine

# Or set header in browser/tools
open http://localhost:8080/debug/pprof/  # Add header: Authotization: your-secret-token
```

### Memory Leak Detection

```bash
# Capture baseline (with authentication)
curl -H "Authotization: your-secret-token" \
  http://localhost:8080/debug/pprof/heap > heap_before.prof

# Run load test
hey -n 10000 -c 100 http://localhost:8080/api/v1/healthcheck

# Capture after load
curl -H "Authotization: your-secret-token" \
  http://localhost:8080/debug/pprof/heap > heap_after.prof

# Compare profiles
go tool pprof -base=heap_before.prof heap_after.prof
```

See [Memory Leak Detection Guide](docs/MEMORY_LEAK_DETECTION.md) for detailed instructions.

## Testing

The project uses mockgen for generating mocks:

```go
// In domain file (e.g., internal/domain/user/service.go)
//go:generate go tool mockgen -source=service.go -destination=../../gen/mockgen/user_service_mock.gen.go -package=mockgen

// Generate mocks
make go_generate
```

## Key Features

- **✅ Clean Architecture** - Clear separation of concerns with domain/module/transport layers
- **✅ Hot-Reload Config** - Configuration changes detected automatically in development
- **✅ Code Generation** - OpenAPI → REST handlers, Protobuf → gRPC stubs, Mockgen → Test mocks
- **✅ Multiple Protocols** - REST API, gRPC API, and background workers in one project
- **✅ Query Builder** - Squirrel SQL query builder for type-safe database queries
- **✅ Error Handling** - Standardized AppError from go-foundation-kit for consistent error responses
- **✅ Structured Logging** - slog-based structured logging throughout
- **✅ Built-in Profiling** - pprof endpoints with hot-reload for performance analysis
- **✅ Graceful Shutdown** - Proper cleanup and shutdown handling
- **✅ Dependency Injection** - Clear dependency wiring in app layer
- **✅ Background Jobs** - Cron scheduler with robfig/cron v3 (supports seconds)
- **✅ Testing Support** - Mock generation for all domain interfaces

## Examples

### Quick Examples

**Create a new feature:**

```bash
# 1. Define API in api/openapi/api.yaml or api/proto/
# 2. Create domain interfaces
mkdir -p internal/domain/product
touch internal/domain/product/{dto.go,service.go,repository.go,value_object.go}

# 3. Implement business logic
mkdir -p internal/module/product/{service,repository}

# 4. Create transport handlers
mkdir -p internal/transport/product
touch internal/transport/product/{grpc_product.go,restapi_product.go}

# 5. Generate code
make generate

# 6. Wire dependencies in internal/app/
```

**Add a cron job:**

```bash
# 1. Create worker
mkdir -p internal/worker/product
touch internal/worker/product/scheduler_product.go

# 2. Implement job method
func (w *SchedulerProduct) CleanupOldData() {
    // Job implementation
}

# 3. Register in internal/app/scheduler.go
c.AddFunc("0 0 2 * * *", productWorker.CleanupOldData)
```

See detailed examples in:

- [Transport Examples](docs/TRANSPORT_EXAMPLES.md) - gRPC & REST API handlers
- [Worker Examples](docs/WORKER_EXAMPLES.md) - Cron jobs & schedulers

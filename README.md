# GO SERVICE BOOTSTRAP

A starter template for building internal backend services with Go.  
This repository provides a ready-to-use project layout, wiring, and tooling so you don't have to set up everything from scratch for every new service.

Designed for a "distributed monolith" architecture: many services, one database, with clear ownership and strict API contracts between services.

## Quick Start

```bash
# 1. Clone the repository
git clone <repository-url>

# 2. Install dependencies
bash scripts/install_dependency.sh

# 3. Create configuration
cp env.example.json env.json

# 4. Generate code
make generate

# 5. Run the application
make run-restapi    # REST API server
make run-grpcapi    # gRPC API server
make run-scheduler  # Background scheduler
```

## Prerequisites

- Go 1.24 or higher
- Node.js and npm (for OpenAPI tooling)
- Make

## Documentation

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

## Contributing

1. Follow the [Naming Conventions](docs/NAMING_CONVENTIONS.md)
2. Keep layers properly separated (no business logic in transport/worker)
3. Write tests with generated mocks
4. Run `make generate` before committing
5. Use structured logging (slog)

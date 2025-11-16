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
  - [Project Structure](#project-structure)
    - [Domain Layer (`internal/domain/`)](#domain-layer-internaldomain)
    - [Module Layer (`internal/module/`)](#module-layer-internalmodule)
    - [Transport Layer (`internal/transport/`)](#transport-layer-internaltransport)
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

### Code Generation

- `make generate` - Run all code generation (API, gRPC, and mocks)
- `make api_generate` - Generate REST API server and types from OpenAPI spec
- `make buf_generate` - Generate gRPC code from Protocol Buffers
- `make go_generate` - Run Go's built-in code generation (e.g., mockgen)

### Utilities

- `make preview_open_api` - Preview OpenAPI documentation in browser using Redocly
- `make clean` - Clean generated files, builds, and artifacts

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

### Other Important Directories

- `api/` - API specifications (OpenAPI YAML, Protocol Buffers)
- `cmd/` - Application entry points and CLI commands
- `internal/app/` - Application initialization (gRPC, REST API, routing)
- `internal/config/` - Configuration loading and types
- `internal/gen/` - Auto-generated code (gRPC, REST API, mocks)
- `internal/provider/` - Infrastructure providers (database, observability)
- `internal/transport/` - Transport layer handlers (gRPC, REST API)

## Architecture Flow

1. **API Specification** (`api/`) → Defines contracts
2. **Code Generation** → Creates server stubs and types
3. **Domain Layer** → Defines business interfaces
4. **Module Layer** → Implements business logic
5. **Transport Layer** → Handles HTTP/gRPC requests
6. **Application Layer** → Wires everything together

## Development Workflow

1. Define your API contract in `api/openapi/api.yaml` or `api/proto/`
2. Create domain interfaces in `internal/domain/<feature>/`
3. Implement business logic in `internal/module/<feature>/`
4. Create transport handlers in `internal/transport/<feature>/`
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
```

### Package Names Rules

- Use lowercase, single word package names when possible
- For domain packages, prefix with `domain`: `domainhealthcheck`, `domainuser`
- For service packages, use feature + `service`: `healthcheckservice`, `userservice`
- For repository packages, use feature + `repository`: `healthcheckrepository`, `userrepository`
- For transport packages, prefix with `transport`: `transporthealthcheck`, `transportuser`

**Example:**

```go
package domainhealthcheck  // domain layer
package healthcheckservice // module service implementation
package transporthealthcheck // transport layer
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

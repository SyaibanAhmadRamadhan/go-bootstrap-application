# Naming Conventions

This document outlines the naming conventions used throughout the project to maintain consistency and readability.

## Table of Contents

- [Naming Conventions](#naming-conventions)
  - [Table of Contents](#table-of-contents)
  - [Folder vs Package Names](#folder-vs-package-names)
    - [Folder Names](#folder-names)
    - [Package Names](#package-names)
    - [Examples](#examples)
  - [Package Names Rules](#package-names-rules)
    - [Example](#example)
  - [Interface Names](#interface-names)
    - [Service Interfaces](#service-interfaces)
    - [Repository Interfaces](#repository-interfaces)
  - [Struct Names](#struct-names)
    - [Implementation Structs](#implementation-structs)
    - [DTO/Value Object Structs](#dtovalue-object-structs)
  - [File Names](#file-names)
  - [Transport Layer Naming](#transport-layer-naming)
    - [Package Name](#package-name)
    - [Struct Names](#struct-names-1)
    - [Constructor Names](#constructor-names)
    - [Method Names](#method-names)
    - [File Organization](#file-organization)
  - [Worker Layer Naming](#worker-layer-naming)
    - [Package Name](#package-name-1)
    - [Struct Names](#struct-names-2)
    - [Constructor Names](#constructor-names-1)
    - [Method Names](#method-names-1)
    - [File Names](#file-names-1)
  - [Method Names](#method-names-2)
    - [Service Methods](#service-methods)
    - [Repository Methods](#repository-methods)
  - [Constants and Enums](#constants-and-enums)
  - [Function Parameters and Return Values](#function-parameters-and-return-values)
    - [Named Return Values](#named-return-values)
    - [Input/Output Pattern](#inputoutput-pattern)
  - [Constructor Functions](#constructor-functions)

## Folder vs Package Names

**IMPORTANT:** Folder names and package names follow different conventions.

### Folder Names

- Use simple, lowercase feature names
- No prefixes or suffixes
- Examples: `healthcheck/`, `user/`, `product/`

### Package Names

- Use descriptive names with prefixes indicating the layer
- Combine folder location + feature name
- Examples: `domainhealthcheck`, `healthcheckservice`, `transporthealthcheck`

### Examples

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

## Package Names Rules

- Use lowercase, single word package names when possible
- For domain packages, prefix with `domain`: `domainhealthcheck`, `domainuser`
- For service packages, use feature + `service`: `healthcheckservice`, `userservice`
- For repository packages, use feature + `repository`: `healthcheckrepository`, `userrepository`
- For transport packages, prefix with `transport`: `transporthealthcheck`, `transportuser`
- For worker packages, prefix with `worker`: `workerhealthcheck`, `workeruser`

### Example

```go
package domainhealthcheck    // domain layer
package healthcheckservice   // module service implementation
package transporthealthcheck // transport layer
package workerhealthcheck    // worker layer
```

## Interface Names

Interfaces should describe behavior and use clear, descriptive names.

### Service Interfaces

```go
type HealthCheckService interface {
    CheckDependencies(ctx context.Context) (output CheckDependenciesOutput)
}

type UserService interface {
    CreateUser(ctx context.Context, input CreateUserInput) (output CreateUserOutput, err error)
    GetUserByID(ctx context.Context, userID string) (output UserOutput, err error)
}
```

### Repository Interfaces

```go
type HealthCheckRepositoryDatastore interface {
    PingDatabase(ctx context.Context) (responseTime time.Duration, err error)
}

type UserRepositoryDatastore interface {
    Create(ctx context.Context, user User) (err error)
    FindByID(ctx context.Context, id string) (user User, err error)
}
```

## Struct Names

### Implementation Structs

Use lowercase `service` or `repository` for concrete implementations:

```go
type service struct {
    healthcheckRepo domainhealthcheck.HealthCheckRepositoryDatastore
}

type repository struct {
    db *sql.DB
}
```

### DTO/Value Object Structs

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

## File Names

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

## Transport Layer Naming

### Package Name

```go
package transporthealthcheck  // prefix 'transport' + feature name
package transportuser
package transportproduct
```

### Struct Names

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

### Constructor Names

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

### Method Names

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

### File Organization

```text
internal/transport/healthcheck/
├── grpc_healthcheck.go      # Contains: TransportHealthCheckGrpc
└── restapi_healthcheck.go   # Contains: TransportHealthCheckRestApi

internal/transport/user/
├── grpc_user.go             # Contains: TransportUserGrpc
└── restapi_user.go          # Contains: TransportUserRestApi
```

## Worker Layer Naming

### Package Name

```go
package workerhealthcheck  // prefix 'worker' + feature name
package workeruser
package workerproduct
```

### Struct Names

Use descriptive names with `Scheduler` prefix for cron jobs:

```go
type SchedulerHealthCheck struct {
    healthcheckService domainhealthcheck.HealthCheckService
}

type SchedulerUser struct {
    userService domainuser.UserService
}
```

### Constructor Names

```go
func NewSchedulerHealthCheck(
    healthcheckService domainhealthcheck.HealthCheckService,
) *SchedulerHealthCheck {
    return &SchedulerHealthCheck{
        healthcheckService: healthcheckService,
    }
}
```

### Method Names

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

### File Names

```text
internal/worker/healthcheck/
└── scheduler_healthcheck.go  # Contains: SchedulerHealthCheck

internal/worker/user/
├── scheduler_user.go         # Contains: SchedulerUser
└── consumer_user.go          # Contains: ConsumerUser (future - for Kafka/RabbitMQ)
```

## Method Names

### Service Methods

Use descriptive verb-noun combinations:

```go
CheckDependencies(ctx context.Context) (output CheckDependenciesOutput)
CreateUser(ctx context.Context, input CreateUserInput) (output CreateUserOutput, err error)
GetUserByID(ctx context.Context, userID string) (output UserOutput, err error)
UpdateUserStatus(ctx context.Context, userID string, status UserStatus) (err error)
DeleteUser(ctx context.Context, userID string) (err error)
```

### Repository Methods

Use simple CRUD operations or specific queries:

```go
Create(ctx context.Context, user User) (err error)
FindByID(ctx context.Context, id string) (user User, err error)
FindByEmail(ctx context.Context, email string) (user User, err error)
Update(ctx context.Context, user User) (err error)
Delete(ctx context.Context, id string) (err error)
PingDatabase(ctx context.Context) (responseTime time.Duration, err error)
```

## Constants and Enums

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

## Function Parameters and Return Values

### Named Return Values

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

### Input/Output Pattern

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

## Constructor Functions

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

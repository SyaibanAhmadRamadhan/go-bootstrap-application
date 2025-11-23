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
    - [Transport Package Names](#transport-package-names)
    - [Transport Struct Names](#transport-struct-names)
    - [Transport Constructor Names](#transport-constructor-names)
    - [Transport Method Names](#transport-method-names)
    - [Transport File Organization](#transport-file-organization)
  - [Worker Layer Naming](#worker-layer-naming)
    - [Worker Package Names](#worker-package-names)
    - [Worker Struct Names](#worker-struct-names)
    - [Worker Constructor Names](#worker-constructor-names)
    - [Worker Method Names](#worker-method-names)
    - [Worker File Names](#worker-file-names)
  - [Repository Layer Naming Convention](#repository-layer-naming-convention)
    - [Repository Parameter Naming](#repository-parameter-naming)
    - [Repository Return Value Naming](#repository-return-value-naming)
    - [Repository Key Principles](#repository-key-principles)
  - [Service Layer Naming Convention](#service-layer-naming-convention)
    - [Service Input Naming](#service-input-naming)
    - [Service Output Naming](#service-output-naming)
    - [Service Key Principles](#service-key-principles)
  - [Method Independence Principle](#method-independence-principle)
    - [Bad Example - Method Depends on Struct from Another Method](#bad-example---method-depends-on-struct-from-another-method)
    - [Good Example - Each Method Has Independent Structs](#good-example---each-method-has-independent-structs)
    - [Why Method Independence?](#why-method-independence)
    - [Independence Rules](#independence-rules)
    - [Complete Flow Example](#complete-flow-example)
  - [Method Names](#method-names)
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

Folder:  internal/module/auth/repository/
Package: package authrepository

Folder:  internal/module/user/repository/
Package: package userrepository

# Transport Layer
Folder:  internal/transport/healthcheck/
Package: package transporthealthcheck

Folder:  internal/transport/auth/
Package: package transportauth

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
package authrepository       // module repository implementation
package transporthealthcheck // transport layer
package transportauth        // transport layer
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

### Transport Package Names

```go
package transporthealthcheck  // prefix 'transport' + feature name
package transportuser
package transportproduct
```

### Transport Struct Names

Use descriptive names with feature name and handler type suffix:

```go
// REST API Transport
type HealthCheckRestAPIHandler struct {
    healthcheckService domainhealthcheck.HealthCheckService
    helper *ginx.GinHelper
}

type AuthRestAPIHandler struct {
    authService domainauth.AuthService
    helper *ginx.GinHelper
}

type UserRestAPIHandler struct {
    userService domainuser.UserService
    helper *ginx.GinHelper
}

// gRPC Transport
type HealthCheckGrpcHandler struct {
    healthcheckService domainhealthcheck.HealthCheckService
    healthcheck.UnimplementedHealthCheckServiceServer
}

type AuthGrpcHandler struct {
    authService domainauth.AuthService
    auth.UnimplementedAuthServiceServer
}

type UserGrpcHandler struct {
    userService domainuser.UserService
    user.UnimplementedUserServiceServer
}
```

### Transport Constructor Names

```go
// REST API Handlers
func NewRestAPIHandler(
    authService domainauth.AuthService,
    helper *ginx.GinHelper,
) *AuthRestAPIHandler {
    return &AuthRestAPIHandler{
        authService: authService,
        helper: helper,
    }
}

func NewRestAPIHandler(
    userService domainuser.UserService,
    helper *ginx.GinHelper,
) *UserRestAPIHandler {
    return &UserRestAPIHandler{
        userService: userService,
        helper: helper,
    }
}

// gRPC Handlers
func NewGrpcHandler(
    healthcheckService domainhealthcheck.HealthCheckService,
) *HealthCheckGrpcHandler {
    return &HealthCheckGrpcHandler{
        healthcheckService: healthcheckService,
    }
}
```

### Transport Method Names

Methods should match the generated API contract:

```go
// REST API - matches OpenAPI operation IDs
func (h *AuthRestAPIHandler) ApiV1PostAuthLogin(c *gin.Context) {
    // Implementation
}

func (h *UserRestAPIHandler) ApiV1PostUsersRegister(c *gin.Context) {
    // Implementation
}

func (h *UserRestAPIHandler) ApiV1GetUsers(
    c *gin.Context,
    params restapigen.ApiV1GetUsersParams,
) {
    // Implementation
}

// gRPC - matches protobuf service definition
func (h *HealthCheckGrpcHandler) ApiV1HealthCheck(
    ctx context.Context,
    req *healthcheck.ApiV1HealthCheckRequest,
) (*healthcheck.ApiV1HealthCheckResponse, error) {
    // Implementation
}

func (h *UserGrpcHandler) ApiV1CreateUser(
    ctx context.Context,
    req *user.ApiV1CreateUserRequest,
) (*user.ApiV1CreateUserResponse, error) {
    // Implementation
}
```

### Transport File Organization

```text
internal/transport/healthcheck/
├── grpc_healthcheck.go      # Contains: HealthCheckGrpcHandler
└── restapi_healthcheck.go   # Contains: HealthCheckRestAPIHandler

internal/transport/auth/
├── grpc_auth.go             # Contains: AuthGrpcHandler
└── restapi_auth.go          # Contains: AuthRestAPIHandler

internal/transport/user/
├── grpc_user.go             # Contains: UserGrpcHandler
└── restapi_user.go          # Contains: UserRestAPIHandler
```

## Worker Layer Naming

### Worker Package Names

```go
package workerhealthcheck  // prefix 'worker' + feature name
package workeruser
package workerproduct
```

### Worker Struct Names

Use descriptive names with `Scheduler` prefix for cron jobs:

```go
type SchedulerHealthCheck struct {
    healthcheckService domainhealthcheck.HealthCheckService
}

type SchedulerUser struct {
    userService domainuser.UserService
}
```

### Worker Constructor Names

```go
func NewSchedulerHealthCheck(
    healthcheckService domainhealthcheck.HealthCheckService,
) *SchedulerHealthCheck {
    return &SchedulerHealthCheck{
        healthcheckService: healthcheckService,
    }
}
```

### Worker Method Names

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

### Worker File Names

```text
internal/worker/healthcheck/
└── scheduler_healthcheck.go  # Contains: SchedulerHealthCheck

internal/worker/user/
├── scheduler_user.go         # Contains: SchedulerUser
└── consumer_user.go          # Contains: ConsumerUser (future - for Kafka/RabbitMQ)
```

## Repository Layer Naming Convention

### Repository Parameter Naming

**Pattern:** `<MethodName><Entity>Params` or `<MethodName><Entity>Filters`

```go
// Write operations use "Params"
type CreateUserParams struct {
    Email        string
    PasswordHash string
    Name         string
    Role         UserRole
    Phone        *string
    Gender       *Gender
}

type UpdateUserParams struct {
    UserID  string
    Name    *string
    Phone   *string
    Gender  *Gender
}

type RevokeTokenParams struct {
    Token  string
    UserID string
}

// Read operations use "Filters"
type GetDetailUserFilters struct {
    UserID *string
    Email  *string
}

type GetListUserFilters struct {
    Page     int64
    PageSize int64
    Search   *string
    Status   *UserStatus
    Role     *UserRole
}

type DeleteExpiredTokensParams struct {
    BeforeDate time.Time
}
```

### Repository Return Value Naming

**Pattern:** `<MethodName><Entity>Result`

```go
type CreateUserResult struct {
    ID        string
    Email     string
    Name      string
    Role      UserRole
    Status    UserStatus
    CreatedAt time.Time
}

type GetDetailUserResult struct {
    ID           string
    Email        string
    PasswordHash string
    Name         string
    Role         UserRole
    Status       UserStatus
    Phone        *string
    Gender       *Gender
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type GetListUserResult struct {
    Users      []GetDetailUserResult
    TotalCount int64
}

type RevokeTokenResult struct {
    Success   bool
    RevokedAt time.Time
}

type DeleteExpiredTokensResult struct {
    DeletedCount int64
}
```

### Repository Key Principles

1. **Each method has its own parameter/result struct** - No reusing structs between methods
2. **Write operations use `Params`** - For Create, Update, Delete operations
3. **Read operations use `Filters`** - For Get/Query operations
4. **Return values use `Result`** - Always use Result suffix for consistency
5. **Pointer for optional fields** - Use `*string`, `*int64` for optional filters/params

## Service Layer Naming Convention

### Service Input Naming

**Pattern:** `<MethodName>Input`

```go
type RegisterInput struct {
    Email    string
    Password string
    Name     string
    Phone    *string
    Gender   *Gender
}

type LoginInput struct {
    Email    string
    Password string
}

type GetProfileInput struct {
    UserID string
}

type UpdateProfileInput struct {
    UserID string
    Name   *string
    Phone  *string
    Gender *Gender
}

type ChangePasswordInput struct {
    UserID      string
    OldPassword string
    NewPassword string
}

type GetListInput struct {
    Page     int64
    PageSize int64
    Search   *string
    Status   *UserStatus
    Role     *UserRole
}
```

### Service Output Naming

**Pattern:** `<MethodName>Output`

```go
type RegisterOutput struct {
    UserID    string
    Email     string
    Name      string
    Role      UserRole
    Status    UserStatus
    CreatedAt time.Time
}

type LoginOutput struct {
    AccessToken  string
    RefreshToken string
    ExpiresIn    int64
    TokenType    string
    User         UserInfo
}

type GetProfileOutput struct {
    ID        string
    Email     string
    Name      string
    Role      UserRole
    Status    UserStatus
    Phone     *string
    Gender    *Gender
    CreatedAt time.Time
    UpdatedAt time.Time
}

type UpdateProfileOutput struct {
    UpdatedAt time.Time
}

type ChangePasswordOutput struct {
    Success   bool
    UpdatedAt time.Time
}

type GetListOutput struct {
    Users      []UserInfo
    TotalCount int64
    Page       int64
    PageSize   int64
}
```

### Service Key Principles

1. **Each method has its own Input/Output struct** - No reusing structs between methods
2. **Input for request data** - Always use Input suffix
3. **Output for response data** - Always use Output suffix
4. **Business logic friendly** - Contains only business-relevant fields
5. **Different from Repository layer** - Service DTOs are business-focused, Repository DTOs are data-focused

## Method Independence Principle

### Bad Example - Method Depends on Struct from Another Method

```go
// WRONG: Reusing CreateUserParams for Update
func (r *repository) UpdateUser(ctx context.Context, params CreateUserParams) error {
    // This creates coupling and confusion
}

// WRONG: Reusing LoginInput for multiple purposes
func (s *service) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
    // Implementation
}

func (s *service) RefreshToken(ctx context.Context, input LoginInput) (LoginOutput, error) {
    // WRONG: LoginInput doesn't make sense for RefreshToken
}
```

### Good Example - Each Method Has Independent Structs

```go
// Repository Layer - Independent structs per method
type CreateUserParams struct {
    Email        string
    PasswordHash string
    Name         string
    Role         UserRole
}

type UpdateUserParams struct {
    UserID string
    Name   *string
    Phone  *string
}

func (r *repository) CreateUser(ctx context.Context, params CreateUserParams) (CreateUserResult, error) {
    // Implementation
}

func (r *repository) UpdateUser(ctx context.Context, params UpdateUserParams) error {
    // Implementation
}

// Service Layer - Independent structs per method
type LoginInput struct {
    Email    string
    Password string
}

type RefreshTokenInput struct {
    RefreshToken string
}

type LoginOutput struct {
    AccessToken  string
    RefreshToken string
    ExpiresIn    int64
}

type RefreshTokenOutput struct {
    AccessToken  string
    RefreshToken string
    ExpiresIn    int64
}

func (s *service) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
    // Implementation
}

func (s *service) RefreshToken(ctx context.Context, input RefreshTokenInput) (RefreshTokenOutput, error) {
    // Implementation with its own struct
}
```

### Why Method Independence?

1. **Clear Intent** - Each method's purpose is obvious from its parameters
2. **Easy to Change** - Modifying one method doesn't affect others
3. **No Confusion** - Developers know exactly what data is needed
4. **Better Validation** - Each struct can have specific validation rules
5. **Self-Documenting** - Code is easier to understand and maintain

### Independence Rules

- ✅ **DO**: Create unique param/result structs for each method
- ✅ **DO**: Use descriptive names that match the method name
- ✅ **DO**: Only include fields relevant to that specific operation
- ❌ **DON'T**: Reuse structs between different methods
- ❌ **DON'T**: Use generic names like `Request`, `Response`, `Data`
- ❌ **DON'T**: Share structs between repository and service layers

### Complete Flow Example

```go
// Domain Layer - Repository Interface
type UserRepositoryDatastore interface {
    CreateUser(ctx context.Context, params CreateUserParams) (CreateUserResult, error)
    GetDetailUser(ctx context.Context, filters GetDetailUserFilters) (GetDetailUserResult, error)
    UpdateUser(ctx context.Context, params UpdateUserParams) error
}

// Domain Layer - Service Interface
type UserService interface {
    Register(ctx context.Context, input RegisterInput) (RegisterOutput, error)
    GetProfile(ctx context.Context, input GetProfileInput) (GetProfileOutput, error)
    UpdateProfile(ctx context.Context, input UpdateProfileInput) (UpdateProfileOutput, error)
}

// Module/Repository - Implementation
func (r *repository) CreateUser(ctx context.Context, params CreateUserParams) (CreateUserResult, error) {
    // Each method has its own params and result
    // Independent from other methods
}

// Module/Service - Implementation
func (s *service) Register(ctx context.Context, input RegisterInput) (RegisterOutput, error) {
    // Transform service input to repository params
    repoParams := CreateUserParams{
        Email:        input.Email,
        PasswordHash: hashedPassword,
        Name:         input.Name,
        Role:         UserRoleUser,
    }
    
    result, err := s.userRepo.CreateUser(ctx, repoParams)
    
    // Transform repository result to service output
    return RegisterOutput{
        UserID:    result.ID,
        Email:     result.Email,
        Name:      result.Name,
        Role:      result.Role,
        Status:    result.Status,
        CreatedAt: result.CreatedAt,
    }, nil
}
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
// Basic CRUD
Create(ctx context.Context, params CreateUserParams) (result CreateUserResult, err error)
FindByID(ctx context.Context, id string) (user User, err error)
FindByEmail(ctx context.Context, email string) (user User, err error)
Update(ctx context.Context, params UpdateUserParams) (err error)
Delete(ctx context.Context, id string) (err error)

// Query with filters
GetDetailUser(ctx context.Context, filters GetDetailUserFilters) (result GetDetailUserResult, err error)
GetListUser(ctx context.Context, filters GetListUserFilters) (result GetListUserResult, err error)

// Specific operations
PingDatabase(ctx context.Context) (responseTime time.Duration, err error)
CreateToken(ctx context.Context, params CreateTokenParams) (result CreateTokenResult, err error)
RevokeToken(ctx context.Context, params RevokeTokenParams) (result RevokeTokenResult, err error)
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

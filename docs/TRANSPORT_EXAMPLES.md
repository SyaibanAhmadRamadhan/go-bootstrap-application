# Transport Layer Examples

This document provides detailed examples for implementing HTTP REST and gRPC transport handlers.

## Table of Contents

- [gRPC Transport](#grpc-transport)
- [REST API Transport](#rest-api-transport)
- [Error Handling](#error-handling)
- [Best Practices](#best-practices)

## gRPC Transport

### Basic gRPC Handler Example

```go
package transporthealthcheck

import (
    "context"
    domainhealthcheck "project/internal/domain/healthcheck"
    "project/internal/gen/grpcgen/healthcheck"
    "google.golang.org/protobuf/types/known/timestamppb"
)

type HealthCheckGrpcHandler struct {
    healthcheckService domainhealthcheck.HealthCheckService
    healthcheck.UnimplementedHealthCheckServiceServer
}

func NewGrpcHandler(
    healthcheckService domainhealthcheck.HealthCheckService,
) *HealthCheckGrpcHandler {
    return &HealthCheckGrpcHandler{
        healthcheckService: healthcheckService,
    }
}

func (h *HealthCheckGrpcHandler) ApiV1HealthCheck(
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

// Helper function to map domain status to protobuf enum
func mapStatusToProto(status domainhealthcheck.StatusHealthCheck) healthcheck.HealthStatus {
    switch status {
    case domainhealthcheck.StatusHealthCheckHealthy:
        return healthcheck.HealthStatus_HEALTHY
    case domainhealthcheck.StatusHealthCheckDegraded:
        return healthcheck.HealthStatus_DEGRADED
    case domainhealthcheck.StatusHealthCheckUnhealthy:
        return healthcheck.HealthStatus_UNHEALTHY
    default:
        return healthcheck.HealthStatus_UNKNOWN
    }
}

func mapDependencyStatusToProto(status domainhealthcheck.StatusDependency) healthcheck.DependencyStatus {
    switch status {
    case domainhealthcheck.StatusDependencyHealthy:
        return healthcheck.DependencyStatus_DEPENDENCY_HEALTHY
    case domainhealthcheck.StatusDependencyUnhealthy:
        return healthcheck.DependencyStatus_DEPENDENCY_UNHEALTHY
    default:
        return healthcheck.DependencyStatus_DEPENDENCY_UNKNOWN
    }
}
```

### Advanced gRPC Handler with Error Handling

```go
package transportuser

import (
    "context"
    "errors"
    domainuser "project/internal/domain/user"
    "project/internal/gen/grpcgen/user"
    
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/types/known/timestamppb"
)

type UserGrpcHandler struct {
    userService domainuser.UserService
    user.UnimplementedUserServiceServer
}

func NewGrpcHandler(
    userService domainuser.UserService,
) *UserGrpcHandler {
    return &UserGrpcHandler{
        userService: userService,
    }
}

func (h *UserGrpcHandler) ApiV1CreateUser(
    ctx context.Context,
    req *user.ApiV1CreateUserRequest,
) (*user.ApiV1CreateUserResponse, error) {
    // Transform request to domain input
    input := domainuser.CreateUserInput{
        Name:  req.Name,
        Email: req.Email,
    }

    // Call domain service
    output, err := t.userService.CreateUser(ctx, input)
    if err != nil {
        // Map domain errors to gRPC status codes
        return nil, mapErrorToGrpcStatus(err)
    }

    // Transform domain output to protobuf response
    return &user.ApiV1CreateUserResponse{
        UserId:    output.UserID,
        Name:      output.Name,
        Email:     output.Email,
        CreatedAt: timestamppb.New(output.CreatedAt),
    }, nil
}

func (h *UserGrpcHandler) ApiV1GetUser(
    ctx context.Context,
    req *user.ApiV1GetUserRequest,
) (*user.ApiV1GetUserResponse, error) {
    // Call domain service
    output, err := h.userService.GetUserByID(ctx, req.UserId)
    if err != nil {
        return nil, mapErrorToGrpcStatus(err)
    }

    // Transform domain output to protobuf response
    return &user.ApiV1GetUserResponse{
        UserId:    output.ID,
        Name:      output.Name,
        Email:     output.Email,
        Status:    mapUserStatusToProto(output.Status),
        CreatedAt: timestamppb.New(output.CreatedAt),
    }, nil
}

// Error mapping helper
func mapErrorToGrpcStatus(err error) error {
    switch {
    case errors.Is(err, domainuser.ErrUserNotFound):
        return status.Error(codes.NotFound, "user not found")
    case errors.Is(err, domainuser.ErrUserAlreadyExists):
        return status.Error(codes.AlreadyExists, "user already exists")
    case errors.Is(err, domainuser.ErrInvalidEmail):
        return status.Error(codes.InvalidArgument, "invalid email format")
    case errors.Is(err, domainuser.ErrUnauthorized):
        return status.Error(codes.Unauthenticated, "unauthorized")
    default:
        return status.Error(codes.Internal, "internal server error")
    }
}

func mapUserStatusToProto(s domainuser.UserStatus) user.UserStatus {
    switch s {
    case domainuser.UserStatusActive:
        return user.UserStatus_ACTIVE
    case domainuser.UserStatusInactive:
        return user.UserStatus_INACTIVE
    case domainuser.UserStatusSuspended:
        return user.UserStatus_SUSPENDED
    default:
        return user.UserStatus_UNKNOWN
    }
}
```

## REST API Transport

### Basic REST Handler Example

```go
package transporthealthcheck

import (
    domainhealthcheck "project/internal/domain/healthcheck"
    "project/internal/gen/restapigen"
    
    "github.com/SyaibanAhmadRamadhan/go-foundation-kit/http/server/ginx"
    "github.com/gin-gonic/gin"
)

type HealthCheckRestAPIHandler struct {
    healthcheckService domainhealthcheck.HealthCheckService
    helper *ginx.GinHelper
}

func NewRestAPIHandler(
    healthcheckService domainhealthcheck.HealthCheckService,
    helper *ginx.GinHelper,
) *HealthCheckRestAPIHandler {
    return &HealthCheckRestAPIHandler{
        healthcheckService: healthcheckService,
        helper: helper,
    }
}

func (h *HealthCheckRestAPIHandler) ApiV1GetHealthCheck(
    c *gin.Context,
) {
    // Call domain service
    output := h.healthcheckService.CheckDependencies(c.Request.Context())

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
    render.Status(r, http.StatusOK)
    render.JSON(w, r, resp)
}
```

### Advanced REST Handler with Error Handling

```go
package transportuser

import (
    "errors"
    domainuser "project/internal/domain/user"
    "project/internal/gen/restapigen"
    
    "github.com/SyaibanAhmadRamadhan/go-foundation-kit/http/server/ginx"
    "github.com/gin-gonic/gin"
)

type UserRestAPIHandler struct {
    userService domainuser.UserService
    helper *ginx.GinHelper
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

func (h *UserRestAPIHandler) ApiV1PostUsers(
    c *gin.Context,
) {
    // Parse request body
    var req restapigen.ApiV1PostUsersRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, r, http.StatusBadRequest, "invalid request body")
        return
    }

    // Transform request to domain input
    input := domainuser.CreateUserInput{
        Name:  req.Name,
        Email: req.Email,
    }

    // Call domain service
    output, err := t.userService.CreateUser(r.Context(), input)
    if err != nil {
        writeErrorFromDomain(w, r, err)
        return
    }

    // Transform domain output to REST API response
    resp := restapigen.ApiV1PostUsersResponse{
        UserId:    output.UserID,
        Name:      output.Name,
        Email:     output.Email,
        CreatedAt: output.CreatedAt,
    }

    render.Status(r, http.StatusCreated)
    render.JSON(w, r, resp)
}

func (t *TransportUserRestApi) ApiV1GetUsersUserId(
    w http.ResponseWriter, 
    r *http.Request,
) {
    // Extract path parameter
    userID := chi.URLParam(r, "userId")
    if userID == "" {
        writeError(w, r, http.StatusBadRequest, "user_id is required")
        return
    }

    // Call domain service
    output, err := t.userService.GetUserByID(r.Context(), userID)
    if err != nil {
        writeErrorFromDomain(w, r, err)
        return
    }

    // Transform domain output to REST API response
    resp := restapigen.ApiV1GetUsersUserIdResponse{
        UserId:    output.ID,
        Name:      output.Name,
        Email:     output.Email,
        Status:    restapigen.UserStatus(output.Status),
        CreatedAt: output.CreatedAt,
    }

    render.Status(r, http.StatusOK)
    render.JSON(w, r, resp)
}

func (t *TransportUserRestApi) ApiV1PatchUsersUserId(
    w http.ResponseWriter, 
    r *http.Request,
) {
    userID := chi.URLParam(r, "userId")
    
    var req restapigen.ApiV1PatchUsersUserIdRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, r, http.StatusBadRequest, "invalid request body")
        return
    }

    input := domainuser.UpdateUserInput{
        UserID: userID,
        Name:   req.Name,
        Email:  req.Email,
    }

    output, err := t.userService.UpdateUser(r.Context(), input)
    if err != nil {
        writeErrorFromDomain(w, r, err)
        return
    }

    resp := restapigen.ApiV1PatchUsersUserIdResponse{
        UserId:    output.ID,
        Name:      output.Name,
        Email:     output.Email,
        UpdatedAt: output.UpdatedAt,
    }

    render.Status(r, http.StatusOK)
    render.JSON(w, r, resp)
}

func (t *TransportUserRestApi) ApiV1DeleteUsersUserId(
    w http.ResponseWriter, 
    r *http.Request,
) {
    userID := chi.URLParam(r, "userId")
    
    err := t.userService.DeleteUser(r.Context(), userID)
    if err != nil {
        writeErrorFromDomain(w, r, err)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
```

## Error Handling

### Error Response Helpers

```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
    Code    string `json:"code,omitempty"`
}

func writeError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
    resp := ErrorResponse{
        Error:   http.StatusText(statusCode),
        Message: message,
    }
    render.Status(r, statusCode)
    render.JSON(w, r, resp)
}

func writeErrorFromDomain(w http.ResponseWriter, r *http.Request, err error) {
    var statusCode int
    var message string

    switch {
    case errors.Is(err, domainuser.ErrUserNotFound):
        statusCode = http.StatusNotFound
        message = "user not found"
    case errors.Is(err, domainuser.ErrUserAlreadyExists):
        statusCode = http.StatusConflict
        message = "user already exists"
    case errors.Is(err, domainuser.ErrInvalidEmail):
        statusCode = http.StatusBadRequest
        message = "invalid email format"
    case errors.Is(err, domainuser.ErrUnauthorized):
        statusCode = http.StatusUnauthorized
        message = "unauthorized"
    default:
        statusCode = http.StatusInternalServerError
        message = "internal server error"
    }

    writeError(w, r, statusCode, message)
}
```

### Domain Error Definitions

```go
package domainuser

import "errors"

var (
    ErrUserNotFound      = errors.New("user not found")
    ErrUserAlreadyExists = errors.New("user already exists")
    ErrInvalidEmail      = errors.New("invalid email format")
    ErrUnauthorized      = errors.New("unauthorized")
)
```

## Best Practices

### 1. Keep Transport Layer Thin

**DO:**

```go
func (t *TransportUserRestApi) ApiV1PostUsers(w http.ResponseWriter, r *http.Request) {
    var req restapigen.ApiV1PostUsersRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    input := domainuser.CreateUserInput{
        Name:  req.Name,
        Email: req.Email,
    }
    
    output, err := t.userService.CreateUser(r.Context(), input)
    // Handle response...
}
```

**DON'T:**

```go
func (t *TransportUserRestApi) ApiV1PostUsers(w http.ResponseWriter, r *http.Request) {
    var req restapigen.ApiV1PostUsersRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // DON'T: Business logic in transport layer
    if req.Email == "" {
        writeError(w, r, 400, "email required")
        return
    }
    if !strings.Contains(req.Email, "@") {
        writeError(w, r, 400, "invalid email")
        return
    }
    
    // DON'T: Direct database access
    user, err := t.db.FindUserByEmail(req.Email)
    // ...
}
```

### 2. Use Consistent Error Mapping

Create helper functions to map domain errors to transport-specific errors:

```go
// For gRPC
func mapToGrpcStatus(err error) error { /* ... */ }

// For REST API
func mapToHttpStatus(err error) int { /* ... */ }
```

### 3. Transform Data Properly

Always transform between protocol types and domain types:

```go
// Request → Domain Input
input := domainuser.CreateUserInput{
    Name:  req.Name,
    Email: req.Email,
}

// Domain Output → Response
resp := restapigen.ApiV1PostUsersResponse{
    UserId: output.UserID,
    Name:   output.Name,
    Email:  output.Email,
}
```

### 4. Handle Context Properly

Pass request context to domain services:

```go
// gRPC - context from request
output, err := t.userService.CreateUser(ctx, input)

// REST API - context from http.Request
output, err := t.userService.CreateUser(r.Context(), input)
```

### 5. Use Structured Responses

```go
// Good: Structured JSON response
render.JSON(w, r, restapigen.ApiV1GetUserResponse{
    UserId: user.ID,
    Name:   user.Name,
    Email:  user.Email,
})

// Bad: Plain text or unstructured response
fmt.Fprintf(w, "User: %s, Email: %s", user.Name, user.Email)
```

### 6. Validate Early, Delegate Logic

```go
func (t *TransportUserRestApi) ApiV1PostUsers(w http.ResponseWriter, r *http.Request) {
    // Validate request format (transport concern)
    var req restapigen.ApiV1PostUsersRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, r, 400, "invalid request body")
        return
    }
    
    // Transform and delegate to service (business logic)
    input := domainuser.CreateUserInput{Name: req.Name, Email: req.Email}
    output, err := t.userService.CreateUser(r.Context(), input)
    // ...
}
```

### 7. Log Appropriately

```go
import "log/slog"

func (t *TransportUserRestApi) ApiV1PostUsers(w http.ResponseWriter, r *http.Request) {
    slog.Info("Creating user", "name", req.Name)
    
    output, err := t.userService.CreateUser(r.Context(), input)
    if err != nil {
        slog.Error("Failed to create user", "error", err, "name", req.Name)
        writeErrorFromDomain(w, r, err)
        return
    }
    
    slog.Info("User created successfully", "user_id", output.UserID)
    // ...
}
```

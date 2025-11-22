//go:generate go tool mockgen -source=repository.go -destination=../../gen/mockgen/user_repository_mock.gen.go -package=mockgen

package domainuser

import (
	"context"
	"time"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/utils/primitive"
)

// UserRepositoryDatastore - for database operations related to users
type UserRepositoryDatastore interface {
	// CreateUser creates a new user record
	CreateUser(ctx context.Context, params CreateUserParams) (CreateUserResult, error)

	// GetDetailUser retrieves user details
	GetDetailUser(ctx context.Context, filters GetDetailUserFilters) (GetDetailUserResult, error)

	// GetListUser retrieves list of users with pagination
	GetListUser(ctx context.Context, filters GetListUserFilters) (GetListUserResult, error)

	// UpdateUser updates user information
	UpdateUser(ctx context.Context, params UpdateUserParams) (UpdateUserResult, error)

	// UpdatePassword updates user password
	UpdatePassword(ctx context.Context, params UpdatePasswordParams) (UpdatePasswordResult, error)

	// UpdateStatus updates user status
	UpdateStatus(ctx context.Context, params UpdateStatusParams) (UpdateStatusResult, error)
}

// ============================================
// CreateUser - Create new user
// ============================================

type CreateUserParams struct {
	Email        string
	PasswordHash string
	Name         string
	Role         UserRole
	Phone        *string
	Gender       *Gender
}

type CreateUserResult struct {
	ID        string
	Email     string
	Name      string
	Role      UserRole
	Status    UserStatus
	CreatedAt time.Time
}

// ============================================
// GetDetailUser - Get user details
// ============================================

type GetDetailUserFilters struct {
	UserID *string
	Email  *string
}

type GetDetailUserResult struct {
	ID           string
	Email        string
	PasswordHash string // for authentication
	Name         string
	Role         UserRole
	Status       UserStatus
	Phone        *string
	Gender       *Gender
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ============================================
// GetListUser - Get list of users
// ============================================

type GetListUserFilters struct {
	Pagination primitive.PaginationInput
	Search     *string // search by name or email
	Status     *UserStatus
	Role       *UserRole
}

type GetListUserResult struct {
	Users      []GetDetailUserResult
	Pagination primitive.PaginationOutput
}

// ============================================
// UpdateUser - Update user profile
// ============================================

type UpdateUserParams struct {
	UserID string
	Name   *string
	Phone  *string
	Gender *Gender
}

type UpdateUserResult struct {
	UpdatedAt time.Time
}

// ============================================
// UpdatePassword - Update user password
// ============================================

type UpdatePasswordParams struct {
	UserID          string
	NewPasswordHash string
}

type UpdatePasswordResult struct {
	UpdatedAt time.Time
}

// ============================================
// UpdateStatus - Update user status
// ============================================

type UpdateStatusParams struct {
	UserID string
	Status UserStatus
}

type UpdateStatusResult struct {
	UpdatedAt time.Time
}

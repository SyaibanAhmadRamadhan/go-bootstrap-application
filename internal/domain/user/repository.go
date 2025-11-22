//go:generate go tool mockgen -source=repository.go -destination=../../gen/mockgen/user_repository_mock.gen.go -package=mockgen

package domainuser

import (
	"context"
	sharedkernel "go-bootstrap/internal/domain/shared"
	"time"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/utils/primitive"
)

type UserRepositoryDatastore interface {
	CreateUser(ctx context.Context, params CreateUserParams) (CreateUserResult, error)

	GetDetailUser(ctx context.Context, filters GetDetailUserFilters) (GetDetailUserResult, error)

	GetListUser(ctx context.Context, filters GetListUserFilters) (GetListUserResult, error)

	UpdateUser(ctx context.Context, params UpdateUserParams) (UpdateUserResult, error)

	UpdatePassword(ctx context.Context, params UpdatePasswordParams) (UpdatePasswordResult, error)

	UpdateStatus(ctx context.Context, params UpdateStatusParams) (UpdateStatusResult, error)
}

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
	Status    sharedkernel.UserStatus
	CreatedAt time.Time
}

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
	Status       sharedkernel.UserStatus
	Phone        *string
	Gender       *Gender
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type GetListUserFilters struct {
	Pagination primitive.PaginationInput
	Search     *string // search by name or email
	Status     *sharedkernel.UserStatus
	Role       *UserRole
}

type GetListUserResult struct {
	Users      []GetDetailUserResult
	Pagination primitive.PaginationOutput
}

type UpdateUserParams struct {
	UserID string
	Name   *string
	Phone  *string
	Gender *Gender
}

type UpdateUserResult struct {
	UpdatedAt time.Time
}

type UpdatePasswordParams struct {
	UserID          string
	NewPasswordHash string
}

type UpdatePasswordResult struct {
	UpdatedAt time.Time
}

type UpdateStatusParams struct {
	UserID string
	Status sharedkernel.UserStatus
}

type UpdateStatusResult struct {
	UpdatedAt time.Time
}

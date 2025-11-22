package domainuser

import (
	sharedkernel "go-bootstrap/internal/domain/shared"
	"time"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/utils/primitive"
)

type RegisterInput struct {
	Email    string
	Password string
	Name     string
	Phone    *string
	Gender   *Gender
}

type RegisterOutput struct {
	UserID    string
	Email     string
	Name      string
	CreatedAt time.Time
}

type GetProfileInput struct {
	UserID string
}

type GetProfileOutput struct {
	User User
}

type GetListInput struct {
	Pagination primitive.PaginationInput
	Search     *string
	Status     *sharedkernel.UserStatus
	Role       *UserRole
}

type GetListOutput struct {
	Users      []User
	Pagination primitive.PaginationOutput
}

type UpdateProfileInput struct {
	UserID string
	Name   *string
	Phone  *string
	Gender *Gender
}

type UpdateProfileOutput struct {
	User      User
	UpdatedAt time.Time
}

type ChangePasswordInput struct {
	UserID      string
	OldPassword string
	NewPassword string
}

type ChangePasswordOutput struct {
	Success   bool
	UpdatedAt time.Time
}

type UpdateStatusInput struct {
	UserID string
	Status sharedkernel.UserStatus
}

type UpdateStatusOutput struct {
	Success   bool
	UpdatedAt time.Time
}

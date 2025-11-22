package domainuser

import (
	"time"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/utils/primitive"
)

// ============================================
// Register Operation
// ============================================

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

// ============================================
// GetProfile Operation
// ============================================

type GetProfileInput struct {
	UserID string
}

type GetProfileOutput struct {
	User User
}

// ============================================
// GetList Operation
// ============================================

type GetListInput struct {
	Pagination primitive.PaginationInput
	Search     *string
	Status     *UserStatus
	Role       *UserRole
}

type GetListOutput struct {
	Users      []User
	Pagination primitive.PaginationOutput
}

// ============================================
// UpdateProfile Operation
// ============================================

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

// ============================================
// ChangePassword Operation
// ============================================

type ChangePasswordInput struct {
	UserID      string
	OldPassword string
	NewPassword string
}

type ChangePasswordOutput struct {
	Success   bool
	UpdatedAt time.Time
}

// ============================================
// UpdateStatus Operation (Admin only)
// ============================================

type UpdateStatusInput struct {
	UserID string
	Status UserStatus
}

type UpdateStatusOutput struct {
	Success   bool
	UpdatedAt time.Time
}

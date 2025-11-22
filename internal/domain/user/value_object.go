package domainuser

import (
	sharedkernel "go-bootstrap/internal/domain/shared"
	"time"
)

// User Role
type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

// Gender
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// User Entity - base user information
type User struct {
	ID        string
	Email     string
	Name      string
	Role      UserRole
	Status    sharedkernel.UserStatus
	Gender    *Gender
	Phone     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

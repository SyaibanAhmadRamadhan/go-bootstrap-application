package domainuser

import "time"

// User Status
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
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
	Status    UserStatus
	Gender    *Gender
	Phone     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

package sharedkernel

import "errors"

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

var (
	ErrUserInactive  = errors.New("user account is inactive")
	ErrUserSuspended = errors.New("user account is suspended")
	ErrUserInvalid   = errors.New("invalid user account status")
)

func (s UserStatus) IsActive() bool {
	return s == UserStatusActive
}

func (s UserStatus) IsInactive() bool {
	return s == UserStatusInactive
}

func (s UserStatus) IsSuspended() bool {
	return s == UserStatusSuspended
}

// CanLogin returns bool + domain error
func (s UserStatus) CanLogin() error {
	switch s {
	case UserStatusActive:
		return nil
	case UserStatusInactive:
		return ErrUserInactive
	case UserStatusSuspended:
		return ErrUserSuspended
	default:
		return ErrUserInvalid
	}
}

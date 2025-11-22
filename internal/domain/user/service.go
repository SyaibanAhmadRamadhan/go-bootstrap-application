//go:generate go tool mockgen -source=service.go -destination=../../gen/mockgen/user_service_mock.gen.go -package=mockgen

package domainuser

import "context"

type UserService interface {
	// Register creates a new user account
	Register(ctx context.Context, input RegisterInput) (RegisterOutput, error)

	// GetProfile retrieves user profile information
	GetProfile(ctx context.Context, input GetProfileInput) (GetProfileOutput, error)

	// GetList retrieves list of users with pagination (admin only)
	GetList(ctx context.Context, input GetListInput) (GetListOutput, error)

	// UpdateProfile updates user profile information
	UpdateProfile(ctx context.Context, input UpdateProfileInput) (UpdateProfileOutput, error)

	// ChangePassword changes user password
	ChangePassword(ctx context.Context, input ChangePasswordInput) (ChangePasswordOutput, error)

	// UpdateStatus updates user status (admin only)
	UpdateStatus(ctx context.Context, input UpdateStatusInput) (UpdateStatusOutput, error)
}

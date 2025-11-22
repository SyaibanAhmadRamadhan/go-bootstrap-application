package domainuser

import "context"

type UserService interface {
	Register(ctx context.Context, input RegisterInput) (RegisterOutput, error)

	GetProfile(ctx context.Context, input GetProfileInput) (GetProfileOutput, error)

	GetList(ctx context.Context, input GetListInput) (GetListOutput, error)

	UpdateProfile(ctx context.Context, input UpdateProfileInput) (UpdateProfileOutput, error)

	ChangePassword(ctx context.Context, input ChangePasswordInput) (ChangePasswordOutput, error)

	UpdateStatus(ctx context.Context, input UpdateStatusInput) (UpdateStatusOutput, error)
}

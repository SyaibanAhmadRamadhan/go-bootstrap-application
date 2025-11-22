//go:generate go tool mockgen -source=service.go -destination=../../gen/mockgen/auth_service_mock.gen.go -package=mockgen

package domainauth

import "context"

type AuthService interface {
	Login(ctx context.Context, input LoginInput) (LoginOutput, error)

	RefreshToken(ctx context.Context, input RefreshTokenInput) (RefreshTokenOutput, error)

	Logout(ctx context.Context, input LogoutInput) (LogoutOutput, error)

	ValidateToken(ctx context.Context, input ValidateTokenInput) (ValidateTokenOutput, error)

	RevokeToken(ctx context.Context, input RevokeTokenInput) (RevokeTokenOutput, error)

	WorkerDeleteExpiredTokens(ctx context.Context)
}

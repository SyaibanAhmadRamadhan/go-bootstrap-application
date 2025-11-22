//go:generate go tool mockgen -source=service.go -destination=../../gen/mockgen/auth_service_mock.gen.go -package=mockgen

package domainauth

import "context"

type AuthService interface {
	// Login authenticates user and returns tokens
	Login(ctx context.Context, input LoginInput) (LoginOutput, error)

	// RefreshToken generates new access token from refresh token
	RefreshToken(ctx context.Context, input RefreshTokenInput) (RefreshTokenOutput, error)

	// Logout revokes user tokens
	Logout(ctx context.Context, input LogoutInput) (LogoutOutput, error)

	// ValidateToken checks if token is valid and returns payload
	ValidateToken(ctx context.Context, input ValidateTokenInput) (ValidateTokenOutput, error)

	// RevokeToken manually revokes a token
	RevokeToken(ctx context.Context, input RevokeTokenInput) (RevokeTokenOutput, error)

	WorkerDeleteExpiredTokens(ctx context.Context)
}

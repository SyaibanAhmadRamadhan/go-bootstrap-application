//go:generate go tool mockgen -source=repository.go -destination=../../gen/mockgen/auth_repository_mock.gen.go -package=mockgen

package domainauth

import (
	"context"
	sharedkernel "go-bootstrap/internal/domain/shared"
	"time"
)

type AuthRepositoryDatastore interface {
	CreateToken(ctx context.Context, params CreateTokenParams) (CreateTokenResult, error)

	GetDetailToken(ctx context.Context, filters GetDetailTokenFilters) (GetDetailTokenResult, error)

	RevokeToken(ctx context.Context, params RevokeTokenParams) (RevokeTokenResult, error)

	DeleteExpiredTokens(ctx context.Context, params DeleteExpiredTokensParams) (DeleteExpiredTokensResult, error)
}

type UserRepositoryDatastore interface {
	GetDetailUser(ctx context.Context, filters GetDetailUserFilters) (GetDetailUserResult, error)
}

type CreateTokenParams struct {
	UserID       string
	Token        string
	TokenType    TokenType
	ExpiresAt    time.Time
	RefreshToken *string
}

type CreateTokenResult struct {
	ID        string
	CreatedAt time.Time
}

type GetDetailTokenFilters struct {
	Token     *string
	TokenID   *string
	UserID    *string
	TokenType *TokenType
}

type GetDetailTokenResult struct {
	ID           string
	UserID       string
	Token        string
	TokenType    TokenType
	Status       TokenStatus
	ExpiresAt    time.Time
	CreatedAt    time.Time
	RefreshToken *string
}

type RevokeTokenParams struct {
	Token  string
	UserID string
}

type RevokeTokenResult struct {
	Success   bool
	RevokedAt time.Time
}

type DeleteExpiredTokensParams struct {
	BeforeDate time.Time
}

type DeleteExpiredTokensResult struct {
	DeletedCount int64
}

type GetDetailUserFilters struct {
	UserID *string
	Email  *string
}

type GetDetailUserResult struct {
	ID           string
	Email        string
	PasswordHash string
	Name         string
	Role         UserRole
	Status       sharedkernel.UserStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

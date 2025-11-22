//go:generate go tool mockgen -source=repository.go -destination=../../gen/mockgen/auth_repository_mock.gen.go -package=mockgen

package domainauth

import (
	"context"
	"time"
)

// AuthRepositoryDatastore - for database operations related to tokens
type AuthRepositoryDatastore interface {
	// CreateToken stores a new token record
	CreateToken(ctx context.Context, params CreateTokenParams) (CreateTokenResult, error)

	// GetDetailToken retrieves token details
	GetDetailToken(ctx context.Context, filters GetDetailTokenFilters) (GetDetailTokenResult, error)

	// RevokeToken marks token as revoked
	RevokeToken(ctx context.Context, params RevokeTokenParams) (RevokeTokenResult, error)

	// DeleteExpiredTokens removes expired tokens (for cleanup worker)
	DeleteExpiredTokens(ctx context.Context, params DeleteExpiredTokensParams) (DeleteExpiredTokensResult, error)
}

// AuthRepositoryDatastore - for database operations related to tokens
type UserRepositoryDatastore interface {
	GetDetailUser(ctx context.Context, filters GetDetailUserFilters) (GetDetailUserResult, error)
}

// ============================================
// CreateToken - Store new token
// ============================================

type CreateTokenParams struct {
	UserID       string
	Token        string
	TokenType    TokenType
	ExpiresAt    time.Time
	RefreshToken *string // optional, only for access tokens
}

type CreateTokenResult struct {
	ID        string
	CreatedAt time.Time
}

// ============================================
// GetDetailToken - Get token info
// ============================================

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

// ============================================
// RevokeToken - Mark token as revoked
// ============================================

type RevokeTokenParams struct {
	Token  string
	UserID string
}

type RevokeTokenResult struct {
	Success   bool
	RevokedAt time.Time
}

// ============================================
// DeleteExpiredTokens - Cleanup old tokens
// ============================================

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
	PasswordHash string // for authentication
	Name         string
	Role         UserRole
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

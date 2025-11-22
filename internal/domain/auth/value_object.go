package domainauth

import "time"

// Value Objects for Token Types
type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

// Value Objects for User Role (used in token payload)
type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

// Token Status
type TokenStatus string

const (
	TokenStatusActive  TokenStatus = "active"
	TokenStatusRevoked TokenStatus = "revoked"
	TokenStatusExpired TokenStatus = "expired"
)

// Token Payload - extracted from JWT
type TokenPayload struct {
	UserID    string
	Email     string
	Role      UserRole
	TokenType TokenType
	IssuedAt  time.Time
	ExpiresAt time.Time
}

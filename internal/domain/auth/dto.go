package domainauth

import "time"

// ============================================
// Login Operation
// ============================================

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // in seconds
	TokenType    string
}

// ============================================
// Refresh Token Operation
// ============================================

type RefreshTokenInput struct {
	RefreshToken string
}

type RefreshTokenOutput struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	TokenType    string
}

// ============================================
// Logout Operation
// ============================================

type LogoutInput struct {
	AccessToken  string
	RefreshToken string
}

type LogoutOutput struct {
	Success bool
	Message string
}

// ============================================
// Validate Token Operation
// ============================================

type ValidateTokenInput struct {
	Token string
}

type ValidateTokenOutput struct {
	Valid     bool
	Payload   *TokenPayload
	ExpiresAt time.Time
}

// ============================================
// Revoke Token Operation
// ============================================

type RevokeTokenInput struct {
	Token string
}

type RevokeTokenOutput struct {
	Success bool
	Message string
}

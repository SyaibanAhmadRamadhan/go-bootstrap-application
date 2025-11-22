package domainauth

import "time"

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

type RefreshTokenInput struct {
	RefreshToken string
}

type RefreshTokenOutput struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	TokenType    string
}

type LogoutInput struct {
	AccessToken  string
	RefreshToken string
}

type LogoutOutput struct {
	Success bool
	Message string
}

type ValidateTokenInput struct {
	Token string
}

type ValidateTokenOutput struct {
	Valid     bool
	Payload   *TokenPayload
	ExpiresAt time.Time
}

type RevokeTokenInput struct {
	Token string
}

type RevokeTokenOutput struct {
	Success bool
	Message string
}

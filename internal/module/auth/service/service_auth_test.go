package authservice_test

import (
	"context"
	"testing"

	domainauth "go-bootstrap/internal/domain/auth"
)

func TestService_Login(t *testing.T) {
	t.Skip("TODO: Implement with mocks")
}

func TestService_RefreshToken(t *testing.T) {
	t.Skip("TODO: Implement with mocks")
}

func TestService_ValidateToken(t *testing.T) {
	ctx := context.Background()
	input := domainauth.ValidateTokenInput{
		Token: "test-token",
	}

	_ = ctx
	_ = input

	t.Skip("TODO: Setup mocks and test validation")
}

func TestService_Logout(t *testing.T) {
	t.Skip("TODO: Implement with mocks")
}

func TestService_RevokeToken(t *testing.T) {
	t.Skip("TODO: Implement with mocks")
}

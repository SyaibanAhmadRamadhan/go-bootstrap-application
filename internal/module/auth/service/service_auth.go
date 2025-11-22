package authservice

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log/slog"
	"time"

	domainauth "go-bootstrap/internal/domain/auth"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/apperror"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	authRepo domainauth.AuthRepositoryDatastore
	userRepo domainauth.UserRepositoryDatastore
}

func NewService(
	authRepo domainauth.AuthRepositoryDatastore,
	userRepo domainauth.UserRepositoryDatastore,
) *service {
	return &service{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

func (s *service) Login(ctx context.Context, input domainauth.LoginInput) (domainauth.LoginOutput, error) {
	user, err := s.userRepo.GetDetailUser(ctx, domainauth.GetDetailUserFilters{
		Email: &input.Email,
	})
	if err != nil {
		return domainauth.LoginOutput{}, apperror.BadRequest("invalid email or password")
	}

	if err = user.Status.CanLogin(); err != nil {
		return domainauth.LoginOutput{}, apperror.BadRequest(err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return domainauth.LoginOutput{}, apperror.BadRequest("invalid email or password")
	}

	accessToken, err := s.generateToken()
	if err != nil {
		return domainauth.LoginOutput{}, apperror.StdUnknown(err)
	}

	refreshToken, err := s.generateToken()
	if err != nil {
		return domainauth.LoginOutput{}, apperror.StdUnknown(err)
	}

	accessTokenExpiry := time.Now().UTC().Add(15 * time.Minute)
	refreshTokenExpiry := time.Now().UTC().Add(7 * 24 * time.Hour)

	_, err = s.authRepo.CreateToken(ctx, domainauth.CreateTokenParams{
		UserID:       user.ID,
		Token:        accessToken,
		TokenType:    domainauth.TokenTypeAccess,
		ExpiresAt:    accessTokenExpiry,
		RefreshToken: &refreshToken,
	})
	if err != nil {
		return domainauth.LoginOutput{}, apperror.StdUnknown(err)
	}

	_, err = s.authRepo.CreateToken(ctx, domainauth.CreateTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		TokenType: domainauth.TokenTypeRefresh,
		ExpiresAt: refreshTokenExpiry,
	})
	if err != nil {
		return domainauth.LoginOutput{}, apperror.StdUnknown(err)
	}

	return domainauth.LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(15 * 60),
		TokenType:    "Bearer",
	}, nil
}

func (s *service) RefreshToken(ctx context.Context, input domainauth.RefreshTokenInput) (domainauth.RefreshTokenOutput, error) {
	tokenData, err := s.authRepo.GetDetailToken(ctx, domainauth.GetDetailTokenFilters{
		Token: &input.RefreshToken,
	})
	if err != nil {
		if errors.Is(err, databases.ErrNoRowFound) {
			return domainauth.RefreshTokenOutput{}, apperror.BadRequest("invalid refresh token")
		}
		return domainauth.RefreshTokenOutput{}, apperror.StdUnknown(err)
	}

	if tokenData.TokenType != domainauth.TokenTypeRefresh {
		return domainauth.RefreshTokenOutput{}, apperror.BadRequest("invalid token type")
	}

	if tokenData.Status != domainauth.TokenStatusActive {
		return domainauth.RefreshTokenOutput{}, apperror.BadRequest("token is not active")
	}

	if time.Now().UTC().After(tokenData.ExpiresAt) {
		return domainauth.RefreshTokenOutput{}, apperror.BadRequest("refresh token expired")
	}

	newAccessToken, err := s.generateToken()
	if err != nil {
		return domainauth.RefreshTokenOutput{}, apperror.StdUnknown(err)
	}

	newRefreshToken, err := s.generateToken()
	if err != nil {
		return domainauth.RefreshTokenOutput{}, apperror.StdUnknown(err)
	}

	_, _ = s.authRepo.RevokeToken(ctx, domainauth.RevokeTokenParams{
		Token:  input.RefreshToken,
		UserID: tokenData.UserID,
	})

	accessTokenExpiry := time.Now().UTC().Add(15 * time.Minute)
	refreshTokenExpiry := time.Now().UTC().Add(7 * 24 * time.Hour)

	_, err = s.authRepo.CreateToken(ctx, domainauth.CreateTokenParams{
		UserID:       tokenData.UserID,
		Token:        newAccessToken,
		TokenType:    domainauth.TokenTypeAccess,
		ExpiresAt:    accessTokenExpiry,
		RefreshToken: &newRefreshToken,
	})
	if err != nil {
		return domainauth.RefreshTokenOutput{}, apperror.StdUnknown(err)
	}

	_, err = s.authRepo.CreateToken(ctx, domainauth.CreateTokenParams{
		UserID:    tokenData.UserID,
		Token:     newRefreshToken,
		TokenType: domainauth.TokenTypeRefresh,
		ExpiresAt: refreshTokenExpiry,
	})
	if err != nil {
		return domainauth.RefreshTokenOutput{}, apperror.StdUnknown(err)
	}

	return domainauth.RefreshTokenOutput{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(15 * 60),
		TokenType:    "Bearer",
	}, nil
}

func (s *service) Logout(ctx context.Context, input domainauth.LogoutInput) (domainauth.LogoutOutput, error) {
	tokenData, err := s.authRepo.GetDetailToken(ctx, domainauth.GetDetailTokenFilters{
		Token: &input.AccessToken,
	})
	if err != nil {
		return domainauth.LogoutOutput{
			Success: false, Message: "Invalid token",
		}, nil
	}

	_, err = s.authRepo.RevokeToken(ctx, domainauth.RevokeTokenParams{
		Token:  input.AccessToken,
		UserID: tokenData.UserID,
	})
	if err != nil {
		return domainauth.LogoutOutput{
			Success: false,
			Message: "Failed to revoke access token",
		}, nil
	}

	if input.RefreshToken != "" {
		_, _ = s.authRepo.RevokeToken(ctx, domainauth.RevokeTokenParams{
			Token:  input.RefreshToken,
			UserID: tokenData.UserID,
		})
	}

	return domainauth.LogoutOutput{
		Success: true,
		Message: "Logged out successfully",
	}, nil
}

func (s *service) ValidateToken(ctx context.Context, input domainauth.ValidateTokenInput) (domainauth.ValidateTokenOutput, error) {
	tokenData, err := s.authRepo.GetDetailToken(ctx, domainauth.GetDetailTokenFilters{
		Token: &input.Token,
	})
	if err != nil {
		return domainauth.ValidateTokenOutput{Valid: false}, nil
	}

	if tokenData.Status != domainauth.TokenStatusActive {
		return domainauth.ValidateTokenOutput{Valid: false}, nil
	}

	if time.Now().UTC().After(tokenData.ExpiresAt) {
		return domainauth.ValidateTokenOutput{Valid: false}, nil
	}

	user, err := s.userRepo.GetDetailUser(ctx, domainauth.GetDetailUserFilters{
		UserID: &tokenData.UserID,
	})
	if err != nil {
		return domainauth.ValidateTokenOutput{Valid: false}, nil
	}

	payload := &domainauth.TokenPayload{
		UserID:    user.ID,
		Email:     user.Email,
		Role:      domainauth.UserRole(user.Role),
		TokenType: tokenData.TokenType,
		IssuedAt:  tokenData.CreatedAt,
		ExpiresAt: tokenData.ExpiresAt,
	}

	return domainauth.ValidateTokenOutput{
		Valid:     true,
		Payload:   payload,
		ExpiresAt: tokenData.ExpiresAt,
	}, nil
}

func (s *service) RevokeToken(ctx context.Context, input domainauth.RevokeTokenInput) (domainauth.RevokeTokenOutput, error) {
	tokenData, err := s.authRepo.GetDetailToken(ctx, domainauth.GetDetailTokenFilters{
		Token: &input.Token,
	})
	if err != nil {
		return domainauth.RevokeTokenOutput{Success: false, Message: "Token not found"}, nil
	}

	result, err := s.authRepo.RevokeToken(ctx, domainauth.RevokeTokenParams{
		Token:  input.Token,
		UserID: tokenData.UserID,
	})
	if err != nil {
		return domainauth.RevokeTokenOutput{Success: false, Message: "Failed to revoke token"}, nil
	}

	if !result.Success {
		return domainauth.RevokeTokenOutput{Success: false, Message: "Token already revoked"}, nil
	}

	return domainauth.RevokeTokenOutput{
		Success: true,
		Message: "Token revoked successfully",
	}, nil
}

func (s *service) WorkerDeleteExpiredTokens(ctx context.Context) {
	beforeDate := time.Now().UTC().Add(-24 * time.Hour)

	result, err := s.authRepo.DeleteExpiredTokens(ctx, domainauth.DeleteExpiredTokensParams{
		BeforeDate: beforeDate,
	})

	if err != nil {
		slog.Error("Failed to cleanup expired tokens",
			"error", err,
			"before_date", beforeDate,
		)
		return
	}

	if result.DeletedCount > 0 {
		slog.Info("Expired tokens cleaned up successfully",
			"deleted_count", result.DeletedCount,
			"before_date", beforeDate,
		)
	} else {
		slog.Info("No expired tokens to cleanup",
			"before_date", beforeDate,
		)
	}
}

func (s *service) generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

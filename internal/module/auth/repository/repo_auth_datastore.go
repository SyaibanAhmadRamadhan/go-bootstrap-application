package authrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	domainauth "go-bootstrap/internal/domain/auth"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases"
)

func (r *repository) CreateToken(ctx context.Context, params domainauth.CreateTokenParams) (domainauth.CreateTokenResult, error) {
	query := `
		INSERT INTO auth_tokens (user_id, token, token_type, expires_at, refresh_token, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	var result domainauth.CreateTokenResult
	err := r.db.RDBMS().QueryRowContext(ctx, query,
		params.UserID,
		params.Token,
		params.TokenType,
		params.ExpiresAt,
		params.RefreshToken,
		domainauth.TokenStatusActive,
		time.Now().UTC(),
	).Scan(&result.ID, &result.CreatedAt)

	if err != nil {
		return domainauth.CreateTokenResult{}, fmt.Errorf("failed to create token: %w", err)
	}

	return result, nil
}

func (r *repository) GetDetailToken(ctx context.Context, filters domainauth.GetDetailTokenFilters) (domainauth.GetDetailTokenResult, error) {
	sq := r.db.Sq().Select(
		"id",
		"user_id",
		"token",
		"token_type",
		"status",
		"expires_at",
		"created_at",
		"refresh_token",
	).From("auth_tokens")

	if filters.Token != nil {
		sq = sq.Where("token = ?", *filters.Token)
	}

	if filters.TokenID != nil {
		sq = sq.Where("id = ?", *filters.TokenID)
	}

	if filters.UserID != nil {
		sq = sq.Where("user_id = ?", *filters.UserID)
	}

	if filters.TokenType != nil {
		sq = sq.Where("token_type = ?", *filters.TokenType)
	}

	sq = sq.Limit(1)

	var result domainauth.GetDetailTokenResult
	row, err := r.db.RDBMS().QueryRowSq(ctx, sq, false)
	if err != nil {
		return domainauth.GetDetailTokenResult{}, fmt.Errorf("failed to get token: %w", err)
	}

	err = row.Scan(
		&result.ID,
		&result.UserID,
		&result.Token,
		&result.TokenType,
		&result.Status,
		&result.ExpiresAt,
		&result.CreatedAt,
		&result.RefreshToken,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domainauth.GetDetailTokenResult{}, databases.ErrNoRowFound
		}
		return domainauth.GetDetailTokenResult{}, fmt.Errorf("token to scan row: %v", err)
	}

	return result, nil
}

func (r *repository) RevokeToken(ctx context.Context, params domainauth.RevokeTokenParams) (domainauth.RevokeTokenResult, error) {
	query := `
		UPDATE auth_tokens
		SET status = $1, updated_at = $2
		WHERE token = $3 AND user_id = $4
	`

	revokedAt := time.Now().UTC()
	result, err := r.db.RDBMS().ExecContext(ctx, query,
		domainauth.TokenStatusRevoked,
		revokedAt,
		params.Token,
		params.UserID,
	)
	if err != nil {
		return domainauth.RevokeTokenResult{}, fmt.Errorf("failed to revoke token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domainauth.RevokeTokenResult{}, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return domainauth.RevokeTokenResult{
			Success: false,
		}, nil
	}

	return domainauth.RevokeTokenResult{
		Success:   true,
		RevokedAt: revokedAt,
	}, nil
}

func (r *repository) DeleteExpiredTokens(ctx context.Context, params domainauth.DeleteExpiredTokensParams) (domainauth.DeleteExpiredTokensResult, error) {
	query := `
		DELETE FROM auth_tokens
		WHERE expires_at < $1
	`

	result, err := r.db.RDBMS().ExecContext(ctx, query, params.BeforeDate)
	if err != nil {
		return domainauth.DeleteExpiredTokensResult{}, fmt.Errorf("failed to delete expired tokens: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domainauth.DeleteExpiredTokensResult{}, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return domainauth.DeleteExpiredTokensResult{
		DeletedCount: rowsAffected,
	}, nil
}

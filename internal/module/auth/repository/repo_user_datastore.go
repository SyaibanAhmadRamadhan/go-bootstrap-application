package authrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases"

	domainauth "go-bootstrap/internal/domain/auth"
)

func (r *repository) GetDetailUser(ctx context.Context, filters domainauth.GetDetailUserFilters) (domainauth.GetDetailUserResult, error) {
	sq := r.db.Sq().Select(
		"id",
		"email",
		"password_hash",
		"name",
		"role",
		"created_at",
		"updated_at",
	).From("users").Where("status = ?", "active")

	if filters.UserID != nil {
		sq = sq.Where("id = ?", *filters.UserID)
	}

	if filters.Email != nil {
		sq = sq.Where("email = ?", *filters.Email)
	}

	sq = sq.Limit(1)

	var result domainauth.GetDetailUserResult
	row, err := r.db.RDBMS().QueryRowSq(ctx, sq, false)
	if err != nil {
		return domainauth.GetDetailUserResult{}, fmt.Errorf("failed to get user: %w", err)
	}

	err = row.Scan(
		&result.ID,
		&result.Email,
		&result.PasswordHash,
		&result.Name,
		&result.Role,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domainauth.GetDetailUserResult{}, databases.ErrNoRowFound
		}
		return domainauth.GetDetailUserResult{}, fmt.Errorf("token to scan row: %v", err)
	}

	return result, nil
}

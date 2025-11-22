package userrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases"

	domainuser "go-bootstrap/internal/domain/user"
)

func (r *repository) CreateUser(ctx context.Context, params domainuser.CreateUserParams) (domainuser.CreateUserResult, error) {
	query := `
		INSERT INTO users (email, password_hash, name, role, status, phone, gender, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now().UTC()
	result, err := r.db.RDBMS().ExecContext(ctx, query,
		params.Email,
		params.PasswordHash,
		params.Name,
		params.Role,
		domainuser.UserStatusActive,
		params.Phone,
		params.Gender,
		now,
		now,
	)

	if err != nil {
		return domainuser.CreateUserResult{}, fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domainuser.CreateUserResult{}, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return domainuser.CreateUserResult{
		ID:        fmt.Sprintf("%d", id),
		Email:     params.Email,
		Name:      params.Name,
		Role:      params.Role,
		Status:    domainuser.UserStatusActive,
		CreatedAt: now,
	}, nil
}

func (r *repository) GetDetailUser(ctx context.Context, filters domainuser.GetDetailUserFilters) (domainuser.GetDetailUserResult, error) {
	sq := r.db.Sq().Select(
		"id",
		"email",
		"password_hash",
		"name",
		"role",
		"status",
		"phone",
		"gender",
		"created_at",
		"updated_at",
	).From("users")

	if filters.UserID != nil {
		sq = sq.Where("id = ?", *filters.UserID)
	}

	if filters.Email != nil {
		sq = sq.Where("email = ?", *filters.Email)
	}

	sq = sq.Limit(1)

	var result domainuser.GetDetailUserResult
	row, err := r.db.RDBMS().QueryRowSq(ctx, sq, false)
	if err != nil {
		return domainuser.GetDetailUserResult{}, fmt.Errorf("failed to get user: %w", err)
	}

	err = row.Scan(
		&result.ID,
		&result.Email,
		&result.PasswordHash,
		&result.Name,
		&result.Role,
		&result.Status,
		&result.Phone,
		&result.Gender,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domainuser.GetDetailUserResult{}, databases.ErrNoRowFound
		}
		return domainuser.GetDetailUserResult{}, fmt.Errorf("token to scan row: %v", err)
	}

	return result, nil
}

func (r *repository) GetListUser(ctx context.Context, filters domainuser.GetListUserFilters) (domainuser.GetListUserResult, error) {
	countSq := r.db.Sq().Select("COUNT(*)").From("users")

	selectSq := r.db.Sq().Select(
		"id",
		"email",
		"password_hash",
		"name",
		"role",
		"status",
		"phone",
		"gender",
		"created_at",
		"updated_at",
	).From("users")

	if filters.Search != nil {
		searchPattern := "%" + *filters.Search + "%"
		searchCondition := sq.Or{
			sq.Like{"name": searchPattern},
			sq.Like{"email": searchPattern},
		}
		countSq = countSq.Where(searchCondition)
		selectSq = selectSq.Where(searchCondition)
	}

	if filters.Status != nil {
		countSq = countSq.Where("status = ?", *filters.Status)
		selectSq = selectSq.Where("status = ?", *filters.Status)
	}

	if filters.Role != nil {
		countSq = countSq.Where("role = ?", *filters.Role)
		selectSq = selectSq.Where("role = ?", *filters.Role)
	}

	selectSq = selectSq.OrderBy("created_at DESC")

	users := []domainuser.GetDetailUserResult{}
	pagination, err := r.db.RDBMS().QuerySqPagination(ctx, countSq, selectSq, false, filters.Pagination, func(rows *sql.Rows) error {
		for rows.Next() {
			var user domainuser.GetDetailUserResult
			err := rows.Scan(
				&user.ID,
				&user.Email,
				&user.PasswordHash,
				&user.Name,
				&user.Role,
				&user.Status,
				&user.Phone,
				&user.Gender,
				&user.CreatedAt,
				&user.UpdatedAt,
			)
			if err != nil {
				return fmt.Errorf("failed to scan user: %w", err)
			}
			users = append(users, user)
		}

		return nil
	})
	if err != nil {
		return domainuser.GetListUserResult{}, err
	}

	return domainuser.GetListUserResult{
		Users:      users,
		Pagination: pagination,
	}, nil
}

func (r *repository) UpdateUser(ctx context.Context, params domainuser.UpdateUserParams) (domainuser.UpdateUserResult, error) {
	updatedAt := time.Now().UTC()

	updateSq := r.db.Sq().Update("users").Set("updated_at", updatedAt)

	if params.Name != nil {
		updateSq = updateSq.Set("name", *params.Name)
	}

	if params.Phone != nil {
		updateSq = updateSq.Set("phone", *params.Phone)
	}

	if params.Gender != nil {
		updateSq = updateSq.Set("gender", *params.Gender)
	}

	updateSq = updateSq.Where("id = ?", params.UserID)

	_, err := r.db.RDBMS().ExecSq(ctx, updateSq, false)
	if err != nil {
		return domainuser.UpdateUserResult{}, fmt.Errorf("failed to update user: %w", err)
	}

	return domainuser.UpdateUserResult{
		UpdatedAt: updatedAt,
	}, nil
}

func (r *repository) UpdatePassword(ctx context.Context, params domainuser.UpdatePasswordParams) (domainuser.UpdatePasswordResult, error) {
	query := `
		UPDATE users
		SET password_hash = ?, updated_at = ?
		WHERE id = ?
	`

	updatedAt := time.Now().UTC()
	_, err := r.db.RDBMS().ExecContext(ctx, query,
		params.NewPasswordHash,
		updatedAt,
		params.UserID,
	)

	if err != nil {
		return domainuser.UpdatePasswordResult{}, fmt.Errorf("failed to update password: %w", err)
	}

	return domainuser.UpdatePasswordResult{
		UpdatedAt: updatedAt,
	}, nil
}

func (r *repository) UpdateStatus(ctx context.Context, params domainuser.UpdateStatusParams) (domainuser.UpdateStatusResult, error) {
	query := `
		UPDATE users
		SET status = ?, updated_at = ?
		WHERE id = ?
	`

	updatedAt := time.Now().UTC()
	_, err := r.db.RDBMS().ExecContext(ctx, query,
		params.Status,
		updatedAt,
		params.UserID,
	)

	if err != nil {
		return domainuser.UpdateStatusResult{}, fmt.Errorf("failed to update status: %w", err)
	}

	return domainuser.UpdateStatusResult{
		UpdatedAt: updatedAt,
	}, nil
}

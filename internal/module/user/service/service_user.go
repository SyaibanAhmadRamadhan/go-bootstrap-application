package userservice

import (
	"context"
	"errors"

	domainuser "go-bootstrap/internal/domain/user"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/apperror"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userRepo domainuser.UserRepositoryDatastore
}

func NewService(
	userRepo domainuser.UserRepositoryDatastore,
) *service {
	return &service{
		userRepo: userRepo,
	}
}

func (s *service) Register(ctx context.Context, input domainuser.RegisterInput) (domainuser.RegisterOutput, error) {
	existingUser, _ := s.userRepo.GetDetailUser(ctx, domainuser.GetDetailUserFilters{
		Email: &input.Email,
	})
	if existingUser.ID != "" {
		return domainuser.RegisterOutput{}, apperror.BadRequest("email already registered")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return domainuser.RegisterOutput{}, apperror.StdUnknown(err)
	}

	result, err := s.userRepo.CreateUser(ctx, domainuser.CreateUserParams{
		Email:        input.Email,
		PasswordHash: string(passwordHash),
		Name:         input.Name,
		Role:         domainuser.UserRoleUser,
		Phone:        input.Phone,
		Gender:       input.Gender,
	})
	if err != nil {
		return domainuser.RegisterOutput{}, apperror.StdUnknown(err)
	}

	return domainuser.RegisterOutput{
		UserID:    result.ID,
		Email:     result.Email,
		Name:      result.Name,
		CreatedAt: result.CreatedAt,
	}, nil
}

func (s *service) GetProfile(ctx context.Context, input domainuser.GetProfileInput) (domainuser.GetProfileOutput, error) {
	user, err := s.userRepo.GetDetailUser(ctx, domainuser.GetDetailUserFilters{
		UserID: &input.UserID,
	})
	if err != nil {
		if errors.Is(err, databases.ErrNoRowFound) {
			return domainuser.GetProfileOutput{}, apperror.BadRequest("user not found")
		}
		return domainuser.GetProfileOutput{}, apperror.StdUnknown(err)
	}

	return domainuser.GetProfileOutput{
		User: domainuser.User{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			Status:    user.Status,
			Gender:    user.Gender,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func (s *service) GetList(ctx context.Context, input domainuser.GetListInput) (domainuser.GetListOutput, error) {
	if input.Pagination.Page <= 0 {
		input.Pagination.Page = 1
	}
	if input.Pagination.PageSize <= 0 {
		input.Pagination.PageSize = 10
	}

	result, err := s.userRepo.GetListUser(ctx, domainuser.GetListUserFilters(input))
	if err != nil {
		return domainuser.GetListOutput{}, apperror.StdUnknown(err)
	}

	users := make([]domainuser.User, 0, len(result.Users))
	for _, u := range result.Users {
		users = append(users, domainuser.User{
			ID:        u.ID,
			Email:     u.Email,
			Name:      u.Name,
			Role:      u.Role,
			Status:    u.Status,
			Gender:    u.Gender,
			Phone:     u.Phone,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	return domainuser.GetListOutput{
		Users:      users,
		Pagination: result.Pagination,
	}, nil
}

func (s *service) UpdateProfile(ctx context.Context, input domainuser.UpdateProfileInput) (domainuser.UpdateProfileOutput, error) {
	user, err := s.userRepo.GetDetailUser(ctx, domainuser.GetDetailUserFilters{
		UserID: &input.UserID,
	})
	if err != nil {
		if errors.Is(err, databases.ErrNoRowFound) {
			return domainuser.UpdateProfileOutput{}, apperror.BadRequest("user not found")
		}
		return domainuser.UpdateProfileOutput{}, apperror.StdUnknown(err)
	}

	result, err := s.userRepo.UpdateUser(ctx, domainuser.UpdateUserParams(input))
	if err != nil {
		return domainuser.UpdateProfileOutput{}, apperror.StdUnknown(err)
	}

	updatedUser, err := s.userRepo.GetDetailUser(ctx, domainuser.GetDetailUserFilters{
		UserID: &input.UserID,
	})
	if err != nil {
		updatedUser = user
	}

	return domainuser.UpdateProfileOutput{
		User: domainuser.User{
			ID:        updatedUser.ID,
			Email:     updatedUser.Email,
			Name:      updatedUser.Name,
			Role:      updatedUser.Role,
			Status:    updatedUser.Status,
			Gender:    updatedUser.Gender,
			Phone:     updatedUser.Phone,
			CreatedAt: updatedUser.CreatedAt,
			UpdatedAt: updatedUser.UpdatedAt,
		},
		UpdatedAt: result.UpdatedAt,
	}, nil
}

func (s *service) ChangePassword(ctx context.Context, input domainuser.ChangePasswordInput) (domainuser.ChangePasswordOutput, error) {
	user, err := s.userRepo.GetDetailUser(ctx, domainuser.GetDetailUserFilters{
		UserID: &input.UserID,
	})
	if err != nil {
		if errors.Is(err, databases.ErrNoRowFound) {
			return domainuser.ChangePasswordOutput{}, apperror.BadRequest("user not found")
		}
		return domainuser.ChangePasswordOutput{}, apperror.StdUnknown(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.OldPassword))
	if err != nil {
		return domainuser.ChangePasswordOutput{}, apperror.BadRequest("invalid old password")
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return domainuser.ChangePasswordOutput{}, apperror.StdUnknown(err)
	}

	result, err := s.userRepo.UpdatePassword(ctx, domainuser.UpdatePasswordParams{
		UserID:          input.UserID,
		NewPasswordHash: string(newPasswordHash),
	})
	if err != nil {
		return domainuser.ChangePasswordOutput{}, apperror.StdUnknown(err)
	}

	return domainuser.ChangePasswordOutput{
		Success:   true,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

func (s *service) UpdateStatus(ctx context.Context, input domainuser.UpdateStatusInput) (domainuser.UpdateStatusOutput, error) {
	_, err := s.userRepo.GetDetailUser(ctx, domainuser.GetDetailUserFilters{
		UserID: &input.UserID,
	})
	if err != nil {
		if errors.Is(err, databases.ErrNoRowFound) {
			return domainuser.UpdateStatusOutput{}, apperror.BadRequest("user not found")
		}
		return domainuser.UpdateStatusOutput{}, apperror.StdUnknown(err)
	}

	result, err := s.userRepo.UpdateStatus(ctx, domainuser.UpdateStatusParams(input))
	if err != nil {
		return domainuser.UpdateStatusOutput{}, apperror.StdUnknown(err)
	}

	return domainuser.UpdateStatusOutput{
		Success:   true,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

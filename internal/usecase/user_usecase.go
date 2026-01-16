package usecase

import (
	"context"

	"go-clean-architecture/internal/domain"
	"go-clean-architecture/internal/dto"
)

type UserUseCase interface {
	Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.UserResponse, error)
	GetAll(ctx context.Context, page, limit int) ([]dto.UserResponse, int64, error)
	Update(ctx context.Context, id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id uint) error
}

type userUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(userRepo domain.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) Create(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, // In production, hash the password
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (u *userUseCase) GetByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (u *userUseCase) GetAll(ctx context.Context, page, limit int) ([]dto.UserResponse, int64, error) {
	users, total, err := u.userRepo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.UserResponse
	for _, user := range users {
		responses = append(responses, dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return responses, total, nil
}

func (u *userUseCase) Update(ctx context.Context, id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (u *userUseCase) Delete(ctx context.Context, id uint) error {
	return u.userRepo.Delete(ctx, id)
}

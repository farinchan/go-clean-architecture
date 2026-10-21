package usecase

import (
	"context"
	"errors"

	"github.com/your-username/go-clean-architecture/internal/dto"
	"github.com/your-username/go-clean-architecture/internal/entity"
	"github.com/your-username/go-clean-architecture/internal/repository"
	"github.com/your-username/go-clean-architecture/pkg/utils"
	"gorm.io/gorm"
)

// UserUseCase defines the user use case interface
type UserUseCase interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.UserResponse, error)
	GetAll(ctx context.Context, page, limit int) ([]dto.UserResponse, int64, error)
	Update(ctx context.Context, id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id uint) error
}

type userUseCase struct {
	userRepo   repository.UserRepository
	jwtManager *utils.JWTManager
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userRepo repository.UserRepository, jwtManager *utils.JWTManager) UserUseCase {
	return &userUseCase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Register registers a new user
func (u *userUseCase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error) {
	// Check if email already exists
	existingUser, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	var phonePtr *string
	if req.Phone != "" {
		phoneUser, err := u.userRepo.FindByPhone(ctx, req.Phone)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if phoneUser != nil {
			return nil, errors.New("phone number already registered")
		}
		phonePtr = &req.Phone
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    phonePtr,
		Password: hashedPassword,
		Role:     "user",
		IsActive: true,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return ToUserResponse(user), nil
}

// Login logs in a user
func (u *userUseCase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find user by email
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is not active")
	}

	// Generate JWT token
	token, err := u.jwtManager.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:             token,
		TokenType:         "Bearer",
		ExpiresIn:         int(u.jwtManager.GetExpiration().Seconds()),
		TwoFactorRequired: false,
		User:              ToUserResponse(user),
	}, nil
}

// GetByID gets a user by ID
func (u *userUseCase) GetByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return ToUserResponse(user), nil
}

// GetAll gets all users with pagination
func (u *userUseCase) GetAll(ctx context.Context, page, limit int) ([]dto.UserResponse, int64, error) {
	users, total, err := u.userRepo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var response []dto.UserResponse
	for _, user := range users {
		res := ToUserResponse(&user)
		if res != nil {
			response = append(response, *res)
		}
	}

	return response, total, nil
}

// Update updates a user
func (u *userUseCase) Update(ctx context.Context, id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" && req.Email != user.Email {
		existingUser, err := u.userRepo.FindByEmail(ctx, req.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already taken")
		}
		user.Email = req.Email
	}
	if req.Phone != "" && (user.Phone == nil || req.Phone != *user.Phone) {
		existingPhoneUser, err := u.userRepo.FindByPhone(ctx, req.Phone)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existingPhoneUser != nil && existingPhoneUser.ID != id {
			return nil, errors.New("phone already taken")
		}
		user.Phone = &req.Phone
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}
	if req.Is2FAEnabled != nil {
		user.Is2FAEnabled = *req.Is2FAEnabled
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return ToUserResponse(user), nil
}

// Delete deletes a user
func (u *userUseCase) Delete(ctx context.Context, id uint) error {
	_, err := u.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return u.userRepo.Delete(ctx, id)
}

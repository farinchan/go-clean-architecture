package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/your-username/go-clean-architecture/internal/dto"
	"github.com/your-username/go-clean-architecture/internal/repository"
	"github.com/your-username/go-clean-architecture/pkg/logger"
	"github.com/your-username/go-clean-architecture/pkg/mail"
	"github.com/your-username/go-clean-architecture/pkg/utils"
	"gorm.io/gorm"
)

// ProfileUseCase defines the profile use case interface
type ProfileUseCase interface {
	GetProfile(ctx context.Context, userID uint) (*dto.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error)
	Toggle2FA(ctx context.Context, userID uint, enable bool) (*dto.UserResponse, error)
	UpdatePhoto(ctx context.Context, userID uint, photoURL string) (*dto.UserResponse, error)
	ChangePassword(ctx context.Context, userID uint, req *dto.ChangePasswordRequest) error
	RegenerateNotificationSecret(ctx context.Context, userID uint) (*dto.UserResponse, error)
}

type profileUseCase struct {
	userRepo repository.UserRepository
	mailer   *mail.Mailer
}

// NewProfileUseCase creates a new profile use case
func NewProfileUseCase(userRepo repository.UserRepository, mailer *mail.Mailer) ProfileUseCase {
	return &profileUseCase{
		userRepo: userRepo,
		mailer:   mailer,
	}
}

func (u *profileUseCase) GetProfile(ctx context.Context, userID uint) (*dto.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return ToUserResponse(user), nil
}

func (u *profileUseCase) UpdateProfile(ctx context.Context, userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Email != "" && req.Email != user.Email {
		existingUser, err := u.userRepo.FindByEmail(ctx, req.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("email already taken")
		}
		user.Email = req.Email
	}

	if req.Phone != "" && (user.Phone == nil || req.Phone != *user.Phone) {
		existingPhoneUser, err := u.userRepo.FindByPhone(ctx, req.Phone)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existingPhoneUser != nil && existingPhoneUser.ID != userID {
			return nil, errors.New("phone already taken")
		}
		user.Phone = &req.Phone
	}

	if req.Is2FAEnabled != nil {
		user.Is2FAEnabled = *req.Is2FAEnabled
	}

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return ToUserResponse(user), nil
}

func (u *profileUseCase) Toggle2FA(ctx context.Context, userID uint, enable bool) (*dto.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.Is2FAEnabled = enable
	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return ToUserResponse(user), nil
}

func (u *profileUseCase) UpdatePhoto(ctx context.Context, userID uint, photoURL string) (*dto.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.Photo = &photoURL
	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return ToUserResponse(user), nil
}

func (u *profileUseCase) ChangePassword(ctx context.Context, userID uint, req *dto.ChangePasswordRequest) error {
	if req.NewPassword != req.ConfirmPassword {
		return errors.New("confirm_password must be equal to new_password")
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if !utils.CheckPassword(req.OldPassword, user.Password) {
		return errors.New("incorrect current password")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	if err := u.userRepo.Update(ctx, user); err != nil {
		return err
	}

	if u.mailer != nil {
		go func(email, name string) {
			if err := u.mailer.SendPasswordChangedAlert(email, name, time.Now()); err != nil {
				logger.Warnf("Failed to send password changed alert to %s: %v", email, err)
			}
		}(user.Email, user.Name)
	}

	return nil
}

func (u *profileUseCase) RegenerateNotificationSecret(ctx context.Context, userID uint) (*dto.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	newSecret := "whsec_" + strings.ReplaceAll(utils.GenerateUUID(), "-", "")
	user.NotificationSecret = &newSecret

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return ToUserResponse(user), nil
}

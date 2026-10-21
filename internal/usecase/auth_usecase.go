package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/your-username/go-clean-architecture/internal/dto"
	"github.com/your-username/go-clean-architecture/internal/entity"
	"github.com/your-username/go-clean-architecture/internal/repository"
	"github.com/your-username/go-clean-architecture/pkg/logger"
	"github.com/your-username/go-clean-architecture/pkg/mail"
	"github.com/your-username/go-clean-architecture/pkg/utils"
	"gorm.io/gorm"
)

// AuthUseCase defines the authentication use case interface
type AuthUseCase interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest, clientIP, userAgent string) (*dto.LoginResponse, error)
	Verify2FA(ctx context.Context, req *dto.Verify2FARequest, clientIP, userAgent string) (*dto.LoginResponse, error)
	ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) error
	VerifyResetToken(ctx context.Context, req *dto.VerifyResetTokenRequest) error
	ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.LoginResponse, error)
}

type authUseCase struct {
	userRepo          repository.UserRepository
	twoFactorRepo     repository.TwoFactorRepository
	passwordResetRepo repository.PasswordResetRepository
	jwtManager        *utils.JWTManager
	mailer            *mail.Mailer
}

// NewAuthUseCase creates a new auth use case
func NewAuthUseCase(
	userRepo repository.UserRepository,
	twoFactorRepo repository.TwoFactorRepository,
	passwordResetRepo repository.PasswordResetRepository,
	jwtManager *utils.JWTManager,
	mailer *mail.Mailer,
) AuthUseCase {
	return &authUseCase{
		userRepo:          userRepo,
		twoFactorRepo:     twoFactorRepo,
		passwordResetRepo: passwordResetRepo,
		jwtManager:        jwtManager,
		mailer:            mailer,
	}
}

// ToUserResponse maps entity.User to dto.UserResponse
func ToUserResponse(user *entity.User) *dto.UserResponse {
	if user == nil {
		return nil
	}
	return &dto.UserResponse{
		ID:                 user.ID,
		Name:               user.Name,
		Email:              user.Email,
		Phone:              user.Phone,
		Photo:              user.Photo,
		Role:               user.Role,
		IsActive:           user.IsActive,
		Is2FAEnabled:       user.Is2FAEnabled,
		NotificationSecret: user.NotificationSecret,
		CreatedAt:          user.CreatedAt,
		UpdatedAt:          user.UpdatedAt,
	}
}

// Register registers a new user
func (u *authUseCase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error) {
	// Check if email already exists
	existingUser, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Check if phone already exists if provided
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

	// Generate notification secret
	notificationSecret := "whsec_" + strings.ReplaceAll(utils.GenerateUUID(), "-", "")

	user := &entity.User{
		Name:               req.Name,
		Email:              req.Email,
		Phone:              phonePtr,
		Password:           hashedPassword,
		Role:               "user",
		IsActive:           true,
		Is2FAEnabled:       false,
		NotificationSecret: &notificationSecret,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return ToUserResponse(user), nil
}

// Login logs in a user, supporting 2FA and session alert email
func (u *authUseCase) Login(ctx context.Context, req *dto.LoginRequest, clientIP, userAgent string) (*dto.LoginResponse, error) {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, errors.New("account is not active")
	}

	// If 2FA is enabled, generate 6-digit OTP and pre_auth_token
	if user.Is2FAEnabled {
		otpCode, err := utils.GenerateOTP(6)
		if err != nil {
			return nil, errors.New("failed to generate verification code")
		}

		preAuthToken := utils.GenerateUUID()
		expiresAt := time.Now().Add(10 * time.Minute)

		// Clear previous 2FA codes for this user
		_ = u.twoFactorRepo.DeleteByUserID(ctx, user.ID)

		tfCode := &entity.TwoFactorCode{
			UserID:       user.ID,
			Code:         otpCode,
			PreAuthToken: preAuthToken,
			ExpiresAt:    expiresAt,
		}

		if err := u.twoFactorRepo.Create(ctx, tfCode); err != nil {
			return nil, err
		}

		// Asynchronously dispatch 2FA email
		if u.mailer != nil {
			go func(email, name, code string) {
				if err := u.mailer.Send2FACode(email, name, code); err != nil {
					logger.Warnf("Failed to send 2FA email to %s: %v", email, err)
				}
			}(user.Email, user.Name, otpCode)
		}

		return &dto.LoginResponse{
			TwoFactorRequired: true,
			PreAuthToken:      preAuthToken,
		}, nil
	}

	// 2FA is disabled, generate JWT Access Token directly
	token, err := u.jwtManager.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	expiresIn := int(u.jwtManager.GetExpiration().Seconds())
	if expiresIn <= 0 {
		expiresIn = 86400
	}

	// Asynchronously dispatch Login Session Alert email
	if u.mailer != nil {
		go func(email, name, ip, ua string) {
			if err := u.mailer.SendLoginAlert(email, name, ip, ua, time.Now()); err != nil {
				logger.Warnf("Failed to send login alert to %s: %v", email, err)
			}
		}(user.Email, user.Name, clientIP, userAgent)
	}

	return &dto.LoginResponse{
		Token:             token,
		TokenType:         "Bearer",
		ExpiresIn:         expiresIn,
		TwoFactorRequired: false,
		User:              ToUserResponse(user),
	}, nil
}

// Verify2FA verifies the 2FA OTP code and returns the JWT token
func (u *authUseCase) Verify2FA(ctx context.Context, req *dto.Verify2FARequest, clientIP, userAgent string) (*dto.LoginResponse, error) {
	tfCode, err := u.twoFactorRepo.FindByPreAuthToken(ctx, req.PreAuthToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Invalid or expired verification code")
		}
		return nil, err
	}

	if time.Now().After(tfCode.ExpiresAt) {
		_ = u.twoFactorRepo.DeleteByPreAuthToken(ctx, req.PreAuthToken)
		return nil, errors.New("Invalid or expired verification code")
	}

	if tfCode.Code != req.Code {
		return nil, errors.New("Invalid or expired verification code")
	}

	// Delete used code
	_ = u.twoFactorRepo.DeleteByPreAuthToken(ctx, req.PreAuthToken)

	user, err := u.userRepo.FindByID(ctx, tfCode.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	token, err := u.jwtManager.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	expiresIn := int(u.jwtManager.GetExpiration().Seconds())
	if expiresIn <= 0 {
		expiresIn = 86400
	}

	// Asynchronously dispatch Login Session Alert email
	if u.mailer != nil {
		go func(email, name, ip, ua string) {
			if err := u.mailer.SendLoginAlert(email, name, ip, ua, time.Now()); err != nil {
				logger.Warnf("Failed to send login alert to %s: %v", email, err)
			}
		}(user.Email, user.Name, clientIP, userAgent)
	}

	return &dto.LoginResponse{
		Token:             token,
		TokenType:         "Bearer",
		ExpiresIn:         expiresIn,
		TwoFactorRequired: false,
		User:              ToUserResponse(user),
	}, nil
}

// ForgotPassword generates a 6-digit OTP and sends email
func (u *authUseCase) ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) error {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return nil so client sees success response without revealing email existence
			return nil
		}
		return err
	}

	otpCode, err := utils.GenerateOTP(6)
	if err != nil {
		return errors.New("failed to generate verification code")
	}

	// Delete previous reset tokens for this email
	_ = u.passwordResetRepo.DeleteByEmail(ctx, req.Email)

	reset := &entity.PasswordReset{
		Email:     req.Email,
		Token:     otpCode,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := u.passwordResetRepo.Create(ctx, reset); err != nil {
		return err
	}

	if u.mailer != nil {
		go func(email, name, code string) {
			if err := u.mailer.SendPasswordResetCode(email, name, code); err != nil {
				logger.Warnf("Failed to send password reset code to %s: %v", email, err)
			}
		}(user.Email, user.Name, otpCode)
	}

	return nil
}

// VerifyResetToken verifies if the reset OTP code is valid and not expired
func (u *authUseCase) VerifyResetToken(ctx context.Context, req *dto.VerifyResetTokenRequest) error {
	reset, err := u.passwordResetRepo.FindByEmailAndToken(ctx, req.Email, req.Token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Invalid or expired verification code")
		}
		return err
	}

	if time.Now().After(reset.ExpiresAt) {
		_ = u.passwordResetRepo.DeleteByEmail(ctx, req.Email)
		return errors.New("Invalid or expired verification code")
	}

	return nil
}

// ResetPassword resets the user's password with the OTP code
func (u *authUseCase) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error {
	if req.Password != req.ConfirmPassword {
		return errors.New("confirm_password must be equal to password")
	}

	reset, err := u.passwordResetRepo.FindByEmailAndToken(ctx, req.Email, req.Token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Invalid or expired verification code")
		}
		return err
	}

	if time.Now().After(reset.ExpiresAt) {
		_ = u.passwordResetRepo.DeleteByEmail(ctx, req.Email)
		return errors.New("Invalid or expired verification code")
	}

	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	if err := u.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Delete used reset tokens
	_ = u.passwordResetRepo.DeleteByEmail(ctx, req.Email)

	return nil
}

// RefreshToken refreshes an existing valid JWT token
func (u *authUseCase) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.LoginResponse, error) {
	claims, err := u.jwtManager.ValidateToken(req.Token)
	if err != nil {
		return nil, errors.New("Invalid or expired token")
	}

	user, err := u.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("account is not active")
	}

	newToken, err := u.jwtManager.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	expiresIn := int(u.jwtManager.GetExpiration().Seconds())
	if expiresIn <= 0 {
		expiresIn = 86400
	}

	return &dto.LoginResponse{
		Token:             newToken,
		TokenType:         "Bearer",
		ExpiresIn:         expiresIn,
		TwoFactorRequired: false,
		User:              ToUserResponse(user),
	}, nil
}

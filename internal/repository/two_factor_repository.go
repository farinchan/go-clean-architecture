package repository

import (
	"context"

	"github.com/your-username/go-clean-architecture/internal/entity"
)

// TwoFactorRepository defines the 2FA repository interface
type TwoFactorRepository interface {
	Create(ctx context.Context, code *entity.TwoFactorCode) error
	FindByPreAuthToken(ctx context.Context, preAuthToken string) (*entity.TwoFactorCode, error)
	DeleteByPreAuthToken(ctx context.Context, preAuthToken string) error
	DeleteByUserID(ctx context.Context, userID uint) error
}

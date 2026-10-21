package repository

import (
	"context"

	"github.com/your-username/go-clean-architecture/internal/entity"
)

// PasswordResetRepository defines the password reset repository interface
type PasswordResetRepository interface {
	Create(ctx context.Context, reset *entity.PasswordReset) error
	FindByEmailAndToken(ctx context.Context, email, token string) (*entity.PasswordReset, error)
	DeleteByEmail(ctx context.Context, email string) error
}

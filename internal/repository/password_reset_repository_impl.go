package repository

import (
	"context"

	"github.com/your-username/go-clean-architecture/internal/entity"
	"gorm.io/gorm"
)

type passwordResetRepository struct {
	db *gorm.DB
}

// NewPasswordResetRepository creates a new password reset repository
func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) Create(ctx context.Context, reset *entity.PasswordReset) error {
	return r.db.WithContext(ctx).Create(reset).Error
}

func (r *passwordResetRepository) FindByEmailAndToken(ctx context.Context, email, token string) (*entity.PasswordReset, error) {
	var reset entity.PasswordReset
	if err := r.db.WithContext(ctx).Where("email = ? AND token = ?", email, token).First(&reset).Error; err != nil {
		return nil, err
	}
	return &reset, nil
}

func (r *passwordResetRepository) DeleteByEmail(ctx context.Context, email string) error {
	return r.db.WithContext(ctx).Where("email = ?", email).Delete(&entity.PasswordReset{}).Error
}

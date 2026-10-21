package repository

import (
	"context"

	"github.com/your-username/go-clean-architecture/internal/entity"
	"gorm.io/gorm"
)

type twoFactorRepository struct {
	db *gorm.DB
}

// NewTwoFactorRepository creates a new 2FA repository
func NewTwoFactorRepository(db *gorm.DB) TwoFactorRepository {
	return &twoFactorRepository{db: db}
}

func (r *twoFactorRepository) Create(ctx context.Context, code *entity.TwoFactorCode) error {
	return r.db.WithContext(ctx).Create(code).Error
}

func (r *twoFactorRepository) FindByPreAuthToken(ctx context.Context, preAuthToken string) (*entity.TwoFactorCode, error) {
	var code entity.TwoFactorCode
	if err := r.db.WithContext(ctx).Where("pre_auth_token = ?", preAuthToken).First(&code).Error; err != nil {
		return nil, err
	}
	return &code, nil
}

func (r *twoFactorRepository) DeleteByPreAuthToken(ctx context.Context, preAuthToken string) error {
	return r.db.WithContext(ctx).Where("pre_auth_token = ?", preAuthToken).Delete(&entity.TwoFactorCode{}).Error
}

func (r *twoFactorRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&entity.TwoFactorCode{}).Error
}

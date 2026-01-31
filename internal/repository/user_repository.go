package repository

import (
	"context"

	"github.com/your-username/go-clean-architecture/internal/entity"
)

// UserRepository defines the user repository interface
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindAll(ctx context.Context, page, limit int) ([]entity.User, int64, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error
}

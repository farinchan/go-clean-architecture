package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uint) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindAll(ctx context.Context, page, limit int) ([]User, int64, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
}

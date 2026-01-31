package entity

import (
	"time"

	"gorm.io/gorm"
)

// User represents the user entity
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:255;not null"`
	Email     string         `json:"email" gorm:"size:255;uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"size:255;not null"`
	Role      string         `json:"role" gorm:"size:50;default:'user'"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for the User model
func (User) TableName() string {
	return "users"
}

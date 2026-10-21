package entity

import (
	"time"

	"gorm.io/gorm"
)

// User represents the user entity
type User struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	Name               string         `json:"name" gorm:"size:255;not null"`
	Email              string         `json:"email" gorm:"size:255;uniqueIndex;not null"`
	Phone              *string        `json:"phone" gorm:"size:50;uniqueIndex"`
	Password           string         `json:"-" gorm:"size:255;not null"`
	Photo              *string        `json:"photo" gorm:"size:255"`
	Role               string         `json:"role" gorm:"size:50;default:'user'"`
	IsActive           bool           `json:"is_active" gorm:"default:true"`
	Is2FAEnabled       bool           `json:"is_2fa_enabled" gorm:"default:false"`
	NotificationSecret *string        `json:"notification_secret,omitempty" gorm:"size:100"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for the User model
func (User) TableName() string {
	return "users"
}

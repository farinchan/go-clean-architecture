package entity

import "time"

// PasswordReset represents a password reset OTP record
type PasswordReset struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"size:255;not null;index"`
	Token     string    `json:"token" gorm:"size:10;not null;index"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName returns the table name for PasswordReset
func (PasswordReset) TableName() string {
	return "password_resets"
}

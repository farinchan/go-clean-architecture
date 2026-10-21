package entity

import "time"

// TwoFactorCode represents a two-factor verification code
type TwoFactorCode struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null;index"`
	Code         string    `json:"code" gorm:"size:10;not null"`
	PreAuthToken string    `json:"pre_auth_token" gorm:"size:100;uniqueIndex;not null"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName returns the table name for TwoFactorCode
func (TwoFactorCode) TableName() string {
	return "two_factor_codes"
}

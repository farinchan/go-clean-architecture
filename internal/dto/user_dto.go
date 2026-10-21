package dto

import "time"

// RegisterRequest represents the register request body
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100" example:"Fajri Ramadhan"`
	Email    string `json:"email" binding:"required,email" example:"fajri@example.com"`
	Phone    string `json:"phone" binding:"omitempty,min=8,max=20" example:"+6281234567890"`
	Password string `json:"password" binding:"required,min=6" example:"SecurePassword123"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"fajri@example.com"`
	Password string `json:"password" binding:"required" example:"SecurePassword123"`
}

// UpdateUserRequest represents the update user request body
type UpdateUserRequest struct {
	Name         string `json:"name" binding:"omitempty,min=2,max=100" example:"John Doe Updated"`
	Email        string `json:"email" binding:"omitempty,email" example:"john.updated@example.com"`
	Phone        string `json:"phone" binding:"omitempty,min=8,max=20" example:"+6281298765432"`
	Password     string `json:"password" binding:"omitempty,min=6" example:"newpassword123"`
	Is2FAEnabled *bool  `json:"is_2fa_enabled"`
}

// UserResponse represents the user response
type UserResponse struct {
	ID                 uint      `json:"id" example:"1"`
	Name               string    `json:"name" example:"Fajri Ramadhan"`
	Email              string    `json:"email" example:"fajri@example.com"`
	Phone              *string   `json:"phone" example:"+6281234567890"`
	Photo              *string   `json:"photo" example:"/uploads/profile/550e8400-e29b-41d4-a716-446655440000.jpg"`
	Role               string    `json:"role" example:"user"`
	IsActive           bool      `json:"is_active" example:"true"`
	Is2FAEnabled       bool      `json:"is_2fa_enabled" example:"false"`
	NotificationSecret *string   `json:"notification_secret,omitempty" example:"whsec_550e8400e29b41d4a716446655440000"`
	CreatedAt          time.Time `json:"created_at" example:"2026-10-16T10:00:00Z"`
	UpdatedAt          time.Time `json:"updated_at" example:"2026-10-16T10:00:00Z"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token             string        `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType         string        `json:"token_type,omitempty" example:"Bearer"`
	ExpiresIn         int           `json:"expires_in,omitempty" example:"86400"`
	TwoFactorRequired bool          `json:"two_factor_required" example:"false"`
	PreAuthToken      string        `json:"pre_auth_token,omitempty" example:"a1b2c3d4-e5f6-7890-abcd-ef1234567890"`
	User              *UserResponse `json:"user,omitempty"`
}

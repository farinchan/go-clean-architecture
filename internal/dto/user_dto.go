package dto

import "time"

// RegisterRequest represents the register request body
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100" example:"John Doe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// UpdateUserRequest represents the update user request body
type UpdateUserRequest struct {
	Name     string `json:"name" binding:"omitempty,min=2,max=100" example:"John Doe Updated"`
	Email    string `json:"email" binding:"omitempty,email" example:"john.updated@example.com"`
	Password string `json:"password" binding:"omitempty,min=6" example:"newpassword123"`
}

// UserResponse represents the user response
type UserResponse struct {
	ID        uint      `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"john@example.com"`
	Role      string    `json:"role" example:"user"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserResponse `json:"user"`
}

package dto

// UpdateProfileRequest represents the update profile request body
type UpdateProfileRequest struct {
	Name         string `json:"name" binding:"omitempty,min=2,max=100" example:"Fajri Ramadhan Updated"`
	Email        string `json:"email" binding:"omitempty,email" example:"fajri.new@example.com"`
	Phone        string `json:"phone" binding:"omitempty,min=8,max=20" example:"+6281298765432"`
	Is2FAEnabled *bool  `json:"is_2fa_enabled"`
}

// Toggle2FARequest represents the toggle 2FA request body
type Toggle2FARequest struct {
	Enable *bool `json:"enable" binding:"required"`
}

// ChangePasswordRequest represents the change password request body
type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required" example:"SecurePassword123"`
	NewPassword     string `json:"new_password" binding:"required,min=6" example:"NewSecretPassword2026!"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword" example:"NewSecretPassword2026!"`
}

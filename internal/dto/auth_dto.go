package dto

// Verify2FARequest represents the 2FA verification request body
type Verify2FARequest struct {
	PreAuthToken string `json:"pre_auth_token" binding:"required" example:"a1b2c3d4-e5f6-7890-abcd-ef1234567890"`
	Code         string `json:"code" binding:"required,min=6,max=6" example:"849201"`
}

// ForgotPasswordRequest represents the forgot password request body
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"fajri@example.com"`
}

// VerifyResetTokenRequest represents the verify reset token request body
type VerifyResetTokenRequest struct {
	Email string `json:"email" binding:"required,email" example:"fajri@example.com"`
	Token string `json:"token" binding:"required" example:"492018"`
}

// ResetPasswordRequest represents the reset password request body
type ResetPasswordRequest struct {
	Email           string `json:"email" binding:"required,email" example:"fajri@example.com"`
	Token           string `json:"token" binding:"required" example:"492018"`
	Password        string `json:"password" binding:"required,min=6" example:"MyNewSecretPassword2026"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password" example:"MyNewSecretPassword2026"`
}

// RefreshTokenRequest represents the refresh token request body
type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

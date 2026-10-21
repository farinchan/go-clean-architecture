package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/your-username/go-clean-architecture/internal/dto"
	"github.com/your-username/go-clean-architecture/internal/usecase"
	"github.com/your-username/go-clean-architecture/pkg/response"
	"github.com/your-username/go-clean-architecture/pkg/validator"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authUseCase usecase.AuthUseCase
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authUseCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, name, phone, and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register request"
// @Success 201 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := validator.FormatValidationErrors(err)
		response.ValidationError(c, errors)
		return
	}

	user, err := h.authUseCase.Register(c.Request.Context(), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "User registered successfully", user)
}

// Login godoc
// @Summary Login user
// @Description Login with email and password, returns JWT token or requires 2FA OTP
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} response.Response{data=dto.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := validator.FormatValidationErrors(err)
		response.ValidationError(c, errors)
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	result, err := h.authUseCase.Login(c.Request.Context(), &req, clientIP, userAgent)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, "Login successful", result)
}

// Verify2FA godoc
// @Summary Verify 2FA OTP
// @Description Verify 2FA OTP code and complete login to get JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.Verify2FARequest true "Verify 2FA request"
// @Success 200 {object} response.Response{data=dto.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /api/v1/auth/verify-2fa [post]
func (h *AuthHandler) Verify2FA(c *gin.Context) {
	var req dto.Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := validator.FormatValidationErrors(err)
		response.ValidationError(c, errors)
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	result, err := h.authUseCase.Verify2FA(c.Request.Context(), &req, clientIP, userAgent)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "2FA verification successful", result)
}

// ForgotPassword godoc
// @Summary Forgot password request
// @Description Send a 6-digit OTP code to the registered email for password reset
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /api/v1/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := validator.FormatValidationErrors(err)
		response.ValidationError(c, errors)
		return
	}

	if err := h.authUseCase.ForgotPassword(c.Request.Context(), &req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Password reset verification code has been sent to your email", nil)
}

// VerifyResetToken godoc
// @Summary Verify reset password token
// @Description Verify the validity of a password reset OTP code
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.VerifyResetTokenRequest true "Verify reset token request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /api/v1/auth/verify-reset-token [post]
func (h *AuthHandler) VerifyResetToken(c *gin.Context) {
	var req dto.VerifyResetTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := validator.FormatValidationErrors(err)
		response.ValidationError(c, errors)
		return
	}

	if err := h.authUseCase.VerifyResetToken(c.Request.Context(), &req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Verification code is valid", nil)
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset user password using the OTP code
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /api/v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := validator.FormatValidationErrors(err)
		response.ValidationError(c, errors)
		return
	}

	if err := h.authUseCase.ResetPassword(c.Request.Context(), &req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Password has been successfully reset. You can now login.", nil)
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Description Generate a new fresh JWT Access Token from an existing valid token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} response.Response{data=dto.LoginResponse}
// @Failure 401 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /api/v1/auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := validator.FormatValidationErrors(err)
		response.ValidationError(c, errors)
		return
	}

	result, err := h.authUseCase.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, "Token refreshed successfully", result)
}

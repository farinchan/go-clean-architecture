package handler

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/your-username/go-clean-architecture/internal/dto"
	"github.com/your-username/go-clean-architecture/internal/usecase"
	"github.com/your-username/go-clean-architecture/pkg/response"
	"github.com/your-username/go-clean-architecture/pkg/utils"
	"github.com/your-username/go-clean-architecture/pkg/validator"
)

// ProfileHandler handles HTTP requests for user profile
type ProfileHandler struct {
	profileUseCase usecase.ProfileUseCase
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(profileUseCase usecase.ProfileUseCase) *ProfileHandler {
	return &ProfileHandler{profileUseCase: profileUseCase}
}

// GetProfile godoc
// @Summary Get current user profile
// @Description Retrieve the profile of the currently authenticated user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/profile [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	profile, err := h.profileUseCase.GetProfile(c.Request.Context(), userID.(uint))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Profile retrieved successfully", profile)
}

// UpdateProfile godoc
// @Summary Update profile
// @Description Update user profile information (name, email, phone, 2FA status)
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /api/v1/profile [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := validator.FormatValidationErrors(err)
		response.ValidationError(c, errors)
		return
	}

	profile, err := h.profileUseCase.UpdateProfile(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Profile updated successfully", profile)
}

// Toggle2FA godoc
// @Summary Toggle 2FA setting
// @Description Enable or disable Two-Factor Authentication via email
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.Toggle2FARequest true "Toggle 2FA request"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /api/v1/profile/2fa [patch]
func (h *ProfileHandler) Toggle2FA(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req dto.Toggle2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := validator.FormatValidationErrors(err)
		response.ValidationError(c, errors)
		return
	}

	profile, err := h.profileUseCase.Toggle2FA(c.Request.Context(), userID.(uint), *req.Enable)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	msg := "Two-factor authentication disabled successfully"
	if *req.Enable {
		msg = "Two-factor authentication enabled successfully"
	}

	response.Success(c, msg, profile)
}

// UpdatePhoto godoc
// @Summary Update profile photo
// @Description Upload and update avatar photo (PNG, JPG, JPEG, WEBP, max 5MB)
// @Tags Profile
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param photo formData file true "Avatar image file"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/profile/photo [patch]
func (h *ProfileHandler) UpdatePhoto(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		response.BadRequest(c, "photo file is required", nil)
		return
	}
	defer file.Close()

	// Max 5MB
	const maxFileSize = 5 * 1024 * 1024
	if header.Size > maxFileSize {
		response.BadRequest(c, "file size exceeds 5MB limit", nil)
		return
	}

	// Allowed extensions
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		response.BadRequest(c, "invalid file type. Allowed: jpg, jpeg, png, webp", nil)
		return
	}

	// Ensure upload directory exists
	uploadDir := filepath.Join(".", "uploads", "profile")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.InternalServerError(c, "failed to create upload directory")
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s%s", utils.GenerateUUID(), ext)
	dstPath := filepath.Join(uploadDir, filename)

	out, err := os.Create(dstPath)
	if err != nil {
		response.InternalServerError(c, "failed to save file")
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		response.InternalServerError(c, "failed to write file")
		return
	}

	photoURL := fmt.Sprintf("/uploads/profile/%s", filename)
	profile, err := h.profileUseCase.UpdatePhoto(c.Request.Context(), userID.(uint), photoURL)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Profile photo updated successfully", profile)
}

// ChangePassword godoc
// @Summary Change password
// @Description Change password for the currently authenticated user
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChangePasswordRequest true "Change password request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /api/v1/profile/change-password [patch]
func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := validator.FormatValidationErrors(err)
		response.ValidationError(c, errors)
		return
	}

	if err := h.profileUseCase.ChangePassword(c.Request.Context(), userID.(uint), &req); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Password changed successfully", nil)
}

// RegenerateNotificationSecret godoc
// @Summary Regenerate notification secret
// @Description Regenerate notification webhook secret (whsec_...)
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/profile/notification-secret/regenerate [post]
func (h *ProfileHandler) RegenerateNotificationSecret(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	profile, err := h.profileUseCase.RegenerateNotificationSecret(c.Request.Context(), userID.(uint))
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Notification webhook secret regenerated successfully", profile)
}

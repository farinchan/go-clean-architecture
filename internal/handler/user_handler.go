package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-username/go-clean-architecture/internal/dto"
	"github.com/your-username/go-clean-architecture/internal/usecase"
	"github.com/your-username/go-clean-architecture/pkg/response"
	"github.com/your-username/go-clean-architecture/pkg/validator"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userUseCase usecase.UserUseCase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register request"
// @Success 201 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /api/v1/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := validator.FormatValidationErrors(err)
		response.ValidationError(c, errors)
		return
	}

	user, err := h.userUseCase.Register(c.Request.Context(), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "User registered successfully", user)
}

// Login godoc
// @Summary Login user
// @Description Login with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} response.Response{data=dto.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := validator.FormatValidationErrors(err)
		response.ValidationError(c, errors)
		return
	}

	result, err := h.userUseCase.Login(c.Request.Context(), &req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, "Login successful", result)
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get a specific user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 404 {object} response.Response
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	user, err := h.userUseCase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "User retrieved successfully", user)
}

// GetUsers godoc
// @Summary Get all users
// @Description Get all users with pagination
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]dto.UserResponse}
// @Failure 500 {object} response.Response
// @Router /api/v1/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := h.userUseCase.GetAll(c.Request.Context(), page, limit)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	meta := response.BuildMeta(page, limit, total)
	response.SuccessWithMeta(c, "Users retrieved successfully", users, meta)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update a specific user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body dto.UpdateUserRequest true "Update user request"
// @Security BearerAuth
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors := validator.FormatValidationErrors(err)
		response.ValidationError(c, errors)
		return
	}

	user, err := h.userUseCase.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "User updated successfully", user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a specific user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	if err := h.userUseCase.Delete(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "User deleted successfully", nil)
}

// GetCurrentUser godoc
// @Summary Get current user
// @Description Get the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 401 {object} response.Response
// @Router /api/v1/users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	user, err := h.userUseCase.GetByID(c.Request.Context(), userID.(uint))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "User retrieved successfully", user)
}

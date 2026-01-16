package handler

import (
	"net/http"
	"strconv"

	"kopeta-backend/internal/dto"
	"kopeta-backend/internal/usecase"
	"kopeta-backend/pkg/response"
	"kopeta-backend/pkg/validator"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
	validator   *validator.CustomValidator
}

func NewUserHandler(userUseCase usecase.UserUseCase, validator *validator.CustomValidator) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		validator:   validator,
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		errors := validator.FormatValidationErrors(err)
		response.BadRequest(c, "Validation failed", errors)
		return
	}

	user, err := h.userUseCase.Create(c.Request.Context(), &req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Created(c, "User created successfully", user)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	user, err := h.userUseCase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	response.Success(c, "User retrieved successfully", user)
}

func (h *UserHandler) GetAll(c *gin.Context) {
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

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	pagination := response.Pagination{
		CurrentPage: page,
		PerPage:     limit,
		TotalPages:  totalPages,
		TotalItems:  total,
	}

	response.SuccessWithPagination(c, "Users retrieved successfully", users, pagination)
}

func (h *UserHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		errors := validator.FormatValidationErrors(err)
		response.BadRequest(c, "Validation failed", errors)
		return
	}

	user, err := h.userUseCase.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	response.Success(c, "User updated successfully", user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	if err := h.userUseCase.Delete(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "User not found")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

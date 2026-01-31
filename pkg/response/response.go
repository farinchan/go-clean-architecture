package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents the standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta holds pagination metadata
type Meta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
}

// Success sends a success response
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMeta sends a success response with pagination meta
func SuccessWithMeta(c *gin.Context, message string, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// Created sends a created response
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// NoContent sends a no content response
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// BadRequest sends a bad request error response
func BadRequest(c *gin.Context, message string, err interface{}) {
	Error(c, http.StatusBadRequest, message, err)
}

// Unauthorized sends an unauthorized error response
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, nil)
}

// Forbidden sends a forbidden error response
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message, nil)
}

// NotFound sends a not found error response
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, nil)
}

// Conflict sends a conflict error response
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, message, nil)
}

// UnprocessableEntity sends an unprocessable entity error response
func UnprocessableEntity(c *gin.Context, message string, err interface{}) {
	Error(c, http.StatusUnprocessableEntity, message, err)
}

// InternalServerError sends an internal server error response
func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message, nil)
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Success: false,
		Message: "Validation failed",
		Error:   errors,
	})
}

// BuildMeta creates pagination metadata
func BuildMeta(page, perPage int, total int64) *Meta {
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &Meta{
		CurrentPage: page,
		PerPage:     perPage,
		Total:       total,
		TotalPages:  totalPages,
	}
}

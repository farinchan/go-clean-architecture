package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError represents an application error
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap implements the errors.Unwrap interface
func (e *AppError) Unwrap() error {
	return e.Err
}

// Common errors
var (
	ErrNotFound          = &AppError{Code: http.StatusNotFound, Message: "Resource not found"}
	ErrBadRequest        = &AppError{Code: http.StatusBadRequest, Message: "Bad request"}
	ErrUnauthorized      = &AppError{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	ErrForbidden         = &AppError{Code: http.StatusForbidden, Message: "Forbidden"}
	ErrConflict          = &AppError{Code: http.StatusConflict, Message: "Resource conflict"}
	ErrInternalServer    = &AppError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	ErrValidation        = &AppError{Code: http.StatusUnprocessableEntity, Message: "Validation error"}
	ErrInvalidCredential = &AppError{Code: http.StatusUnauthorized, Message: "Invalid email or password"}
	ErrUserNotActive     = &AppError{Code: http.StatusForbidden, Message: "User account is not active"}
	ErrEmailTaken        = &AppError{Code: http.StatusConflict, Message: "Email is already registered"}
)

// NewAppError creates a new AppError
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// WrapError wraps an error with an AppError
func WrapError(appErr *AppError, err error) *AppError {
	return &AppError{
		Code:    appErr.Code,
		Message: appErr.Message,
		Err:     err,
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError gets the AppError from an error
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return ErrInternalServer
}

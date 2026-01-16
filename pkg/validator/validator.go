package validator

import (
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func NewValidator() *CustomValidator {
	return &CustomValidator{
		Validator: validator.New(),
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func (cv *CustomValidator) ValidateVar(field interface{}, tag string) error {
	return cv.Validator.Var(field, tag)
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, ValidationError{
				Field:   e.Field(),
				Message: getErrorMessage(e),
			})
		}
	}

	return errors
}

func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	case "gte":
		return "Value must be greater than or equal to " + e.Param()
	case "lte":
		return "Value must be less than or equal to " + e.Param()
	default:
		return "Invalid value"
	}
}

package validator

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// CustomValidator holds custom validators
type CustomValidator struct {
	validate *validator.Validate
}

// NewValidator creates a new custom validator
func NewValidator() *CustomValidator {
	v := validator.New()

	// Use JSON tag names in validation errors
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validators here
	// v.RegisterValidation("custom_tag", customValidationFunc)

	return &CustomValidator{validate: v}
}

// Validate validates a struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validate.Struct(i)
}

// RegisterGinValidator registers the custom validator with Gin
func RegisterGinValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Use JSON tag names
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		// Register custom validators
		// v.RegisterValidation("custom_tag", customValidationFunc)
	}
}

// FormatValidationErrors formats validation errors to a map
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors[e.Field()] = getErrorMessage(e)
		}
	}

	return errors
}

// getErrorMessage returns a human-readable error message
func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	case "gte":
		return "Value must be greater than or equal to " + fe.Param()
	case "lte":
		return "Value must be less than or equal to " + fe.Param()
	case "oneof":
		return "Value must be one of: " + fe.Param()
	case "url":
		return "Invalid URL format"
	case "uuid":
		return "Invalid UUID format"
	case "numeric":
		return "Value must be numeric"
	case "alpha":
		return "Value must contain only letters"
	case "alphanum":
		return "Value must contain only letters and numbers"
	default:
		return "Invalid value for " + fe.Field()
	}
}

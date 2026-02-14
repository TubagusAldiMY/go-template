package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() error {
	validate = validator.New()

	// Register custom validators
	if err := validate.RegisterValidation("password", validatePassword); err != nil {
		return fmt.Errorf("failed to register password validator: %w", err)
	}

	if err := validate.RegisterValidation("username", validateUsername); err != nil {
		return fmt.Errorf("failed to register username validator: %w", err)
	}

	return nil
}

func Validate(i interface{}) error {
	return validate.Struct(i)
}

func ValidateVar(field interface{}, tag string) error {
	return validate.Var(field, tag)
}

func GetValidator() *validator.Validate {
	return validate
}

// Custom validators

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Minimum 8 characters
	if len(password) < 8 {
		return false
	}

	// At least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// At least one lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// At least one digit
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	// At least one special character
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	return hasUpper && hasLower && hasDigit && hasSpecial
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	// Length between 3 and 30
	if len(username) < 3 || len(username) > 30 {
		return false
	}

	// Only alphanumeric, underscore, and hyphen
	matched := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(username)
	return matched
}

// FormatValidationErrors formats validation errors into readable messages
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())

			switch e.Tag() {
			case "required":
				errors[field] = fmt.Sprintf("%s is required", field)
			case "email":
				errors[field] = "invalid email format"
			case "min":
				errors[field] = fmt.Sprintf("%s must be at least %s characters", field, e.Param())
			case "max":
				errors[field] = fmt.Sprintf("%s must not exceed %s characters", field, e.Param())
			case "password":
				errors[field] = "password must be at least 8 characters and contain uppercase, lowercase, digit, and special character"
			case "username":
				errors[field] = "username must be 3-30 characters and contain only alphanumeric, underscore, or hyphen"
			case "uuid":
				errors[field] = "invalid UUID format"
			default:
				errors[field] = fmt.Sprintf("%s is invalid", field)
			}
		}
	}

	return errors
}

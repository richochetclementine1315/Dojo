package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator is the global validator instance
var Validator *validator.Validate

// InitValidator initializes the validator
func InitValidator() {
	Validator = validator.New()

	// Register custom validators here if needed
	// Example: Validator.RegisterValidation("custom_tag", customValidatorFunc)
}

// ValidateStruct validates a struct using the validator
func ValidateStruct(s interface{}) error {
	if Validator == nil {
		InitValidator()
	}

	err := Validator.Struct(s)
	if err != nil {
		return FormatValidationError(err)
	}
	return nil
}

// FormatValidationError formats validation errors into readable messages
func FormatValidationError(err error) error {
	if err == nil {
		return nil
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	var messages []string
	for _, e := range validationErrs {
		message := formatFieldError(e)
		messages = append(messages, message)
	}

	return fmt.Errorf("%s", strings.Join(messages, "; "))
}

// formatFieldError formats a single field validation error
func formatFieldError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, e.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, e.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

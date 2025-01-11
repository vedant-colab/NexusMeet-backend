package utils

import (
	"github.com/go-playground/validator/v10"
)

// ErrorResponse represents a validation error with details.
type ErrorResponse struct {
	Error       bool        `json:"error"`       // Indicates if there's an error
	FailedField string      `json:"failedField"` // The struct field that failed validation
	Tag         string      `json:"tag"`         // The validation tag that failed
	Value       interface{} `json:"value"`       // The value that caused the validation failure
}

// XValidator wraps the validator.Validate instance.
type XValidator struct {
	validator *validator.Validate
}

// validate is a singleton instance of validator.Validate.
var validate = validator.New()

// NewValidator creates and returns a new XValidator instance.
func NewValidator() *XValidator {
	return &XValidator{
		validator: validate,
	}
}

// Validate validates the given struct and returns a slice of ErrorResponse.
func (v *XValidator) Validate(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	// Perform validation
	errs := v.validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ErrorResponse{
				Error:       true,
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Value:       err.Value(),
			})
		}
	}

	return validationErrors
}

// FormatValidationErrors formats validation errors into a user-friendly map.
func FormatValidationErrors(errors []ErrorResponse) map[string]string {
	errorMap := make(map[string]string)
	for _, err := range errors {
		errorMap[err.FailedField] = err.Tag
	}
	return errorMap
}

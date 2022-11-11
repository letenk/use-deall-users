package web

import "github.com/go-playground/validator/v10"

// FormatValidationError for iteration error from package validator
// Because when error, will the return many error
func FormatValidationError(err error) []string {
	// Create variable err with data type slice string
	var errors []string

	// Process iteration errors
	for _, e := range err.(validator.ValidationErrors) {
		// Append every error message to var errors
		errors = append(errors, e.Error())
	}

	return errors
}

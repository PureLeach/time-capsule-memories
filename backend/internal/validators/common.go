package validators

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Common validator instance
var commonValidate = validator.New()

// ValidateStruct validates the provided struct and returns a description of validation errors.
func ValidateStruct(data interface{}) error {
	if err := commonValidate.Struct(data); err != nil {
		// Cast the error to a ValidationErrors type to extract field-level errors
		validationErrors := err.(validator.ValidationErrors)
		var errorMessages []string
		// Iterate through validation errors and format the error messages
		for _, ve := range validationErrors {
			errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' failed validation: rule '%s'", ve.Field(), ve.Tag()))
		}
		// Return the concatenated error messages
		return fmt.Errorf("%s", strings.Join(errorMessages, "; "))
	}
	return nil
}

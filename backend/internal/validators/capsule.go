package validators

import (
	"errors"
	"time"

	"time_capsule_memories/internal/models"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate *validator.Validate

// Initialize the validator and register custom validation rules.
func init() {
	validate = validator.New()

	// Register custom validation rules
	validate.RegisterValidation("send_at_date_format", sendAtDateFormat)
	validate.RegisterValidation("future_date", futureDate)
}

// Custom validator: checks if the date format is correct (YYYY-MM-DD).
func sendAtDateFormat(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().String())
	return err == nil
}

// Custom validator: checks if the date is in the future.
func futureDate(fl validator.FieldLevel) bool {
	parsedDate, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}
	today := time.Now().Truncate(24 * time.Hour)
	return parsedDate.After(today)
}

// ValidateCapsule performs validation of the CreateCapsuleRequest model.
func ValidateCapsule(capsule models.CreateCapsuleRequest) error {
	err := validate.Struct(capsule)
	if err == nil {
		return nil
	}

	validationErrors := err.(validator.ValidationErrors)
	for _, ve := range validationErrors {
		// Switch based on validation tag to return appropriate error messages.
		switch ve.Tag() {
		case "send_at_date_format":
			return errors.New("invalid field `send_at` - The date must be in YYYY-MM-DD format")
		case "future_date":
			return errors.New("invalid field `send_at` - The date must be in the future")
		case "max":
			if ve.Field() == "Message" {
				return errors.New("invalid field `message` - The message cannot exceed 4096 characters")
			}
		default:
			return errors.New("invalid value in field " + ve.Field())
		}
	}

	return nil
}

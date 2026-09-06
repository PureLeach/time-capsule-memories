// Package validators adds the project's rules to go-playground/validator and
// renders its field errors as client-safe messages.
package validators

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

const dateLayout = "2006-01-02"

var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Report wire names, so a client can map messages onto the JSON it sent.
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if name == "" || name == "-" {
			return field.Name
		}
		return name
	})

	for tag, fn := range map[string]validator.Func{
		"send_at_date_format": sendAtDateFormat,
		"future_date":         futureDate,
		"image_content_type":  imageContentType,
	} {
		if err := v.RegisterValidation(tag, fn); err != nil {
			panic(fmt.Sprintf("register validation %q: %v", tag, err))
		}
	}

	return v
}

func sendAtDateFormat(fl validator.FieldLevel) bool {
	_, err := time.Parse(dateLayout, fl.Field().String())
	return err == nil
}

// futureDate works in UTC, matching the dispatcher, so a date the UI offers is
// never rejected as past.
func futureDate(fl validator.FieldLevel) bool {
	parsed, err := time.Parse(dateLayout, fl.Field().String())
	if err != nil {
		return false
	}
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return parsed.After(today)
}

// An allowlist rather than an "image/*" prefix: a client must not be able to get
// a signature for image/svg+xml, which browsers treat as active content.
var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
}

func imageContentType(fl validator.FieldLevel) bool {
	return allowedImageTypes[strings.ToLower(strings.TrimSpace(fl.Field().String()))]
}

// ValidateStruct reports every failing field, not just the first.
func ValidateStruct(data any) error {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	var fieldErrors validator.ValidationErrors
	if !errors.As(err, &fieldErrors) {
		return fmt.Errorf("validate: %w", err)
	}

	messages := make([]string, 0, len(fieldErrors))
	for _, fe := range fieldErrors {
		messages = append(messages, describe(fe))
	}
	return errors.New(strings.Join(messages, "; "))
}

func describe(fe validator.FieldError) string {
	field := fe.Field()

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("field `%s` is required", field)
	case "send_at_date_format":
		return fmt.Sprintf("field `%s` must be a date in YYYY-MM-DD format", field)
	case "future_date":
		return fmt.Sprintf("field `%s` must be a date in the future", field)
	case "email":
		return fmt.Sprintf("field `%s` must be a valid email address", field)
	case "uuid4":
		return fmt.Sprintf("field `%s` must be a valid UUID", field)
	case "image_content_type":
		return fmt.Sprintf("field `%s` must be one of: image/jpeg, image/png, image/webp, image/gif", field)
	case "max":
		return fmt.Sprintf("field `%s` must be at most %s characters long", field, fe.Param())
	default:
		return fmt.Sprintf("field `%s` failed the %q rule", field, fe.Tag())
	}
}

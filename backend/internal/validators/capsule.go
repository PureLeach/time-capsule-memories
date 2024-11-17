package validators

import (
	"errors"
	"time"

	"time_capsule_memories/internal/models"

	"github.com/go-playground/validator/v10"
)

// Валидатор
var validate *validator.Validate

// Инициализация валидатора
func init() {
	validate = validator.New()

	// Регистрация кастомных правил
	validate.RegisterValidation("send_at_date_format", sendAtDateFormat)
	validate.RegisterValidation("future_date", futureDate)
}

// Кастомный валидатор: формат даты
func sendAtDateFormat(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().String())
	return err == nil
}

// Кастомный валидатор: дата в будущем
func futureDate(fl validator.FieldLevel) bool {
	parsedDate, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}
	today := time.Now().Truncate(24 * time.Hour)
	return parsedDate.After(today)
}

// ValidateCapsule выполняет валидацию модели CreateCapsule.
func ValidateCapsule(capsule models.CreateCapsule) error {
	err := validate.Struct(capsule)
	if err == nil {
		return nil
	}

	validationErrors := err.(validator.ValidationErrors)
	for _, ve := range validationErrors {
		switch ve.Tag() {
		case "send_at_date_format":
			return errors.New("invalid field `send_at` - Дата должна быть в формате YYYY-MM-DD")
		case "future_date":
			return errors.New("invalid field `send_at` - Дата должна быть больше текущей")
		case "uuid_or_null":
			return errors.New("invalid field `files_folder_uuid` - Поле должно быть валидным UUID или null")
		default:
			return errors.New("Некорректное значение в поле " + ve.Field())
		}
	}

	return nil
}

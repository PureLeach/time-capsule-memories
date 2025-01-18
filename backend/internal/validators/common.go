package validators

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var common_validate = validator.New()

// ValidateStruct выполняет валидацию структуры и возвращает описание ошибок
func ValidateStruct(data interface{}) error {
	if err := common_validate.Struct(data); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errorMessages []string
		for _, ve := range validationErrors {
			errorMessages = append(errorMessages, fmt.Sprintf("Поле '%s' не прошло валидацию: правило '%s'", ve.Field(), ve.Tag()))
		}
		return fmt.Errorf("%s", strings.Join(errorMessages, "; "))
	}
	return nil
}

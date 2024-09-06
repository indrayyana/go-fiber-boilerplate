package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func Password(field validator.FieldLevel) bool {
	value, ok := field.Field().Interface().(string)
	if ok {
		hasDigit := regexp.MustCompile(`[0-9]`).MatchString(value)
		hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(value)

		if !hasDigit || !hasLetter {
			return false
		}
	}

	return true
}

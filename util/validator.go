package util

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func IsNumeric(fl validator.FieldLevel) bool {
	re := regexp.MustCompile("^[0-9]+$")
	return re.MatchString(fl.Field().String())
}

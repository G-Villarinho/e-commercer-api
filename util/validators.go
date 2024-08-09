package util

import (
	"errors"
	"mime/multipart"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

func IsNumeric(fl validator.FieldLevel) bool {
	re := regexp.MustCompile("^[0-9]+$")
	return re.MatchString(fl.Field().String())
}

func ValidateFile(file *multipart.FileHeader) error {
	if file == nil {
		return errors.New("file is required")
	}

	ext := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, ".")+1:])
	validExtensions := []string{"jpg", "jpeg", "png", "gif", "bmp", "tiff", "webp"}

	isValidExtension := false
	for _, validExt := range validExtensions {
		if ext == validExt {
			isValidExtension = true
			break
		}
	}

	if !isValidExtension {
		return errors.New("invalid file type")
	}

	return nil
}

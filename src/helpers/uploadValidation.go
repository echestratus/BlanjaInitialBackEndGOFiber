package helpers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SizeUploadValidation(fileSize int64, maxFileSize int64) bool {
	if fileSize > maxFileSize {
		return true
	} else {
		return false
	}
}

func isValidFileType(validFileTypes []string, fileType string) bool {
	for _, validType := range validFileTypes {
		if validType == fileType {
			return true
		}
	}
	return false
}

func TypeUploadValidation(buffer []byte, validFileTypes []string) error {
	fileType := http.DetectContentType(buffer)
	if !isValidFileType(validFileTypes, fileType) {
		return fiber.NewError(fiber.StatusBadRequest, "File type is not valid. Only png, jpg, jpeg, and pdf are allowed")
	}
	return nil
}

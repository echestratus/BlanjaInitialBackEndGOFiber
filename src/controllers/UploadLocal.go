package controllers

import (
	"BlanjaInitialBackEndGOFiber/src/helpers"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func UploadLocal(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to upload file: " + err.Error())
	}

	maxFileSize := int64(2 << 20)
	if isSizeExceed := helpers.SizeUploadValidation(file.Size, maxFileSize); isSizeExceed {
		return c.Status(fiber.StatusBadRequest).SendString("File exceeds 2MB")
	}

	fileHeader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to open file: " + err.Error())
	}
	defer fileHeader.Close()

	buffer := make([]byte, 512)
	if _, err := fileHeader.Read(buffer); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read file: " + err.Error())
	}

	validateFileTypes := []string{"image/png", "image/jpeg", "image/jpg", "application/pdf"}
	if err := helpers.TypeUploadValidation(buffer, validateFileTypes); err != nil {
		return err
	}

	filePath := helpers.UploadFile(file)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save file: " + err.Error())
	}

	return c.SendString(fmt.Sprintf("%s File uploaded to %s", file.Filename, filePath))
}

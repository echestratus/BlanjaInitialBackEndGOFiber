package helpers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func UploadFile(file *multipart.FileHeader) string {
	uploadDir := "src/uploads"
	epochTime := time.Now().Unix()
	filePath := filepath.Join(uploadDir, fmt.Sprintf("%d_%s", epochTime, file.Filename))

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}
	return filePath
}

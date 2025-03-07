package service

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/fauzancodes/videoverse-api/app/pkg/upload"
)

func Upload(file *multipart.FileHeader, userID string) (responseURL string, statusCode int, err error) {
	allowedExtensions := []string{
		".png", ".jpg", ".jpeg", ".webp",
		".mp4", ".m4v", ".mkv",
	}
	extension := filepath.Ext(file.Filename)
	var isExtensionAllowed bool
	for _, item := range allowedExtensions {
		if extension == item {
			isExtensionAllowed = true
			break
		}
	}
	if !isExtensionAllowed {
		err = fmt.Errorf("the file extension is not allowed. allowed file extensions are %s", strings.Join(allowedExtensions, ", "))
		statusCode = http.StatusBadRequest
		return
	}

	var src multipart.File
	src, err = file.Open()
	if err != nil {
		err = fmt.Errorf("faield to open file: %s", err.Error())
		statusCode = http.StatusInternalServerError
		return
	}
	defer src.Close()

	responseURL, _, _, err = upload.UploadFile(src, userID, "")
	if err != nil {
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

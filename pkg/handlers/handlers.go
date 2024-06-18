package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/chigaji/file_sharing_service/pkg/security"
	"github.com/chigaji/file_sharing_service/pkg/storage"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func UploadHandler(c echo.Context) error {

	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error retrieving the file")
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error Opening the file")
	}
	defer src.Close()

	// generate a random file ID

	fileID := uuid.New().String()

	filePath := fmt.Sprintf("uploads/%s", fileID)

	// create the destination file
	dst, err := os.Create(filePath)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to create the file")
	}

	defer dst.Close()

	// Encrypt the file

	if err := security.Encrypt(src, dst); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to encrypt the file")
	}

	// save the file

	if err := storage.SaveFileData(fileID, file.Filename, time.Now().Add(24*time.Hour)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error Saving the file")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "File uploaded!",
		"FileID":  fileID,
	})
}

func DowloadHandler(c echo.Context) error {
	fileID := c.QueryParam("fileID")
	if fileID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "File ID is required")
	}

	fileData, err := storage.GetFileData(fileID)

	if err != nil || fileData.Expiry.Before(time.Now()) {
		return echo.NewHTTPError(http.StatusNotFound, "File not found or Expired")
	}

	encryptedFile, err := os.Open(fmt.Sprintf("uploads/%s", fileID))

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "File not found")
	}

	defer encryptedFile.Close()

	decryptedFile := security.Decrypt(encryptedFile)

	return c.Stream(http.StatusOK, "application/octet-stream", decryptedFile)
}

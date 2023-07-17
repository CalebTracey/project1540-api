package routes

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func accessFile(r *http.Request) (multipart.File, *multipart.FileHeader, error) {
	if formFile, header, err := r.FormFile("file"); err == nil {
		defer func(file multipart.File) {
			if closeErr := file.Close(); closeErr != nil {
				log.Errorf("accessFile: %v", closeErr)
			}
		}(formFile)
		return formFile, header, err
	} else {
		return formFile, header, fmt.Errorf("accessFile: %w", err)
	}
}

func createTempFile() (*os.File, error) {
	if tempFile, err := os.CreateTemp("", "temp-*"); err != nil {
		return nil, err
	} else {
		defer func(name string) {
			if removeErr := os.Remove(name); removeErr != nil {
				log.Errorf("createTempFile: %v", removeErr)
			}
		}(tempFile.Name())
		return tempFile, nil
	}
}

// copyFile copies the file data to the temporary file.
// on success, it returns the file info after opening the file for reading.
// needed for conversion of multipart.File to *os.File
func copyFile(tempFile *os.File, file multipart.File) (os.FileInfo, error) {
	if _, err := io.Copy(tempFile, file); err != nil {
		return nil, fmt.Errorf("copyFile: %w", err)
	} else {
		if fileInfo, statErr := tempFile.Stat(); statErr == nil {
			return fileInfo, nil
		} else {
			return fileInfo, fmt.Errorf("copyFile: %w", statErr)
		}
	}
}

// openUploadFile opens the file for reading
func openUploadFile(uploadFile *os.File) error {
	if _, err := uploadFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("openUploadFile: %w", err)
	}
	return nil // success
}

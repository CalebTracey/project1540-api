package routes

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func accessFile(r *http.Request) (multipart.File, *multipart.FileHeader, error) {
	if formFile, header, err := r.FormFile("file"); err == nil {
		defer func(file multipart.File) {
			_ = file.Close()
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
			_ = os.Remove(name)
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
		fileInfo, _ := tempFile.Stat()
		return fileInfo, nil
	}
}

// openUploadFile opens the file for reading
func openUploadFile(uploadFile *os.File) error {
	if _, err := uploadFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("openUploadFile: %w", err)
	}
	return nil // success
}

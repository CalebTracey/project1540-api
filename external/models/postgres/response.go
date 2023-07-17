package postgres

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project1540-api/external/models"
)

type FileResponse struct {
	Files   Files          `json:"files"`
	Message models.Message `json:"message"`
}

func (f *FileResponse) ToJSON(w http.ResponseWriter) error {
	if err := json.NewEncoder(w).Encode(f); err != nil {
		return fmt.Errorf("ToJSON: %w", err)
	}
	return nil
}

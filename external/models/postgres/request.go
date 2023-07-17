package postgres

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type NewFileRequest struct {
	Name string `json:"fileName"`
	Url  string `json:"fileUrl"`
}

func (r *NewFileRequest) FromJSON(req *http.Request) error {
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		return fmt.Errorf("FromJSON: %w", err)
	}
	return nil // success
}

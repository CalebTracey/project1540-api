package s3

import "project1540-api/external/models"

type DownloadS3Response struct {
	Message models.Message `json:"message,omitempty"`
}

type UploadS3Response struct {
	Message models.Message `json:"message,omitempty"`
}

type APIResponse struct {
	Results Responses      `json:"results,omitempty"`
	Message models.Message `json:"message,omitempty"`
}

type Responses []Response
type Response struct {
	ID string `json:"id"`
}

package models

import (
	"io"
)

type Request struct {
	ID string `json:"id,omitempty"`
}

type UploadS3 struct {
}

type InputFile struct {
	ID   string        `json:"id"`
	Name string        `json:"name,omitempty"`
	Size string        `json:"size,omitempty"`
	Type string        `json:"type,omitempty"`
	URL  string        `json:"url,omitempty"`
	Tags Tags          `json:"tags,omitempty"`
	File io.ReadCloser `json:"file,omitempty"`
}

type Tags []Tag

type Tag struct {
	Name string `json:"name"`
}

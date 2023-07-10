package models

type Request struct {
	ID string `json:"id,omitempty"`
}
type File struct {
	ID   string `json:"id"`
	Type string `json:"type,omitempty"`
	URL  string `json:"url,omitempty"`
	Tags Tags   `json:"tags,omitempty"`
}

type Tags []Tag

type Tag struct {
	Name string `json:"name"`
}

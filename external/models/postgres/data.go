package postgres

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

type Files []*File

type File struct {
	ID          *uuid.UUID `db:"id"`
	Name        string     `db:"name"`
	Tags        []string   `db:"tags"`
	CreatedDate *time.Time `db:"created_on"`
	UpdatedDate *time.Time `db:"updated_on"`
	URL         string     `db:"url"`
	Type        string     `db:"type"`
}

type FileOption func(*File)

func NewFile(options ...FileOption) *File {
	newFile := new(File)
	newFile.generateFileID()

	for _, option := range options {
		option(newFile)
	}

	createdDate := time.Now()
	newFile.CreatedDate = &createdDate

	return newFile
}

func WithName(name string) FileOption {
	return func(f *File) {
		f.Name = name
	}
}

func WithURL(url string) FileOption {
	return func(f *File) {
		f.URL = url
	}
}

func WithTags(tags []string) FileOption {
	return func(f *File) {
		f.Tags = tags
	}
}

func WithType(fileType string) FileOption {
	return func(f *File) {
		f.Type = fileType
	}
}

func (f *File) generateFileID() {
	if u, err := uuid.NewRandom(); err != nil {
		log.Errorf("generateFileID: %s", err.Error())
	} else {
		f.ID = &u
	}
}

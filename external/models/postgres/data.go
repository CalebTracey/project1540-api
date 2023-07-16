package postgres

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type File struct {
	ID   uuid.UUID
	Name string
}

type FileOption func(*File)

func New(options ...FileOption) (file *File) {
	file = new(File)
	file.generateFileID()

	for _, option := range options {
		option(file)
	}
	return file
}

func WithName(name string) FileOption {
	return func(f *File) {
		f.Name = name
	}
}

func (f *File) generateFileID() {
	if u, err := uuid.NewRandom(); err != nil {
		log.Errorf("generateFileID: %s", err.Error())
	} else {
		f.ID = u
	}
}

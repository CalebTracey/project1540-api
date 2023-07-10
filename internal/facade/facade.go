package facade

import (
	"context"
	log "github.com/sirupsen/logrus"
	"project1540-api/internal/dao/s3"
)

type IFacade interface {
	TestFacade(ctx context.Context) string
}

type Service struct {
	S3DAO s3.DAO
}

func (s Service) TestFacade(ctx context.Context) string {
	if err := s.S3DAO.PutObject(ctx); err != nil {
		log.Error(err)
		return ""
	}
	return "Test!"
}

package facade

import (
	"context"
	log "github.com/sirupsen/logrus"
	"project1540-api/external/models"
	daoS3 "project1540-api/internal/dao/s3"
)

type IFacade interface {
	UploadS3(ctx context.Context, input models.InputFile) *models.ErrorLog
}

type Service struct {
	S3DAO daoS3.DAO
}

func (s Service) UploadS3(ctx context.Context, input models.InputFile) *models.ErrorLog {
	if err := s.S3DAO.PutObject(ctx, input); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

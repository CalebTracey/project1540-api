package s3

import (
	"context"
	svcS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
	"project1540-api/external/models"
	"project1540-api/external/models/s3"
	daoS3 "project1540-api/internal/dao/s3"
)

//go:generate mockgen -source=facade.go -destination=mock/facade.go -package=s3
type IS3Facade interface {
	UploadS3Object(ctx context.Context, request s3.UploadS3Request) *models.ErrorLog
	DownloadS3Object(ctx context.Context, request s3.DownloadS3Request) (*svcS3.GetObjectOutput, *models.ErrorLog)
	GetS3ObjectNames(ctx context.Context, bucketName string) ([]string, *models.ErrorLog)
}

type Service struct {
	S3DAO daoS3.IDAO
}

func (s Service) UploadS3Object(ctx context.Context, request s3.UploadS3Request) *models.ErrorLog {
	// TODO: validate request
	if err := s.S3DAO.PutObject(ctx, request); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (s Service) DownloadS3Object(ctx context.Context, request s3.DownloadS3Request) (*svcS3.GetObjectOutput, *models.ErrorLog) {
	// TODO: validate request
	return s.S3DAO.GetObject(ctx, request)
}

func (s Service) GetS3ObjectNames(ctx context.Context, bucketName string) ([]string, *models.ErrorLog) {
	// TODO: validate request
	return s.S3DAO.GetAllObjectNames(ctx, bucketName)
}

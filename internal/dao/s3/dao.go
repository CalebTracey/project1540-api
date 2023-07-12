package s3

import (
	"context"
	svcS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"
	"net/http"
	"project1540-api/external/models"
	"project1540-api/external/models/s3"
)

type IDAO interface {
	PutObject(ctx context.Context, input s3.UploadS3Request) *models.ErrorLog
	GetObject(ctx context.Context, request s3.DownloadS3Request) (*svcS3.GetObjectOutput, *models.ErrorLog)
}

type DAO struct {
	*svcS3.Client
}

func (s DAO) PutObject(ctx context.Context, input s3.UploadS3Request) *models.ErrorLog {
	if _, err := s.Client.PutObject(
		ctx, &svcS3.PutObjectInput{
			Bucket: &input.DestBucket,
			Key:    &input.Name,
			Body:   input.File,
		},
	); err != nil {
		return &models.ErrorLog{
			StatusCode: http.StatusInternalServerError,
			RootCause:  err.Error(),
			Trace:      "PutObject",
		}
	}
	return nil //success
}

func (s DAO) GetObject(ctx context.Context, request s3.DownloadS3Request) (*svcS3.GetObjectOutput, *models.ErrorLog) {
	if object, err := s.Client.GetObject(
		ctx, &svcS3.GetObjectInput{
			Bucket: &request.BucketName,
			Key:    &request.FileName,
		},
	); err == nil {
		log.Printf("successfully downloaded file from %s/%s\n", request.BucketName, request.FileName)
		return object, nil
	} else {
		log.Error(err)
		return nil, &models.ErrorLog{
			RootCause:  err.Error(),
			Trace:      "GetObject",
			StatusCode: http.StatusInternalServerError,
		}
	}
}

package s3

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"project1540-api/external/models"
)

const (
	devBucket = "project1540-dev"
)

type IDAO interface {
	PutObject(ctx context.Context, input models.InputFile) *models.ErrorLog
	GetObject(ctx context.Context, input models.InputFile) (*os.File, *models.ErrorLog)
}

type DAO struct {
	*s3manager.Uploader
	*s3manager.Downloader
}

func (s DAO) PutObject(ctx context.Context, input models.InputFile) *models.ErrorLog {
	var bucket, key string
	//TODO: update
	bucket = devBucket
	key = input.Name

	if _, err := s.Uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   input.File,
	}); err != nil {

		if awsErr, isAwsError := err.(awserr.Error); isAwsError {
			return &models.ErrorLog{
				Status:     awsErr.Code(),
				StatusCode: getStatusCode(awsErr.Code()),
				RootCause:  awsErr.Error(),
				Trace:      "PutObject",
			}
		}
	}

	log.Printf("successfully uploaded file to %s/%s\n", bucket, key)
	return nil
}

func (s DAO) GetObject(ctx context.Context, input models.InputFile) (*os.File, *models.ErrorLog) {
	var bucket, key string
	//TODO: update
	bucket = devBucket
	key = input.Name

	file := new(os.File)
	//Writer
	if _, err := s.Downloader.DownloadWithContext(ctx, file, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}); err != nil {

		if awsErr, isAwsError := err.(awserr.Error); isAwsError {
			return file, &models.ErrorLog{
				Status:     awsErr.Code(),
				StatusCode: getStatusCode(awsErr.Code()),
				RootCause:  awsErr.Error(),
				Trace:      "GetObject",
			}
		}
	}

	log.Printf("successfully downloaded file from %s/%s\n", bucket, key)
	return file, nil
}

func getStatusCode(code string) int {
	switch code {
	case "AccessDenied":
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

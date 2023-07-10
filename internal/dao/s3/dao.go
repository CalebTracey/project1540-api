package s3

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
	"net/http"
	"project1540-api/external/models"
)

const (
	devBucket = "project1540-dev"
)

type IDAO interface {
	PutObject(ctx context.Context)
}

type DAO struct {
	*s3.S3
}

func (s DAO) PutObject(ctx context.Context) *models.ErrorLog {
	var bucket, key string
	//TODO: update
	bucket = devBucket
	key = "test key"

	if _, err := s.S3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		//Body:   os.Stdin,
	}); err != nil {

		if awsErr, ok := err.(awserr.Error); ok {

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

func getStatusCode(code string) int {
	switch code {
	case "AccessDenied":
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

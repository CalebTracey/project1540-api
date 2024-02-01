package s3

import (
	"context"
	"fmt"
	svcS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/calebtracey/project1540-api/external/models/s3"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -source=dao.go -destination=mock/dao.go -package=s3
type IDAO interface {
	PutObject(ctx context.Context, input s3.UploadS3Request) error
	GetObject(ctx context.Context, request s3.DownloadS3Request) (*svcS3.GetObjectOutput, error)
	GetAllObjectNames(ctx context.Context, bucketName string) ([]string, error)
}

type DAO struct {
	*svcS3.Client
}

func (s DAO) PutObject(ctx context.Context, input s3.UploadS3Request) error {
	if _, err := s.Client.PutObject(
		ctx, &svcS3.PutObjectInput{
			Bucket: &input.DestBucket,
			Key:    &input.Name,
			Body:   input.File,
		},
	); err != nil {
		return fmt.Errorf("PutObject: %w", err)
	}
	return nil // success
}

func (s DAO) GetObject(ctx context.Context, request s3.DownloadS3Request) (*svcS3.GetObjectOutput, error) {
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
		return nil, fmt.Errorf("GetObject: %w", err)
	}
}

func (s DAO) GetAllObjectNames(ctx context.Context, bucketName string) ([]string, error) {
	var objectNames []string
	paginator := svcS3.NewListObjectsV2Paginator(
		s.Client, &svcS3.ListObjectsV2Input{Bucket: &bucketName},
	)
	// paginate through the list of objects and collect object names
	for paginator.HasMorePages() {
		if resp, err := paginator.NextPage(ctx); err == nil {
			// add the object names from the current page
			for _, obj := range resp.Contents {
				objectNames = append(objectNames, *obj.Key)
			}
		} else {
			return nil, fmt.Errorf("GetAllObjectNames: %w", err)
		}
	}
	return objectNames, nil
}

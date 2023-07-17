package facade

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"project1540-api/external/models"
	postgres2 "project1540-api/external/models/postgres"
	"project1540-api/internal/facade/postgres"
	"project1540-api/internal/facade/s3"
	"project1540-api/internal/services/parser"
)

type IFacade interface {
	UpdateDatabaseFromS3Bucket(ctx context.Context, bucket string) *models.ErrorLog
	InsertNewFileByS3Bucket(ctx context.Context, req postgres2.NewFileRequest) *models.ErrorLog

	s3.IS3Facade
	postgres.IFacade
	parser.IParser
}

type Service struct {
	S3         s3.IS3Facade
	PostgresQL postgres.IFacade
	Parser     parser.IParser
}

func (s Service) UpdateDatabaseFromS3Bucket(ctx context.Context, bucket string) *models.ErrorLog {
	g, ctx := errgroup.WithContext(ctx)

	if objectNames, err := s.S3.GetS3ObjectNames(
		ctx, bucket,
	); err == nil {
		g.SetLimit(len(objectNames))

		for _, fileName := range objectNames {
			fileName := fileName

			g.Go(func() error {
				if tags, fileType, parseErr := s.Parser.ExtractTags(fileName); parseErr == nil {
					if postgresErr := s.PostgresQL.InsertNewFileDetails(
						ctx, fileName, fileType, "url", tags,
					); postgresErr == nil {
						return nil // success
					} else {
						return errors.New(fmt.Sprintf("UpdateDatabaseWithS3Data: %v", postgresErr))
					}
				} else {
					return fmt.Errorf("UpdateDatabaseFromS3Bucket: %w", parseErr)
				}
			})
		}

		if serviceErr := g.Wait(); serviceErr != nil {
			log.Error(err)
			return &models.ErrorLog{
				RootCause:  serviceErr.Error(),
				Trace:      "UpdateDatabaseFromS3Bucket",
				StatusCode: http.StatusInternalServerError,
			}
		}
	}
	return nil // success
}

func (s Service) InsertNewFileByS3Bucket(ctx context.Context, req postgres2.NewFileRequest) *models.ErrorLog {
	if tags, fileType, parseErr := s.Parser.ExtractTags(req.Name); parseErr == nil {
		if dbErr := s.PostgresQL.InsertNewFileDetails(
			ctx, req.Name, fileType, req.Url, tags,
		); dbErr != nil {
			dbErr.Trace = fmt.Sprintf("InsertNewFileByS3Bucket: %s", dbErr.Trace)
			return dbErr
		}
	} else {
		return &models.ErrorLog{
			RootCause:  parseErr.Error(),
			Trace:      "InsertNewFileByS3Bucket",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return nil // success
}

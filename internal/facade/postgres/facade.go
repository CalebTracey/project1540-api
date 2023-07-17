package postgres

import (
	"context"
	"net/http"
	"project1540-api/external/models"
	"project1540-api/external/models/postgres"
	postgresSrc "project1540-api/internal/dao/postgres"
)

type IFacade interface {
	InsertNewFileDetails(ctx context.Context, fileName, fileType, url string, tags []string) *models.ErrorLog
}

type Service struct {
	PSQLDAO postgresSrc.IDAO
}

func (s Service) InsertNewFileDetails(ctx context.Context, fileName, fileType, url string, tags []string) *models.ErrorLog {

	query := `INSERT INTO file.video_file (id, name, url, tags, type, created_on) VALUES ($1, $2, $3, $4, $5, $6)`

	if err := s.PSQLDAO.InsertOneFile(
		ctx, query, postgres.NewFile(
			postgres.WithName(fileName),
			postgres.WithURL(url),
			postgres.WithTags(tags),
			postgres.WithType(fileType),
		),
	); err != nil {
		return &models.ErrorLog{
			RootCause:  err.Error(),
			StatusCode: http.StatusInternalServerError,
			Status:     "ERROR",
			Trace:      "InsertNewFileDetails",
		}
	}

	return nil
}

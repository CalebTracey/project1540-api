package postgres

import (
	"context"
	"github.com/calebtracey/project1540-api/external/models"
	"github.com/calebtracey/project1540-api/external/models/postgres"
	postgresSrc "github.com/calebtracey/project1540-api/internal/dao/postgres"
	"net/http"
)

//go:generate mockgen -source=facade.go -destination=mock/facade.go -package=postgres
type IFacade interface {
	InsertNewFileDetails(ctx context.Context, fileName, fileType, url string, tags []string) *models.ErrorLog
	SearchFilesByTag(ctx context.Context, tags []string) postgres.FileResponse
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

func (s Service) SearchFilesByTag(ctx context.Context, tags []string) postgres.FileResponse {
	// query := `SELECT * FROM file.video_file WHERE $1 = ANY(tags)`
	if tags == nil {
		return postgres.FileResponse{
			Message: models.Message{
				ErrorLogs: models.ErrorLogs{
					{
						RootCause:  "SearchFilesByTag: tags parameter is required",
						StatusCode: http.StatusBadRequest,
						Status:     "ERROR",
						Trace:      "SearchFilesForTag",
					},
				},
			},
		}
	}

	if results, searchErr := s.PSQLDAO.SearchFilesByTag(ctx, SearchByTagQuery, tags); searchErr == nil {

		return postgres.FileResponse{Files: results}
	} else {
		return postgres.FileResponse{
			Message: models.Message{
				ErrorLogs: models.ErrorLogs{
					{
						RootCause:  searchErr.Error(),
						StatusCode: http.StatusInternalServerError,
						Status:     "ERROR",
						Trace:      "SearchFilesForTag",
					},
				},
			},
		}
	}
}

const SearchByTagQuery = `SELECT * FROM file.video_file WHERE tags && $1`

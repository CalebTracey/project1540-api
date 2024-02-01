package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/calebtracey/project1540-api/external/models/postgres"
)

//go:generate mockgen -source=dao.go -destination=mock/dao.go -package=postgres
type IDAO interface {
	InsertOneFile(ctx context.Context, query string, payload *postgres.File) error
	SearchFilesByTag(ctx context.Context, query string, tags []string) (files []*postgres.File, err error)
}

type DAO struct {
	*pgxpool.Pool
}

func (s DAO) InsertOneFile(ctx context.Context, query string, payload *postgres.File) error {
	if _, err := s.Pool.Exec(
		ctx, query,
		&payload.ID,
		&payload.Name,
		&payload.URL,
		&payload.Tags,
		&payload.Type,
		&payload.CreatedDate,
	); err == nil {
		return nil // success
	} else {
		return err
	}
}

func (s DAO) SearchFilesByTag(ctx context.Context, query string, tags []string) (files []*postgres.File, err error) {
	rows, err := s.Pool.Query(ctx, query, tags)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var file postgres.File
		if scanErr := rows.Scan(
			&file.ID,
			&file.Name,
			&file.Tags,
			&file.CreatedDate,
			&file.UpdatedDate,
			&file.URL,
			&file.Type,
		); scanErr != nil {
			return nil, scanErr
		}
		files = append(files, &file)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

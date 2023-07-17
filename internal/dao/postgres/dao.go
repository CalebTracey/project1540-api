package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"project1540-api/external/models/postgres"
)

type IDAO interface {
	InsertOneFile(ctx context.Context, query string, payload *postgres.File) error
}

type DAO struct {
	*pgxpool.Pool
}

func (s DAO) InsertOneFile(ctx context.Context, query string, payload *postgres.File) error {

	if result, err := s.Pool.Exec(
		ctx, query, &payload.ID, &payload.Name, &payload.URL, &payload.Tags, &payload.Type, &payload.CreatedDate,
	); err == nil {
		log.Infoln("InsertOneFile: success;")
		log.Infoln(result.String())
		log.Infoln(result.RowsAffected())
		return nil // success

	} else {
		return err
	}
}

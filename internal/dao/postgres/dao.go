package postgres

import (
	"context"
	"database/sql"
)

type IDAO interface {
	InsertOne(ctx context.Context, query string) error
}

type DAO struct {
	*sql.DB
}

func (s DAO) InsertOne(ctx context.Context, query string) error {

	return nil // success
}

package postgres

import (
	"project1540-api/internal/dao/postgres"
)

type IPostgresFacade interface {
}

type Service struct {
	PSQLDAO postgres.IDAO
}

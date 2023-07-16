package facade

import (
	"project1540-api/internal/facade/postgres"
	"project1540-api/internal/facade/s3"
)

type IFacade interface {
	s3.IS3Facade
	postgres.IPostgresFacade
}

type Service struct {
	S3   s3.IS3Facade
	PSQL postgres.IPostgresFacade
}

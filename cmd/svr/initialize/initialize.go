package initialize

import (
	"context"
	initialize "github.com/calebtracey/project1540-api/cmd/svr/initialize/s3"
	"github.com/calebtracey/project1540-api/internal/dao/postgres"
	daoS3 "github.com/calebtracey/project1540-api/internal/dao/s3"
	"github.com/calebtracey/project1540-api/internal/facade"
	facadePsql "github.com/calebtracey/project1540-api/internal/facade/postgres"
	facadeS3 "github.com/calebtracey/project1540-api/internal/facade/s3"
	"github.com/calebtracey/project1540-api/internal/routes"
	"github.com/calebtracey/project1540-api/internal/services/parser"
	"github.com/calebtraceyco/config"
	log "github.com/sirupsen/logrus"
)

const (
	postgresDB = "POSTGRES"
)

func NewService(ctx context.Context, cfg *config.Config) (service routes.Handler, err error) {
	var psqlConfig *config.DatabaseConfig

	if psqlConfig, err = cfg.Database(postgresDB); err == nil {
		log.Infof("established source connection: \"%s\"\n", postgresDB)
	} else {
		panic(err)
	}

	return routes.Handler{
		Service: &facade.Service{
			S3: facadeS3.Service{
				S3DAO: daoS3.DAO{
					Client: initialize.NewClient(ctx),
				},
			},
			PostgresQL: facadePsql.Service{
				PSQLDAO: postgres.DAO{
					Pool: psqlConfig.Pool,
				},
			},
			Parser: parser.Service{},
		},
	}, nil
}

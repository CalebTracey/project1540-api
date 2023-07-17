package initialize

import (
	"context"
	config "github.com/calebtracey/config-yaml"
	log "github.com/sirupsen/logrus"
	initialize "project1540-api/cmd/svr/initialize/s3"
	"project1540-api/internal/dao/postgres"
	daoS3 "project1540-api/internal/dao/s3"
	"project1540-api/internal/facade"
	facadePsql "project1540-api/internal/facade/postgres"
	facadeS3 "project1540-api/internal/facade/s3"
	"project1540-api/internal/parser"
	"project1540-api/internal/routes"
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

package initialize

import (
	"context"
	"fmt"
	config "github.com/calebtracey/config-yaml"
	initialize "project1540-api/cmd/svr/initialize/s3"
	"project1540-api/internal/dao/postgres"
	daoS3 "project1540-api/internal/dao/s3"
	"project1540-api/internal/facade"
	facadePsql "project1540-api/internal/facade/postgres"
	facadeS3 "project1540-api/internal/facade/s3"
	"project1540-api/internal/routes"
)

func NewService(ctx context.Context, cfg *config.Config) (service routes.Handler, err error) {
	var psqlConfig *config.DatabaseConfig

	if psqlConfig, err = cfg.Database("PSQL"); err != nil {
		return routes.Handler{}, fmt.Errorf("NewService: %w", err)
	}

	return routes.Handler{
		Service: &facade.Service{
			S3: facadeS3.Service{
				S3DAO: daoS3.DAO{
					Client: initialize.NewClient(ctx),
				},
			},
			PSQL: facadePsql.Service{
				PSQLDAO: postgres.DAO{
					DB: psqlConfig.DB,
				},
			},
		},
	}, nil
}

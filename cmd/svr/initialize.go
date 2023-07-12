package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	log "github.com/sirupsen/logrus"
	daoS3 "project1540-api/internal/dao/s3"
	"project1540-api/internal/facade"
	facadeS3 "project1540-api/internal/facade/s3"
	"project1540-api/internal/routes"
)

func initializeService(ctx context.Context) routes.Handler {
	return routes.Handler{
		Service: &facade.Service{
			S3: facadeS3.Service{
				S3DAO: daoS3.DAO{Client: s3Client(ctx)},
			},
		},
	}
}

func s3Client(ctx context.Context) *s3.Client {
	if cfg, err := config.LoadDefaultConfig(ctx); err == nil {
		return s3.NewFromConfig(cfg, s3ConfigOptions(cfg))
	} else {
		log.Panicln("failed to load AWS default config: %v", err)
	}
	return nil
}

type s3Options func(o *s3.Options)

func s3ConfigOptions(cfg aws.Config) s3Options {
	return func(o *s3.Options) {
		// make credentials with ARN from IAM S3 role
		o.Credentials = stscreds.NewAssumeRoleProvider(sts.NewFromConfig(cfg), roleARN)
		o.Region = "us-east-2"
	}
}

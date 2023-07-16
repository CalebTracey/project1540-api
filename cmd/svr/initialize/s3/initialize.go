package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/sirupsen/logrus"
)

func NewClient(ctx context.Context) *s3.Client {
	if cfg, err := config.LoadDefaultConfig(ctx); err == nil {
		return s3.NewFromConfig(cfg, s3ConfigOptions(cfg))
	} else {
		logrus.Panicln("failed to load AWS default config: %v", err)
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

const roleARN = "arn:aws:iam::128120887705:role/s3-dev"

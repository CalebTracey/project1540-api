package main

import (
	"github.com/NYTimes/gziphandler"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/calebtraceyco/http/server"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"os"
	s3DAO "project1540-api/internal/dao/s3"
	"project1540-api/internal/facade"
	"project1540-api/internal/routes"
	"time"
)

const (
	defaultPort = "8080"
	roleARN     = "arn:aws:iam::128120887705:role/s3-dev"
	portEnv     = "PORT"
)

func main() {
	defer panicQuit()
	var port string

	if port = os.Getenv(portEnv); port == "" {
		port = defaultPort
	}

	sess := session.Must(session.NewSession())
	s3Client := initializeS3(sess)

	handler := routes.Handler{
		Service: facade.Service{
			S3DAO: s3DAO.DAO{
				Uploader:   s3Uploader(sess, s3Client),
				Downloader: s3Downloader(sess, s3Client),
			},
		},
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	log.Fatal(server.ListenAndServe(
		port, "dev", gziphandler.GzipHandler(
			handler.InitializeRoutes(
				middleware.RequestID,
				middleware.RealIP,
				middleware.Logger,
				middleware.Recoverer,
				middleware.Timeout(60*time.Second),
			),
		)),
	)
}

func initializeS3(sess *session.Session) *s3.S3 {
	region := endpoints.UsEast2RegionID
	return s3.New(sess, &aws.Config{
		Region: &region,
		Credentials: stscreds.NewCredentials(
			sess, roleARN,
		),
	})
}

func s3Uploader(sess *session.Session, s3Client *s3.S3) *s3manager.Uploader {
	return s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 64 * 1024 * 1024 // 64MB per part
		u.S3 = s3Client
	})
}

func s3Downloader(sess *session.Session, s3Client *s3.S3) *s3manager.Downloader {
	return s3manager.NewDownloader(sess, func(d *s3manager.Downloader) {
		d.PartSize = 64 * 1024 * 1024 // 64MB per part
		d.S3 = s3Client
	})
}

func panicQuit() {
	if r := recover(); r != nil {
		log.Errorf("I panicked and am quitting: %v", r)
		log.Error("I should be alerting someone...")
	}
}

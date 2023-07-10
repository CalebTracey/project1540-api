package main

import (
	"github.com/NYTimes/gziphandler"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/calebtraceyco/http/server"
	log "github.com/sirupsen/logrus"
	"os"
	s3DAO "project1540-api/internal/dao/s3"
	"project1540-api/internal/facade"
	"project1540-api/internal/routes"
)

const (
	defaultPort = "8080"
	roleARN     = "arn:aws:iam::128120887705:role/s3-dev"
	PortEnv     = "PORT"
)

func main() {
	defer panicQuit()
	port := os.Getenv(PortEnv)

	if port == "" {
		port = defaultPort
	}

	handler := routes.Handler{
		Service: facade.Service{
			S3DAO: s3DAO.DAO{
				S3: initializeS3(),
			},
		},
	}

	router := handler.InitializeRoutes()

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(server.ListenAndServe(port, "dev", gziphandler.GzipHandler(router)))
}

func initializeS3() *s3.S3 {

	sess := session.Must(session.NewSession())
	region := endpoints.UsEast2RegionID

	return s3.New(sess, &aws.Config{
		Region: &region,
		Credentials: stscreds.NewCredentials(
			sess, roleARN,
		),
	})
	//return s3.New(sess)
}

func panicQuit() {
	if r := recover(); r != nil {
		log.Errorf("I panicked and am quitting: %v", r)
		log.Error("I should be alerting someone...")
	}
}

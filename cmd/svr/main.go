package main

import (
	"context"
	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi/v5/middleware"
	"time"

	"github.com/calebtraceyco/http/server"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	defer panicQuit()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if port = os.Getenv(portEnv); port == "" {
		port = defaultPort
	}

	handler := initializeService(ctx)

	log.Fatal(
		server.ListenAndServe(port, "dev", gziphandler.GzipHandler(
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

func panicQuit() {
	if r := recover(); r != nil {
		log.Errorf("I panicked and am quitting: %v", r)
		log.Error("I should be alerting someone...")
	}
}

var port string

const (
	defaultPort = "8080"
	roleARN     = "arn:aws:iam::128120887705:role/s3-dev"
	portEnv     = "PORT"
)

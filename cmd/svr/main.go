package main

import (
	"context"
	"github.com/NYTimes/gziphandler"
	config "github.com/calebtracey/config-yaml"
	"github.com/go-chi/chi/v5/middleware"
	"project1540-api/cmd/svr/initialize"
	"time"

	"github.com/calebtraceyco/http/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	defer panicQuit()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appConfig := config.New(configPath)

	if service, svcErr := initialize.NewService(ctx, appConfig); svcErr != nil {
		log.Panicln(svcErr)
	} else {
		log.Fatal(server.ListenAndServe(
			port, "dev", gziphandler.GzipHandler(
				service.InitializeRoutes(
					middleware.RequestID,
					middleware.RealIP,
					middleware.Logger,
					middleware.Recoverer,
					middleware.Timeout(60*time.Second),
				),
			),
		))
	}
}

func panicQuit() {
	if r := recover(); r != nil {
		log.Errorf("I panicked and am quitting: %v", r)
		log.Error("I should be alerting someone...")
	}
}

var port string

const (
	configPath  = "config.yaml"
	defaultPort = "8080"
	portEnv     = "PORT"
)

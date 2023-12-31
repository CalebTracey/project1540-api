package main

import (
	"context"
	"github.com/NYTimes/gziphandler"
	"github.com/calebtracey/project1540-api/cmd/svr/initialize"
	"github.com/calebtracey/project1540-api/internal/routes"
	"github.com/calebtraceyco/config"
	"github.com/go-chi/chi/v5/middleware"
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
		handler := routes.Handler(service)
		log.Fatal(server.ListenAndServe(
			appConfig.Port, appConfig.Env, gziphandler.GzipHandler(
				handler.InitializeRoutes(
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

const (
	configPath = "config.yaml"
)

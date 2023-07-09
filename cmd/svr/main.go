package main

import (
	"github.com/NYTimes/gziphandler"
	"github.com/calebtraceyco/http/server"
	log "github.com/sirupsen/logrus"
	"project1540-api/internal/facade"
	"project1540-api/internal/routes"
)

const (
	defaultPort = "8080"
)

func main() {
	defer panicQuit()

	handler := routes.Handler{
		Service: facade.Service{},
	}

	router := handler.InitializeRoutes()

	//router.Use(middlewa)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", defaultPort)
	log.Fatal(server.ListenAndServe(defaultPort, "dev", gziphandler.GzipHandler(router)))

}

func panicQuit() {
	if r := recover(); r != nil {
		log.Errorf("I panicked and am quitting: %v", r)
		log.Error("I should be alerting someone...")
	}
}

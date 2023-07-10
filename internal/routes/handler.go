package routes

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"project1540-api/graph"
	"project1540-api/graph/generated"
	"project1540-api/internal/facade"
	"time"
)

type Handler struct {
	Resolver graph.Resolver
	Service  facade.IFacade
}

func (h *Handler) InitializeRoutes() *chi.Mux {

	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	h.Resolver.IFacade = h.Service

	gqlServer := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &h.Resolver,
			},
		),
	)

	r.Handle("/graphql", gqlServer)
	r.Handle("/", playground.Handler("GraphQL playground", "/graphql"))

	r.Post("/put", h.PutHandler())

	return r
}

func (h *Handler) PutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		h.Service.TestFacade(r.Context())
	}
}

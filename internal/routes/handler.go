package routes

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"project1540-api/graph"
	"project1540-api/graph/generated"
	"project1540-api/internal/facade"
)

type Handler struct {
	Resolver graph.Resolver
	Service  facade.IFacade
}

func (h *Handler) InitializeRoutes() *chi.Mux {

	r := chi.NewRouter()

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

	return r
}

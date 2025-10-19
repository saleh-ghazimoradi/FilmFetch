package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/gateway/handlers"
	"net/http"
)

type HealthRoutes struct {
	healthHandler *handlers.HealthHandler
}

func (h *HealthRoutes) HealthRoute(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", h.healthHandler.HealthCheck)
}

func NewHealthRoute(healthHandler *handlers.HealthHandler) *HealthRoutes {
	return &HealthRoutes{healthHandler: healthHandler}
}

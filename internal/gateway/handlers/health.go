package handlers

import (
	"github.com/saleh-ghazimoradi/FilmFetch/config"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/helper"
	"log/slog"
	"net/http"
)

type HealthHandler struct {
	logger      *slog.Logger
	config      *config.Config
	customError helper.CustomError
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	env := helper.Envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": h.config.Application.Environment,
			"version":     h.config.Application.Version,
		},
	}

	if err := helper.WriteJSON(w, http.StatusOK, env, nil); err != nil {
		h.customError.ServerErrorResponse(w, r, err)
	}
}

func NewHealthHandler(config *config.Config, logger *slog.Logger, customError *helper.CustomError) *HealthHandler {
	return &HealthHandler{
		config:      config,
		logger:      logger,
		customError: *customError,
	}
}

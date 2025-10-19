package handlers

import (
	"fmt"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/domain"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/helper"
	"log/slog"
	"net/http"
	"time"
)

type MovieHandler struct {
	logger      *slog.Logger
	customError *helper.CustomError
}

func (m *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

func (m *MovieHandler) GetMovieById(w http.ResponseWriter, r *http.Request) {
	id, err := helper.ReadIdParam(r)
	if err != nil {
		m.customError.NotFoundResponse(w, r)
		return
	}

	movie := &domain.Movie{
		Id:        id,
		CreatedAt: time.Now(),
		Title:     "BlackList",
		Year:      2010,
		Runtime:   1560,
		Genres:    []string{"criminal", "FBI"},
		Version:   1,
	}

	if err = helper.WriteJSON(w, http.StatusOK, helper.Envelope{"movie": movie}, nil); err != nil {
		m.customError.ServerErrorResponse(w, r, err)
	}
}

func NewMovieHandler(logger *slog.Logger, customError *helper.CustomError) *MovieHandler {
	return &MovieHandler{
		logger:      logger,
		customError: customError,
	}
}

package handlers

import (
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/dto"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/helper"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/repository"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/service"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/validator"
	"log/slog"
	"net/http"
)

type MovieHandler struct {
	logger       *slog.Logger
	customError  *helper.CustomError
	movieService service.MovieService
}

func (m *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var payload *dto.Movie
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		m.customError.BadRequestResponse(w, r, err)
		return
	}

	v := validator.NewValidator()
	dto.ValidateMovie(v, payload)
	if !v.Valid() {
		m.customError.FailedValidationResponse(w, r, v.Errors)
		return
	}

	movie, err := m.movieService.CreateMovie(r.Context(), payload)
	if err != nil {
		m.customError.ServerErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.Id))

	if err = helper.WriteJSON(w, http.StatusCreated, helper.Envelope{"movie": movie}, headers); err != nil {
		m.customError.ServerErrorResponse(w, r, err)
	}
}

func (m *MovieHandler) GetMovieById(w http.ResponseWriter, r *http.Request) {
	id, err := helper.ReadIdParam(r)
	if err != nil {
		m.customError.NotFoundResponse(w, r)
		return
	}

	movie, err := m.movieService.GetMovieById(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			m.customError.NotFoundResponse(w, r)
		default:
			m.customError.ServerErrorResponse(w, r, err)
		}
		return
	}

	if err = helper.WriteJSON(w, http.StatusOK, helper.Envelope{"movie": movie}, nil); err != nil {
		m.customError.ServerErrorResponse(w, r, err)
	}
}

func (m *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	payload := &dto.QueryMovie{}
	v := validator.NewValidator()

	qs := r.URL.Query()
	payload.Title = helper.ReadString(qs, "title", "")
	payload.Genres = helper.ReadCSV(qs, "genres", []string{})
	payload.Filters.Page = helper.ReadInt(qs, "page", 1, v)
	payload.Filters.PageSize = helper.ReadInt(qs, "page_size", 20, v)
	payload.Filters.Sort = helper.ReadString(qs, "sort", "id")
	payload.Filters.SortSafeList = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	dto.ValidateFilters(v, payload.Filters)
	if !v.Valid() {
		m.customError.FailedValidationResponse(w, r, v.Errors)
		return
	}

	movies, metadata, err := m.movieService.GetMovies(r.Context(), payload)
	if err != nil {
		m.customError.ServerErrorResponse(w, r, err)
		return
	}

	if err = helper.WriteJSON(w, http.StatusOK, helper.Envelope{"movies": movies, "metadata": metadata}, nil); err != nil {
		m.customError.ServerErrorResponse(w, r, err)
	}
}

func (m *MovieHandler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	id, err := helper.ReadIdParam(r)
	if err != nil {
		m.customError.NotFoundResponse(w, r)
		return
	}

	var payload *dto.UpdateMovie
	if err = helper.ReadJSON(w, r, &payload); err != nil {
		m.customError.BadRequestResponse(w, r, err)
		return
	}

	v := validator.NewValidator()
	dto.ValidateUpdateMovie(v, payload)
	if !v.Valid() {
		m.customError.FailedValidationResponse(w, r, v.Errors)
		return
	}

	updatedMovie, err := m.movieService.UpdateMovie(r.Context(), id, payload)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrEditConflict):
			m.customError.EditConflictResponse(w, r)
		default:
			m.customError.ServerErrorResponse(w, r, err)
		}
		return
	}

	if err = helper.WriteJSON(w, http.StatusOK, helper.Envelope{"movie": updatedMovie}, nil); err != nil {
		m.customError.ServerErrorResponse(w, r, err)
	}
}

func (m *MovieHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := helper.ReadIdParam(r)
	if err != nil {
		m.customError.NotFoundResponse(w, r)
		return
	}

	if err = m.movieService.DeleteMovie(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			m.customError.NotFoundResponse(w, r)
		default:
			m.customError.ServerErrorResponse(w, r, err)
		}
		return
	}

	if err := helper.WriteJSON(w, http.StatusOK, helper.Envelope{"message": "movie successfully deleted"}, nil); err != nil {
		m.customError.ServerErrorResponse(w, r, err)
	}
}

func NewMovieHandler(logger *slog.Logger, customError *helper.CustomError, movieService service.MovieService) *MovieHandler {
	return &MovieHandler{
		logger:       logger,
		customError:  customError,
		movieService: movieService,
	}
}

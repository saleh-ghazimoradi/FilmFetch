package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/gateway/handlers"
	"net/http"
)

type MovieRoutes struct {
	movieHandler *handlers.MovieHandler
}

func (m *MovieRoutes) MovieRoute(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/v1/movies", m.movieHandler.CreateMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", m.movieHandler.GetMovieById)
	router.HandlerFunc(http.MethodPut, "/v1/movies/:id", m.movieHandler.UpdateMovie)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", m.movieHandler.DeleteMovie)
}

func NewMovieRoutes(movieHandler *handlers.MovieHandler) *MovieRoutes {
	return &MovieRoutes{
		movieHandler: movieHandler,
	}
}

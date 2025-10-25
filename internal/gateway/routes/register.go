package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/helper"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/middleware"
	"net/http"
)

type Register struct {
	customError  *helper.CustomError
	middleware   *middleware.Middleware
	healthRoutes *HealthRoutes
	movieRoutes  *MovieRoutes
	userRoutes   *UserRoutes
}

type Options func(*Register)

func WithCustomError(customError *helper.CustomError) Options {
	return func(r *Register) {
		r.customError = customError
	}
}

func WithMiddleware(middleware *middleware.Middleware) Options {
	return func(r *Register) {
		r.middleware = middleware
	}
}

func WithHealthRoutes(healthRoutes *HealthRoutes) Options {
	return func(r *Register) {
		r.healthRoutes = healthRoutes
	}
}

func WithMovieRoutes(movieRoutes *MovieRoutes) Options {
	return func(r *Register) {
		r.movieRoutes = movieRoutes
	}
}

func WithUserRoutes(userRoutes *UserRoutes) Options {
	return func(r *Register) {
		r.userRoutes = userRoutes
	}
}

func (r *Register) RegisterRoutes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(r.customError.NotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(r.customError.MethodNotAllowedResponse)

	r.healthRoutes.HealthRoute(router)
	r.movieRoutes.MovieRoute(router)
	r.userRoutes.UserRoutes(router)

	return r.middleware.RecoverPanic(r.middleware.RateLimit(router))
}

func NewRegister(opts ...Options) *Register {
	r := &Register{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

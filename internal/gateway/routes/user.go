package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/gateway/handlers"
	"net/http"
)

type UserRoutes struct {
	userHandler *handlers.UserHandler
}

func (u *UserRoutes) UserRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/v1/users", u.userHandler.CreateUser)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", u.userHandler.ActivateUser)
}

func NewUserRoutes(userHandler *handlers.UserHandler) *UserRoutes {
	return &UserRoutes{
		userHandler: userHandler,
	}
}

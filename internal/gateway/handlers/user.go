package handlers

import (
	"errors"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/dto"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/helper"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/repository"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/service"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/validator"
	"github.com/saleh-ghazimoradi/FilmFetch/utils/email"
	"log/slog"
	"net/http"
)

type UserHandler struct {
	customErr   *helper.CustomError
	userService service.UserService
	mailService email.MailSender
	logger      *slog.Logger
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload *dto.User

	if err := helper.ReadJSON(w, r, &payload); err != nil {
		u.customErr.BadRequestResponse(w, r, err)
		return
	}

	v := validator.NewValidator()
	dto.ValidateUser(v, payload)
	if !v.Valid() {
		u.customErr.FailedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := u.userService.CreateUser(r.Context(), payload)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			u.customErr.FailedValidationResponse(w, r, v.Errors)
		default:
			u.customErr.ServerErrorResponse(w, r, err)
		}
		return
	}

	helper.Background(func() {
		if err := u.mailService.Send(user.Email, "user_welcome.tmpl", user); err != nil {
			u.logger.Error(err.Error())
		}
	})

	if err = helper.WriteJSON(w, http.StatusAccepted, helper.Envelope{"user": user}, nil); err != nil {
		u.customErr.ServerErrorResponse(w, r, err)
	}
}

func NewUserHandler(customErr *helper.CustomError, userService service.UserService, mailService email.MailSender, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		customErr:   customErr,
		userService: userService,
		mailService: mailService,
		logger:      logger,
	}
}

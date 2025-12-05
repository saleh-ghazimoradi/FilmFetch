package handlers

import (
	"errors"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/domain"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/dto"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/helper"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/repository"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/service"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/validator"
	"github.com/saleh-ghazimoradi/FilmFetch/utils/email"
	"log/slog"
	"net/http"
	"time"
)

type UserHandler struct {
	customErr    *helper.CustomError
	userService  service.UserService
	tokenService service.TokenService
	mailService  email.MailSender
	logger       *slog.Logger
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

	token, err := u.tokenService.Tokenize(r.Context(), user.Id, 3*24*time.Hour, domain.ScopeActivation)
	if err != nil {
		u.customErr.ServerErrorResponse(w, r, err)
		return
	}

	helper.Background(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userId":          user.Id,
		}

		if err := u.mailService.Send(user.Email, "user_welcome.tmpl", data); err != nil {
			u.logger.Error(err.Error())
		}
	})

	if err = helper.WriteJSON(w, http.StatusAccepted, helper.Envelope{"user": user}, nil); err != nil {
		u.customErr.ServerErrorResponse(w, r, err)
	}
}

func (u *UserHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	var payload *dto.ActivateUser

	if err := helper.ReadJSON(w, r, &payload); err != nil {
		u.customErr.BadRequestResponse(w, r, err)
		return
	}

	v := validator.NewValidator()
	dto.ValidateTokenPlaintext(v, payload)
	if !v.Valid() {
		u.customErr.FailedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := u.tokenService.ActivateUser(r.Context(), payload)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			u.customErr.FailedValidationResponse(w, r, v.Errors)
		case errors.Is(err, repository.ErrEditConflict):
			u.customErr.EditConflictResponse(w, r)
		case err != nil && err.Error() == "validation failed":
			u.customErr.FailedValidationResponse(w, r, v.Errors)
		default:
			u.customErr.ServerErrorResponse(w, r, err)
		}
		return
	}

	if err := helper.WriteJSON(w, http.StatusOK, helper.Envelope{"user": user}, nil); err != nil {
		u.customErr.ServerErrorResponse(w, r, err)
	}
}

func NewUserHandler(customErr *helper.CustomError, userService service.UserService, tokenService service.TokenService, mailService email.MailSender, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		customErr:    customErr,
		userService:  userService,
		tokenService: tokenService,
		mailService:  mailService,
		logger:       logger,
	}
}

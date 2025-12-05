package dto

import (
	"github.com/saleh-ghazimoradi/FilmFetch/internal/validator"
)

type ActivateUser struct {
	TokenPlaintext string `json:"token"`
}

func ValidateTokenPlaintext(v *validator.Validator, user *ActivateUser) {
	v.Check(user.TokenPlaintext != "", "token", "must be provided")
	v.Check(len(user.TokenPlaintext) == 26, "token", "must contain 26 characters")
}

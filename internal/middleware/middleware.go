package middleware

import (
	"fmt"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/helper"
	"net/http"
)

type Middleware struct {
	customError *helper.CustomError
}

func (m *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				m.customError.ServerErrorResponse(w, r, fmt.Errorf("%v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func NewMiddleware(customError *helper.CustomError) *Middleware {
	return &Middleware{
		customError: customError,
	}
}

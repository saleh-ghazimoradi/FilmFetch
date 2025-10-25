package helper

import (
	"fmt"
	"log/slog"
	"net/http"
)

type CustomError struct {
	logger *slog.Logger
}

func (c *CustomError) LogError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	c.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (c *CustomError) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := Envelope{"error": message}

	if err := WriteJSON(w, status, env, nil); err != nil {
		c.LogError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *CustomError) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	c.LogError(r, err)
	message := "the server encountered a problem and could not process your request"
	c.ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func (c *CustomError) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	c.ErrorResponse(w, r, http.StatusNotFound, message)
}

func (c *CustomError) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	c.ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (c *CustomError) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	c.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (c *CustomError) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	c.ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (c *CustomError) EditConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	c.ErrorResponse(w, r, http.StatusConflict, message)
}

func (c *CustomError) RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded, please try again"
	c.ErrorResponse(w, r, http.StatusTooManyRequests, message)
}

func NewCustomErr(logger *slog.Logger) *CustomError {
	return &CustomError{
		logger: logger,
	}
}

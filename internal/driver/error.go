package driver

import (
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	StatusCode int
	Err        string
	Message    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v - %v - %v", e.StatusCode, e.Err, e.Message)
}

func (e *Error) Is(target error) bool {
	var restErr *Error
	if !errors.As(target, &restErr) {
		return false
	}
	return *e == *restErr
}

func NewInternalServerError(err error) *Error {
	return &Error{
		StatusCode: http.StatusInternalServerError,
		Err:        "ERR_INTERNAL_ERROR",
		Message:    err.Error(),
	}
}

func NewBadRequestError(msg string) *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,
		Err:        "ERR_BAD_REQUEST",
		Message:    msg,
	}
}

func NewInvalidCredsError() *Error {
	return &Error{
		StatusCode: http.StatusUnauthorized,
		Err:        "ERR_INVALID_CREDENTIALS",
		Message:    "invalid credentials",
	}
}

func NewNotFoundError(msg string) *Error {
	return &Error{
		StatusCode: http.StatusNotFound,
		Err:        "ERR_MESSAGES_NOT_FOUND",
		Message:    msg,
	}
}

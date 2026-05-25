package core

import (
	"errors"
)

var (
	ErrCodeBadRequest    = errors.New("ERR_BAD_REQUEST")
	ErrCodeInternalError = errors.New("ERR_INTERNAL_ERROR")
	ErrSessionExpired    = errors.New("ERR_SESSION_EXPIRED")
	ErrMessageNotFound   = errors.New("ERR_MESSAGE_NOT_FOUND")
	ErrRetryNonFailed    = errors.New("ERR_RETRY_NON_FAILED")
)

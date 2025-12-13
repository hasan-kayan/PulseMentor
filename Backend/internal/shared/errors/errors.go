package errors

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrInvalidInput  = errors.New("invalid input")
	ErrAlreadyExists = errors.New("already exists")
	ErrInternal      = errors.New("internal error")
)


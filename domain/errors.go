package domain

import (
	"errors"
)

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal Server Error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested Item is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("your Item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("given Param is not valid")
)

// RequestErr is meant to be the error returned from the usecases
type RequestErr interface {
	GetStatus() int
	error
}

type uCaseErr struct {
	Status int
	Err    error
}

func (u uCaseErr) GetStatus() int {
	return u.Status
}

func (u uCaseErr) Error() string {
	return u.Err.Error()
}

// Create a basic usecase error
func NewUCaseErr(status int, err error) uCaseErr {
	return uCaseErr{status, err}
}

package rest

import "net/http"

const (
	BadMethodMessage  = "method not allowed"
	BadRequestMessage = "bad request"
)

type CodeError interface {
	error
	Code() int
}

type Error struct {
	error
	code int
}

func NewError(err error, code int) *Error { return &Error{err, code} }
func (e *Error) Code() int                { return e.code }

func NewNotFoundError(err error) *Error { return NewError(err, http.StatusNotFound) }
func NewInternalError(err error) *Error { return NewError(err, http.StatusInternalServerError) }

type BadRequestError struct{}

func (e *BadRequestError) Code() int     { return http.StatusBadRequest }
func (e *BadRequestError) Error() string { return BadRequestMessage }

type MethodNotAllowedError struct{}

func (e *MethodNotAllowedError) Code() int     { return http.StatusMethodNotAllowed }
func (e *MethodNotAllowedError) Error() string { return BadMethodMessage }

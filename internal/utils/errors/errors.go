package errors

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// Error holds http status and message values
type Error struct {
	httpStatus int
	message    string
}

// Error returns message
func (e *Error) Error() string {
	return e.message
}

// Status returns http status code
func (e *Error) Status() int {
	return e.httpStatus
}

// Option to modify Error message
type Option func(e *Error)

// WithReason adds reason to Error message
func WithReason(reason string) Option {
	return func(e *Error) {
		e.message = fmt.Sprintf("%s %s", e.message, reason)
	}
}

// WithMessage replaces Error message
func WithMessage(msg string) Option {
	return func(e *Error) {
		e.message = msg
	}
}

// New creates a copy from Error with Options to modify message
func (e *Error) New(options ...Option) *Error {
	err := &Error{
		httpStatus: e.httpStatus,
		message:    e.message,
	}

	for _, option := range options {
		option(err)
	}

	return err
}

// New returns Error with specified http status and message
func New(httpStatus int, msg string) *Error {
	return &Error{
		httpStatus: httpStatus,
		message:    msg,
	}
}

const (
	defaultErrorMessage4xx = "Problem with user request. Please make sure the information you send is correct."
	defaultTimeOutMessage  = "client closed connection."
	defaultErrorMessage5xx = "Problem with connection to server. Please try again or contact relevant person-in-charge."
)

// https://juloprojects.atlassian.net/browse/J360-227
var (
	// 4xx
	BadRequest   = New(http.StatusBadRequest, defaultErrorMessage4xx)
	Unauthorized = New(http.StatusUnauthorized, defaultErrorMessage4xx)
	NotFound     = New(http.StatusNotFound, defaultErrorMessage4xx)
	Conflict     = New(http.StatusConflict, defaultErrorMessage4xx)

	// timeout
	ClientClosedConnection = New(499, defaultTimeOutMessage)

	// 5xx
	InternalServerError = New(http.StatusInternalServerError, defaultErrorMessage5xx)
)

// As acts as go errors.As() wrapper to target Error. defaults to InternalServerError
func As(err error) *Error {
	if err == nil {
		return InternalServerError
	}

	if errors.Is(err, context.Canceled) {
		return ClientClosedConnection
	}

	target := InternalServerError
	errors.As(err, &target)
	return target
}

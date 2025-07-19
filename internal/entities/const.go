package entities

import "errors"

const (
	StatusActive   = 1
	StatusInactive = 0
	StatusDeleted  = 2
)

// ErrInvalidInput returns an error indicating invalid input
func ErrInvalidInput(msg string) error {
	return errors.New(msg)
}

var ErrNotFound = errors.New("record not found")
var SqlNoRows = errors.New("no rows in result set")

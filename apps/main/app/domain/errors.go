package domain

import "errors"

var (
	// ErrNotFound is returned when a requested user or hair history does not exist.
	ErrNotFound = errors.New("not found")
	// ErrInvalidInput is returned for malformed IDs or other bad client input.
	ErrInvalidInput = errors.New("invalid input")
)

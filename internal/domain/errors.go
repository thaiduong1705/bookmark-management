package domain

import "errors"

var (
	// ErrNotFound is returned when not found.
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists is returned when a unique constraint rejects a record.
	ErrAlreadyExists = errors.New("already exists")
	// ErrInvalidInput is returned when the database rejects the data as invalid
	// (not null, check constraint, value too long...).
	ErrInvalidInput = errors.New("invalid input")
)

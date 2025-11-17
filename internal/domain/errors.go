package domain

import "errors"

var (
	// ErrNotFound represents missing entities.
	ErrNotFound = errors.New("not found")
	// ErrConflict represents unique constraint violations.
	ErrConflict = errors.New("conflict")
	// ErrInvalidCredentials indicates authentication failure.
	ErrInvalidCredentials = errors.New("invalid credentials")
)

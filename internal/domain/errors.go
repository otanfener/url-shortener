package domain

import "errors"

var (
	//ErrShortCodeNotFound is returned when short code is not found
	ErrShortCodeNotFound = errors.New("short code not found")
	//ErrShortCodeExists is returned when short code already exists
	ErrShortCodeExists = errors.New("short code already exists")
	// ErrStorageFailure is returned when storage operation fails
	ErrStorageFailure = errors.New("storage failure")
	// ErrCounterFailure is returned when counter operation fails
	ErrCounterFailure = errors.New("counter failure")
	// ErrInvalidInput is returned when input is invalid
	ErrInvalidInput = errors.New("invalid input")
	// ErrInternal is returned when there is internal server error
	ErrInternal = errors.New("internal error")
)

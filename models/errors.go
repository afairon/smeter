package models

import (
	"errors"
)

var (
	// ErrInvalidName means that the name is invalid.
	ErrInvalidName = errors.New("models: invalid name")

	// ErrInvalidReq means that the ID is invalid.
	ErrInvalidReq = errors.New("models: invalid request")

	// ErrNotFound means that the record is not found.
	ErrNotFound = errors.New("models: not found")

	// ErrUnknownSensorType means that the type is set to 0.
	ErrUnknownSensorType = errors.New("models: unknown sensor type")
)

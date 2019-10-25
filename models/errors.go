package models

import (
	"errors"
)

var (
	// ErrInvalidName means that the name is invalid.
	ErrInvalidName = errors.New("models: invalid name")
)

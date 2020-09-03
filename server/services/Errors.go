package services

import "errors"

var (
	ErrAlreadyInitialized = errors.New("Service already initialized")
	ErrNotInitialized     = errors.New("Service not yet initialized")
)

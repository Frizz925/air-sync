package services

import (
	"errors"
	"fmt"
)

var (
	ErrAlreadyInitialized = errors.New("Service already initialized")
	ErrNotInitialized     = errors.New("Service not yet initialized")
)

type CronRequestError struct {
	error
}

func NewCronRequestError(format string, a ...interface{}) CronRequestError {
	return CronRequestError{
		error: fmt.Errorf(format, a...),
	}
}

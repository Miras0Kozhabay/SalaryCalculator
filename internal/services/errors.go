package services

import "errors"

// Typed errors for proper error handling in handlers
var (
	ErrInvalidInput = errors.New("invalid input")
	ErrCalculation  = errors.New("calculation failed")
)

package services

import (
	"errors"
	"fmt"
)

// Typed errors for proper error handling in handlers
var (
	ErrInvalidInput = errors.New("invalid input")
	ErrCalculation  = errors.New("calculation failed")
)

// WrapCalculationError wraps a calculation error while preserving the original message
func WrapCalculationError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: %v", ErrCalculation, err)
}

package validator

import "errors"

var (
	ErrRequired    = errors.New("required field is missing")
	ErrValueTooSmall = errors.New("value is too small")
	ErrValueTooLarge = errors.New("value is too large")
	ErrInvalidFormat = errors.New("invalid format")
)


package validator

import (
	"strconv"
	"strings"
)

// ValidateInt64 int64 값 검증
func ValidateInt64(value string, min, max int64) (int64, error) {
	if value == "" {
		return 0, nil
	}
	
	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	
	if min > 0 && val < min {
		return 0, ErrValueTooSmall
	}
	
	if max > 0 && val > max {
		return 0, ErrValueTooLarge
	}
	
	return val, nil
}

// ValidateRequired 필수 값 검증
func ValidateRequired(value string) error {
	if strings.TrimSpace(value) == "" {
		return ErrRequired
	}
	return nil
}

// ValidateRange 범위 검증
func ValidateRange(value, min, max int64) error {
	if value < min {
		return ErrValueTooSmall
	}
	if value > max {
		return ErrValueTooLarge
	}
	return nil
}


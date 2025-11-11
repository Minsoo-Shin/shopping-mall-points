package errors

import (
	"fmt"
	"net/http"
)

// AppError 애플리케이션 에러
type AppError struct {
	Code    int
	Message string
	Err     error
}

// Error 에러 메시지 반환
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewAppError 새로운 애플리케이션 에러 생성
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewBadRequestError Bad Request 에러 생성
func NewBadRequestError(message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, message, err)
}

// NewNotFoundError Not Found 에러 생성
func NewNotFoundError(message string, err error) *AppError {
	return NewAppError(http.StatusNotFound, message, err)
}

// NewInternalServerError Internal Server Error 생성
func NewInternalServerError(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, message, err)
}


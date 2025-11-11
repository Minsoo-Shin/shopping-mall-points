package point

import "errors"

var (
	// ErrInsufficientPoints 보유 포인트 부족
	ErrInsufficientPoints = errors.New("insufficient points")
	
	// ErrBelowMinUseAmount 최소 사용 금액 미만
	ErrBelowMinUseAmount = errors.New("below minimum use amount")
	
	// ErrInvalidUseUnit 사용 단위 오류
	ErrInvalidUseUnit = errors.New("invalid use unit")
	
	// ErrExceedMaxUseRate 최대 사용 비율 초과
	ErrExceedMaxUseRate = errors.New("exceed maximum use rate")
	
	// ErrBelowMinPayment 최소 결제 금액 미만
	ErrBelowMinPayment = errors.New("below minimum payment amount")
	
	// ErrPointNotFound 포인트 정보 없음
	ErrPointNotFound = errors.New("point not found")
	
	// ErrTransactionNotFound 거래 내역 없음
	ErrTransactionNotFound = errors.New("transaction not found")
)


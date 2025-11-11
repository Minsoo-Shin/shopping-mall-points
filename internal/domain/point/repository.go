package point

import (
	"context"
	"time"
)

// Repository 포인트 리포지토리 인터페이스
type Repository interface {
	// GetUserPoint 사용자 포인트 조회 (락 포함)
	GetUserPoint(ctx context.Context, userID int64) (*UserPoint, error)
	
	// UpdateUserPoint 사용자 포인트 업데이트
	UpdateUserPoint(ctx context.Context, userPoint *UserPoint) error
	
	// CreateTransaction 거래 내역 생성
	CreateTransaction(ctx context.Context, tx *Transaction) error
	
	// GetEarnedTransactions 적립 거래 내역 조회 (FIFO용, 만료일 순)
	GetEarnedTransactions(ctx context.Context, userID int64, limit int) ([]*Transaction, error)
	
	// UpdateTransaction 거래 내역 업데이트
	UpdateTransaction(ctx context.Context, tx *Transaction) error
	
	// GetExpiringTransactions 만료 예정 거래 내역 조회
	GetExpiringTransactions(ctx context.Context, before time.Time, limit int) ([]*Transaction, error)
	
	// GetTransactionsByUser 사용자 거래 내역 조회
	GetTransactionsByUser(ctx context.Context, userID int64, limit, offset int) ([]*Transaction, error)
	
	// GetTransactionByID 거래 내역 ID로 조회
	GetTransactionByID(ctx context.Context, id int64) (*Transaction, error)
	
	// GetTransactionsByOrderID 주문 ID로 거래 내역 조회
	GetTransactionsByOrderID(ctx context.Context, orderID int64) ([]*Transaction, error)
}

// TransactionManager 트랜잭션 관리자 인터페이스
type TransactionManager interface {
	// WithTransaction 트랜잭션 내에서 함수 실행
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
}


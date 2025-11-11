package mysql

import (
	"context"
	"database/sql"
)

// TransactionManager 트랜잭션 관리자
type TransactionManager struct {
	db *sql.DB
}

// NewTransactionManager 트랜잭션 관리자 생성
func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// WithTransaction 트랜잭션 내에서 함수 실행
func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// 트랜잭션 컨텍스트 생성
	txCtx := context.WithValue(ctx, "tx", tx)

	// 함수 실행
	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	// 커밋
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetTx 컨텍스트에서 트랜잭션 추출
func GetTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value("tx").(*sql.Tx); ok {
		return tx
	}
	return nil
}

// GetDBOrTx 컨텍스트에서 트랜잭션이 있으면 반환, 없으면 DB 반환
func (tm *TransactionManager) GetDBOrTx(ctx context.Context) interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
} {
	if tx := GetTx(ctx); tx != nil {
		return tx
	}
	return tm.db
}

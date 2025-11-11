package mysql

import (
	"context"
	"database/sql"
	"shopping-mall/internal/domain/point"
	"time"
)

// PointRepository 포인트 리포지토리 구현
type PointRepository struct {
	tm *TransactionManager
}

// NewPointRepository 포인트 리포지토리 생성
func NewPointRepository(tm *TransactionManager) *PointRepository {
	return &PointRepository{tm: tm}
}

// GetUserPoint 사용자 포인트 조회 (락 포함)
func (r *PointRepository) GetUserPoint(ctx context.Context, userID int64) (*point.UserPoint, error) {
	query := `
		SELECT user_id, available_balance, pending_balance, total_earned, total_used, updated_at
		FROM user_points
		WHERE user_id = ?
		FOR UPDATE
	`

	db := r.tm.GetDBOrTx(ctx)
	row := db.QueryRowContext(ctx, query, userID)

	var up point.UserPoint
	var updatedAt time.Time
	err := row.Scan(
		&up.UserID,
		&up.AvailableBalance,
		&up.PendingBalance,
		&up.TotalEarned,
		&up.TotalUsed,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, point.ErrPointNotFound
	}
	if err != nil {
		return nil, err
	}

	up.UpdatedAt = updatedAt
	return &up, nil
}

// CreateUserPoint 사용자 포인트 생성
func (r *PointRepository) CreateUserPoint(ctx context.Context, userPoint *point.UserPoint) error {
	query := `
		INSERT INTO user_points (user_id, available_balance, pending_balance, total_earned, total_used, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	
	db := r.tm.GetDBOrTx(ctx)
	_, err := db.ExecContext(ctx, query,
		userPoint.UserID,
		userPoint.AvailableBalance,
		userPoint.PendingBalance,
		userPoint.TotalEarned,
		userPoint.TotalUsed,
		time.Now(),
	)
	return err
}

// UpdateUserPoint 사용자 포인트 업데이트
func (r *PointRepository) UpdateUserPoint(ctx context.Context, userPoint *point.UserPoint) error {
	query := `
		UPDATE user_points
		SET available_balance = ?, pending_balance = ?, total_earned = ?, total_used = ?, updated_at = ?
		WHERE user_id = ?
	`

	db := r.tm.GetDBOrTx(ctx)
	_, err := db.ExecContext(ctx, query,
		userPoint.AvailableBalance,
		userPoint.PendingBalance,
		userPoint.TotalEarned,
		userPoint.TotalUsed,
		time.Now(),
		userPoint.UserID,
	)
	return err
}

// CreateTransaction 거래 내역 생성
func (r *PointRepository) CreateTransaction(ctx context.Context, tx *point.Transaction) error {
	query := `
		INSERT INTO point_transactions 
		(user_id, transaction_type, amount, balance_after, reason_type, reason_detail, 
		 order_id, earned_at, expires_at, expired, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	db := r.tm.GetDBOrTx(ctx)
	result, err := db.ExecContext(ctx, query,
		tx.UserID,
		tx.Type,
		tx.Amount,
		tx.BalanceAfter,
		tx.ReasonType,
		tx.ReasonDetail,
		tx.OrderID,
		tx.EarnedAt,
		tx.ExpiresAt,
		tx.Expired,
		tx.Status,
		time.Now(),
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	tx.ID = id
	return nil
}

// GetEarnedTransactions 적립 거래 내역 조회 (FIFO용, 만료일 순)
func (r *PointRepository) GetEarnedTransactions(ctx context.Context, userID int64, limit int) ([]*point.Transaction, error) {
	query := `
		SELECT id, user_id, transaction_type, amount, balance_after, reason_type, reason_detail,
		       order_id, earned_at, expires_at, expired, status, created_at
		FROM point_transactions
		WHERE user_id = ? 
		  AND transaction_type = 'EARN'
		  AND expired = false
		  AND status = 'CONFIRMED'
		ORDER BY expires_at ASC, created_at ASC
		LIMIT ?
	`

	db := r.tm.GetDBOrTx(ctx)
	rows, err := db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*point.Transaction
	for rows.Next() {
		var tx point.Transaction
		var earnedAt, expiresAt sql.NullTime
		var orderID sql.NullInt64

		err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.Type,
			&tx.Amount,
			&tx.BalanceAfter,
			&tx.ReasonType,
			&tx.ReasonDetail,
			&orderID,
			&earnedAt,
			&expiresAt,
			&tx.Expired,
			&tx.Status,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if orderID.Valid {
			tx.OrderID = &orderID.Int64
		}
		if earnedAt.Valid {
			tx.EarnedAt = &earnedAt.Time
		}
		if expiresAt.Valid {
			tx.ExpiresAt = &expiresAt.Time
		}

		transactions = append(transactions, &tx)
	}

	return transactions, rows.Err()
}

// UpdateTransaction 거래 내역 업데이트
func (r *PointRepository) UpdateTransaction(ctx context.Context, tx *point.Transaction) error {
	query := `
		UPDATE point_transactions
		SET expired = ?, status = ?
		WHERE id = ?
	`

	db := r.tm.GetDBOrTx(ctx)
	_, err := db.ExecContext(ctx, query, tx.Expired, tx.Status, tx.ID)
	return err
}

// GetExpiringTransactions 만료 예정 거래 내역 조회
func (r *PointRepository) GetExpiringTransactions(ctx context.Context, before time.Time, limit int) ([]*point.Transaction, error) {
	query := `
		SELECT id, user_id, transaction_type, amount, balance_after, reason_type, reason_detail,
		       order_id, earned_at, expires_at, expired, status, created_at
		FROM point_transactions
		WHERE transaction_type = 'EARN'
		  AND expired = false
		  AND status = 'CONFIRMED'
		  AND expires_at IS NOT NULL
		  AND expires_at <= ?
		ORDER BY expires_at ASC
		LIMIT ?
	`

	db := r.tm.GetDBOrTx(ctx)
	rows, err := db.QueryContext(ctx, query, before, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*point.Transaction
	for rows.Next() {
		var tx point.Transaction
		var earnedAt, expiresAt sql.NullTime
		var orderID sql.NullInt64

		err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.Type,
			&tx.Amount,
			&tx.BalanceAfter,
			&tx.ReasonType,
			&tx.ReasonDetail,
			&orderID,
			&earnedAt,
			&expiresAt,
			&tx.Expired,
			&tx.Status,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if orderID.Valid {
			tx.OrderID = &orderID.Int64
		}
		if earnedAt.Valid {
			tx.EarnedAt = &earnedAt.Time
		}
		if expiresAt.Valid {
			tx.ExpiresAt = &expiresAt.Time
		}

		transactions = append(transactions, &tx)
	}

	return transactions, rows.Err()
}

// GetTransactionsByUser 사용자 거래 내역 조회
func (r *PointRepository) GetTransactionsByUser(ctx context.Context, userID int64, limit, offset int) ([]*point.Transaction, error) {
	query := `
		SELECT id, user_id, transaction_type, amount, balance_after, reason_type, reason_detail,
		       order_id, earned_at, expires_at, expired, status, created_at
		FROM point_transactions
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	db := r.tm.GetDBOrTx(ctx)
	rows, err := db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*point.Transaction
	for rows.Next() {
		var tx point.Transaction
		var earnedAt, expiresAt sql.NullTime
		var orderID sql.NullInt64

		err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.Type,
			&tx.Amount,
			&tx.BalanceAfter,
			&tx.ReasonType,
			&tx.ReasonDetail,
			&orderID,
			&earnedAt,
			&expiresAt,
			&tx.Expired,
			&tx.Status,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if orderID.Valid {
			tx.OrderID = &orderID.Int64
		}
		if earnedAt.Valid {
			tx.EarnedAt = &earnedAt.Time
		}
		if expiresAt.Valid {
			tx.ExpiresAt = &expiresAt.Time
		}

		transactions = append(transactions, &tx)
	}

	return transactions, rows.Err()
}

// GetTransactionByID 거래 내역 ID로 조회
func (r *PointRepository) GetTransactionByID(ctx context.Context, id int64) (*point.Transaction, error) {
	query := `
		SELECT id, user_id, transaction_type, amount, balance_after, reason_type, reason_detail,
		       order_id, earned_at, expires_at, expired, status, created_at
		FROM point_transactions
		WHERE id = ?
	`

	db := r.tm.GetDBOrTx(ctx)
	row := db.QueryRowContext(ctx, query, id)

	var tx point.Transaction
	var earnedAt, expiresAt sql.NullTime
	var orderID sql.NullInt64

	err := row.Scan(
		&tx.ID,
		&tx.UserID,
		&tx.Type,
		&tx.Amount,
		&tx.BalanceAfter,
		&tx.ReasonType,
		&tx.ReasonDetail,
		&orderID,
		&earnedAt,
		&expiresAt,
		&tx.Expired,
		&tx.Status,
		&tx.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, point.ErrTransactionNotFound
	}
	if err != nil {
		return nil, err
	}

	if orderID.Valid {
		tx.OrderID = &orderID.Int64
	}
	if earnedAt.Valid {
		tx.EarnedAt = &earnedAt.Time
	}
	if expiresAt.Valid {
		tx.ExpiresAt = &expiresAt.Time
	}

	return &tx, nil
}

// GetTransactionsByOrderID 주문 ID로 거래 내역 조회
func (r *PointRepository) GetTransactionsByOrderID(ctx context.Context, orderID int64) ([]*point.Transaction, error) {
	query := `
		SELECT id, user_id, transaction_type, amount, balance_after, reason_type, reason_detail,
		       order_id, earned_at, expires_at, expired, status, created_at
		FROM point_transactions
		WHERE order_id = ?
		ORDER BY created_at ASC
	`

	db := r.tm.GetDBOrTx(ctx)
	rows, err := db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*point.Transaction
	for rows.Next() {
		var tx point.Transaction
		var earnedAt, expiresAt sql.NullTime
		var orderIDVal sql.NullInt64

		err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.Type,
			&tx.Amount,
			&tx.BalanceAfter,
			&tx.ReasonType,
			&tx.ReasonDetail,
			&orderIDVal,
			&earnedAt,
			&expiresAt,
			&tx.Expired,
			&tx.Status,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if orderIDVal.Valid {
			tx.OrderID = &orderIDVal.Int64
		}
		if earnedAt.Valid {
			tx.EarnedAt = &earnedAt.Time
		}
		if expiresAt.Valid {
			tx.ExpiresAt = &expiresAt.Time
		}

		transactions = append(transactions, &tx)
	}

	return transactions, rows.Err()
}

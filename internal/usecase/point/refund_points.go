package point

import (
	"context"
	"shopping-mall/internal/domain/point"
	"time"
)

// RefundPointsUseCase 포인트 환불 유스케이스
type RefundPointsUseCase struct {
	repo point.Repository
	tm   point.TransactionManager
}

// NewRefundPointsUseCase 포인트 환불 유스케이스 생성
func NewRefundPointsUseCase(repo point.Repository, tm point.TransactionManager) *RefundPointsUseCase {
	return &RefundPointsUseCase{
		repo: repo,
		tm:   tm,
	}
}

// RefundPoints 포인트 환불
func (uc *RefundPointsUseCase) RefundPoints(ctx context.Context, userID int64, orderID int64) error {
	return uc.tm.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. 주문 관련 거래 내역 조회
		transactions, err := uc.repo.GetTransactionsByOrderID(txCtx, orderID)
		if err != nil {
			return err
		}

		// 2. 포인트 잔액 조회
		userPoint, err := uc.repo.GetUserPoint(txCtx, userID)
		if err != nil {
			return err
		}

		// 3. 사용했던 포인트 복구
		var usedAmount int64
		for _, tx := range transactions {
			if tx.Type == point.TransactionTypeUse && tx.Status == point.TransactionStatusConfirmed {
				usedAmount += tx.Amount
			}
		}

		if usedAmount > 0 {
			// 포인트 환불
			userPoint.Refund(usedAmount)

			// 환불 거래 내역 생성
			transaction := &point.Transaction{
				UserID:       userID,
				Type:         point.TransactionTypeEarn,
				Amount:       usedAmount,
				BalanceAfter: userPoint.AvailableBalance,
				ReasonType:   point.ReasonTypeRefund,
				ReasonDetail: "주문 환불",
				OrderID:      &orderID,
				Status:       point.TransactionStatusConfirmed,
				CreatedAt:    time.Now(),
			}

			if err := uc.repo.CreateTransaction(txCtx, transaction); err != nil {
				return err
			}
		}

		// 4. 이미 적립된 포인트 회수
		var earnedAmount int64
		for _, tx := range transactions {
			if tx.Type == point.TransactionTypeEarn && tx.Status == point.TransactionStatusConfirmed {
				earnedAmount += tx.Amount
			}
		}

		if earnedAmount > 0 {
			// 포인트 회수
			userPoint.Expire(earnedAmount)

			// 취소 거래 내역 생성
			transaction := &point.Transaction{
				UserID:       userID,
				Type:         point.TransactionTypeCancel,
				Amount:       earnedAmount,
				BalanceAfter: userPoint.AvailableBalance,
				ReasonType:   point.ReasonTypeRefund,
				ReasonDetail: "주문 환불로 인한 적립 취소",
				OrderID:      &orderID,
				Status:       point.TransactionStatusCancelled,
				CreatedAt:    time.Now(),
			}

			if err := uc.repo.CreateTransaction(txCtx, transaction); err != nil {
				return err
			}

			// 적립 거래 내역 취소 처리
			for _, tx := range transactions {
				if tx.Type == point.TransactionTypeEarn && tx.Status == point.TransactionStatusConfirmed {
					tx.Status = point.TransactionStatusCancelled
					if err := uc.repo.UpdateTransaction(txCtx, tx); err != nil {
						return err
					}
				}
			}
		}

		// 5. 잔액 업데이트
		return uc.repo.UpdateUserPoint(txCtx, userPoint)
	})
}

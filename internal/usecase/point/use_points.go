package point

import (
	"context"
	"shopping-mall/internal/domain/point"
	"time"
)

// UsePointsUseCase 포인트 사용 유스케이스
type UsePointsUseCase struct {
	repo   point.Repository
	tm     point.TransactionManager
	policy *point.Policy
}

// NewUsePointsUseCase 포인트 사용 유스케이스 생성
func NewUsePointsUseCase(repo point.Repository, tm point.TransactionManager, policy *point.Policy) *UsePointsUseCase {
	return &UsePointsUseCase{
		repo:   repo,
		tm:     tm,
		policy: policy,
	}
}

// UsePoints 포인트 사용
func (uc *UsePointsUseCase) UsePoints(ctx context.Context, userID int64, useAmount, orderAmount int64, orderID int64) error {
	return uc.tm.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. 포인트 잔액 조회 (FOR UPDATE 락)
		userPoint, err := uc.repo.GetUserPoint(txCtx, userID)
		if err != nil {
			return err
		}

		// 2. 사용 유효성 검증
		if err := uc.policy.ValidateUse(useAmount, orderAmount, userPoint.AvailableBalance); err != nil {
			return err
		}

		// 3. FIFO 방식으로 적립 내역에서 차감
		earnedTransactions, err := uc.repo.GetEarnedTransactions(txCtx, userID, 100)
		if err != nil {
			return err
		}

		remainingAmount := useAmount
		for _, tx := range earnedTransactions {
			if remainingAmount <= 0 {
				break
			}

			availableAmount := tx.Amount
			if availableAmount > remainingAmount {
				availableAmount = remainingAmount
			}

			remainingAmount -= availableAmount
		}

		if remainingAmount > 0 {
			return point.ErrInsufficientPoints
		}

		// 4. 포인트 차감 (도메인 로직)
		if err := userPoint.Use(useAmount); err != nil {
			return err
		}

		// 5. 사용 거래 내역 생성
		transaction := &point.Transaction{
			UserID:       userID,
			Type:         point.TransactionTypeUse,
			Amount:       useAmount,
			BalanceAfter: userPoint.AvailableBalance,
			ReasonType:   point.ReasonTypePurchase,
			ReasonDetail: "주문 결제",
			OrderID:      &orderID,
			Status:       point.TransactionStatusConfirmed,
			CreatedAt:    time.Now(),
		}

		if err := uc.repo.CreateTransaction(txCtx, transaction); err != nil {
			return err
		}

		// 6. 잔액 업데이트
		return uc.repo.UpdateUserPoint(txCtx, userPoint)
	})
}

package point

import (
	"context"
	"shopping-mall/internal/domain/point"
	"time"
)

// ExpirePointsUseCase 포인트 만료 유스케이스
type ExpirePointsUseCase struct {
	repo point.Repository
	tm   point.TransactionManager
}

// NewExpirePointsUseCase 포인트 만료 유스케이스 생성
func NewExpirePointsUseCase(repo point.Repository, tm point.TransactionManager) *ExpirePointsUseCase {
	return &ExpirePointsUseCase{
		repo: repo,
		tm:   tm,
	}
}

// ExpirePoints 만료 포인트 처리
func (uc *ExpirePointsUseCase) ExpirePoints(ctx context.Context, before time.Time, limit int) error {
	return uc.tm.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. 만료 대상 조회
		transactions, err := uc.repo.GetExpiringTransactions(txCtx, before, limit)
		if err != nil {
			return err
		}

		// 2. 사용자별로 그룹화
		userTransactions := make(map[int64][]*point.Transaction)
		for _, tx := range transactions {
			userTransactions[tx.UserID] = append(userTransactions[tx.UserID], tx)
		}

		// 3. 사용자별로 만료 처리
		for userID, txs := range userTransactions {
			// 포인트 잔액 조회
			userPoint, err := uc.repo.GetUserPoint(txCtx, userID)
			if err != nil {
				continue
			}

			// 만료 포인트 계산
			var expireAmount int64
			for _, tx := range txs {
				if !tx.Expired {
					expireAmount += tx.Amount
					tx.Expired = true
					if err := uc.repo.UpdateTransaction(txCtx, tx); err != nil {
						return err
					}
				}
			}

			if expireAmount > 0 {
				// 포인트 만료
				userPoint.Expire(expireAmount)

				// 만료 거래 내역 생성
				transaction := &point.Transaction{
					UserID:       userID,
					Type:         point.TransactionTypeExpire,
					Amount:       expireAmount,
					BalanceAfter: userPoint.AvailableBalance,
					ReasonType:   point.ReasonTypeAdmin,
					ReasonDetail: "포인트 만료",
					Status:       point.TransactionStatusConfirmed,
					CreatedAt:    time.Now(),
				}

				if err := uc.repo.CreateTransaction(txCtx, transaction); err != nil {
					return err
				}

				// 잔액 업데이트
				if err := uc.repo.UpdateUserPoint(txCtx, userPoint); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

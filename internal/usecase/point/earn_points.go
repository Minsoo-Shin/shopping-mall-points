package point

import (
	"context"
	"shopping-mall/internal/domain/point"
	"time"
)

// EarnPointsUseCase 포인트 적립 유스케이스
type EarnPointsUseCase struct {
	repo   point.Repository
	tm     point.TransactionManager
	policy *point.Policy
}

// NewEarnPointsUseCase 포인트 적립 유스케이스 생성
func NewEarnPointsUseCase(repo point.Repository, tm point.TransactionManager, policy *point.Policy) *EarnPointsUseCase {
	return &EarnPointsUseCase{
		repo:   repo,
		tm:     tm,
		policy: policy,
	}
}

// EarnPointsFromPurchase 구매 적립
func (uc *EarnPointsUseCase) EarnPointsFromPurchase(ctx context.Context, userID int64, paymentAmount int64, orderID int64) error {
	return uc.tm.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. 포인트 잔액 조회
		userPoint, err := uc.repo.GetUserPoint(txCtx, userID)
		if err != nil {
			// 없으면 생성
			userPoint = &point.UserPoint{
				UserID:           userID,
				AvailableBalance: 0,
				PendingBalance:   0,
				TotalEarned:      0,
				TotalUsed:        0,
				UpdatedAt:        time.Now(),
			}
		}

		// 2. 적립 포인트 계산
		earnAmount := uc.policy.CalculateEarnPoints(paymentAmount)
		if earnAmount <= 0 {
			return nil
		}

		// 3. 포인트 적립 (도메인 로직)
		userPoint.Earn(earnAmount)

		// 4. 적립 거래 내역 생성
		now := time.Now()
		expiresAt := uc.policy.CalculateExpiryDate(now)

		transaction := &point.Transaction{
			UserID:       userID,
			Type:         point.TransactionTypeEarn,
			Amount:       earnAmount,
			BalanceAfter: userPoint.AvailableBalance,
			ReasonType:   point.ReasonTypePurchase,
			ReasonDetail: "구매 적립",
			OrderID:      &orderID,
			EarnedAt:     &now,
			ExpiresAt:    &expiresAt,
			Expired:      false,
			Status:       point.TransactionStatusConfirmed,
			CreatedAt:    now,
		}

		if err := uc.repo.CreateTransaction(txCtx, transaction); err != nil {
			return err
		}

		// 5. 잔액 업데이트
		return uc.repo.UpdateUserPoint(txCtx, userPoint)
	})
}

// EarnPointsFromReview 리뷰 적립
func (uc *EarnPointsUseCase) EarnPointsFromReview(ctx context.Context, userID int64, isPhoto bool) error {
	return uc.tm.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. 포인트 잔액 조회
		userPoint, err := uc.repo.GetUserPoint(txCtx, userID)
		if err != nil {
			userPoint = &point.UserPoint{
				UserID:           userID,
				AvailableBalance: 0,
				PendingBalance:   0,
				TotalEarned:      0,
				TotalUsed:        0,
				UpdatedAt:        time.Now(),
			}
		}

		// 2. 적립 포인트 계산
		var earnAmount int64
		var reasonDetail string
		if isPhoto {
			earnAmount = uc.policy.ReviewPhotoPoints
			reasonDetail = "포토 리뷰 적립"
		} else {
			earnAmount = uc.policy.ReviewTextPoints
			reasonDetail = "텍스트 리뷰 적립"
		}

		// 3. 포인트 적립
		userPoint.Earn(earnAmount)

		// 4. 적립 거래 내역 생성
		now := time.Now()
		expiresAt := uc.policy.CalculateExpiryDate(now)

		transaction := &point.Transaction{
			UserID:       userID,
			Type:         point.TransactionTypeEarn,
			Amount:       earnAmount,
			BalanceAfter: userPoint.AvailableBalance,
			ReasonType:   point.ReasonTypeReview,
			ReasonDetail: reasonDetail,
			EarnedAt:     &now,
			ExpiresAt:    &expiresAt,
			Expired:      false,
			Status:       point.TransactionStatusConfirmed,
			CreatedAt:    now,
		}

		if err := uc.repo.CreateTransaction(txCtx, transaction); err != nil {
			return err
		}

		// 5. 잔액 업데이트
		return uc.repo.UpdateUserPoint(txCtx, userPoint)
	})
}

// EarnSignupBonus 가입 보너스 적립
func (uc *EarnPointsUseCase) EarnSignupBonus(ctx context.Context, userID int64) error {
	return uc.tm.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. 포인트 잔액 조회
		userPoint, err := uc.repo.GetUserPoint(txCtx, userID)
		if err != nil {
			userPoint = &point.UserPoint{
				UserID:           userID,
				AvailableBalance: 0,
				PendingBalance:   0,
				TotalEarned:      0,
				TotalUsed:        0,
				UpdatedAt:        time.Now(),
			}
		}

		// 2. 가입 보너스 적립
		earnAmount := uc.policy.SignupBonus
		userPoint.Earn(earnAmount)

		// 3. 적립 거래 내역 생성
		now := time.Now()
		expiresAt := uc.policy.CalculateExpiryDate(now)

		transaction := &point.Transaction{
			UserID:       userID,
			Type:         point.TransactionTypeEarn,
			Amount:       earnAmount,
			BalanceAfter: userPoint.AvailableBalance,
			ReasonType:   point.ReasonTypeSignup,
			ReasonDetail: "가입 보너스",
			EarnedAt:     &now,
			ExpiresAt:    &expiresAt,
			Expired:      false,
			Status:       point.TransactionStatusConfirmed,
			CreatedAt:    now,
		}

		if err := uc.repo.CreateTransaction(txCtx, transaction); err != nil {
			return err
		}

		// 4. 잔액 업데이트
		return uc.repo.UpdateUserPoint(txCtx, userPoint)
	})
}

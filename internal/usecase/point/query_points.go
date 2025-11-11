package point

import (
	"context"
	"shopping-mall/internal/domain/point"
	"shopping-mall/internal/repository/redis"
)

// QueryPointsUseCase 포인트 조회 유스케이스
type QueryPointsUseCase struct {
	repo  point.Repository
	cache *redis.PointCache
}

// NewQueryPointsUseCase 포인트 조회 유스케이스 생성
func NewQueryPointsUseCase(repo point.Repository, cache *redis.PointCache) *QueryPointsUseCase {
	return &QueryPointsUseCase{
		repo:  repo,
		cache: cache,
	}
}

// GetBalance 잔액 조회
func (uc *QueryPointsUseCase) GetBalance(ctx context.Context, userID int64) (*point.UserPoint, error) {
	// 캐시에서 조회 시도
	if uc.cache != nil {
		cached, err := uc.cache.GetBalance(ctx, userID)
		if err == nil && cached != nil {
			return &point.UserPoint{
				UserID:           userID,
				AvailableBalance: cached.AvailableBalance,
				PendingBalance:   cached.PendingBalance,
				TotalEarned:      cached.TotalEarned,
				TotalUsed:        cached.TotalUsed,
			}, nil
		}
	}

	// DB에서 조회
	userPoint, err := uc.repo.GetUserPoint(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 캐시에 저장
	if uc.cache != nil {
		_ = uc.cache.SetBalance(ctx, userID, &redis.BalanceCache{
			AvailableBalance: userPoint.AvailableBalance,
			PendingBalance:   userPoint.PendingBalance,
			TotalEarned:      userPoint.TotalEarned,
			TotalUsed:        userPoint.TotalUsed,
		})
	}

	return userPoint, nil
}

// GetTransactions 거래 내역 조회
func (uc *QueryPointsUseCase) GetTransactions(ctx context.Context, userID int64, limit, offset int) ([]*point.Transaction, error) {
	return uc.repo.GetTransactionsByUser(ctx, userID, limit, offset)
}

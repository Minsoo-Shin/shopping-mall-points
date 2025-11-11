package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// PointCache 포인트 캐시
type PointCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewPointCache 포인트 캐시 생성
func NewPointCache(client *redis.Client) *PointCache {
	return &PointCache{
		client: client,
		ttl:    5 * time.Minute, // 기본 TTL 5분
	}
}

// CacheKey 캐시 키 생성
func (c *PointCache) CacheKey(userID int64) string {
	return fmt.Sprintf("point:balance:%d", userID)
}

// GetBalance 잔액 조회
func (c *PointCache) GetBalance(ctx context.Context, userID int64) (*BalanceCache, error) {
	key := c.CacheKey(userID)
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var balance BalanceCache
	if err := json.Unmarshal([]byte(val), &balance); err != nil {
		return nil, err
	}

	return &balance, nil
}

// SetBalance 잔액 캐싱
func (c *PointCache) SetBalance(ctx context.Context, userID int64, balance *BalanceCache) error {
	key := c.CacheKey(userID)
	val, err := json.Marshal(balance)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, val, c.ttl).Err()
}

// DeleteBalance 잔액 캐시 삭제
func (c *PointCache) DeleteBalance(ctx context.Context, userID int64) error {
	key := c.CacheKey(userID)
	return c.client.Del(ctx, key).Err()
}

// BalanceCache 잔액 캐시 구조체
type BalanceCache struct {
	AvailableBalance int64 `json:"available_balance"`
	PendingBalance   int64 `json:"pending_balance"`
	TotalEarned      int64 `json:"total_earned"`
	TotalUsed        int64 `json:"total_used"`
}

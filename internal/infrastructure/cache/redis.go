package cache

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Config Redis 설정
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// NewRedis Redis 연결 생성
func NewRedis(cfg Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return client, nil
}

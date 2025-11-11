package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"shopping-mall/config"
	"shopping-mall/internal/domain/point"
	"shopping-mall/internal/infrastructure/cache"
	"shopping-mall/internal/infrastructure/database"
	"shopping-mall/internal/infrastructure/logger"
	"shopping-mall/internal/repository/mysql"
	"shopping-mall/internal/usecase/point"
)

func main() {
	// 설정 로드
	cfg := config.Load()
	
	// 로거 초기화
	zapLogger, err := logger.NewLogger(cfg.Server.Env)
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer zapLogger.Sync()
	
	// MySQL 연결
	db, err := database.NewMySQL(database.Config{
		Host:     cfg.MySQL.Host,
		Port:     cfg.MySQL.Port,
		User:     cfg.MySQL.User,
		Password: cfg.MySQL.Password,
		Database: cfg.MySQL.Database,
	})
	if err != nil {
		zapLogger.Fatal("Failed to connect to MySQL", zap.Error(err))
	}
	defer db.Close()
	
	// Repository 초기화
	tm := mysql.NewTransactionManager(db)
	pointRepo := mysql.NewPointRepository(tm)
	
	// UseCase 초기화
	expireUseCase := point.NewExpirePointsUseCase(pointRepo, tm)
	
	zapLogger.Info("Point expiration worker started")
	
	// 매일 자정에 실행되는 틱커
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	
	// 즉시 한 번 실행
	runExpiration(zapLogger, expireUseCase)
	
	// 시그널 대기 및 주기적 실행
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	for {
		select {
		case <-ticker.C:
			runExpiration(zapLogger, expireUseCase)
		case <-quit:
			zapLogger.Info("Worker shutting down...")
			return
		}
	}
}

func runExpiration(logger *zap.Logger, expireUseCase *point.ExpirePointsUseCase) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	
	now := time.Now()
	limit := 1000 // 한 번에 처리할 최대 개수
	
	logger.Info("Running point expiration", zap.Time("before", now))
	
	if err := expireUseCase.ExpirePoints(ctx, now, limit); err != nil {
		logger.Error("Failed to expire points", zap.Error(err))
		return
	}
	
	logger.Info("Point expiration completed")
}


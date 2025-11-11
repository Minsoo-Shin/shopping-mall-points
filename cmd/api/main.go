package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"shopping-mall/config"
	"shopping-mall/internal/domain/point"
	"shopping-mall/internal/handler/http"
	"shopping-mall/internal/infrastructure/cache"
	"shopping-mall/internal/infrastructure/database"
	"shopping-mall/internal/infrastructure/logger"
	"shopping-mall/internal/repository/mysql"
	"shopping-mall/internal/repository/redis"
	pointUseCase "shopping-mall/internal/usecase/point"
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
	
	// Redis 연결
	redisClient, err := cache.NewRedis(cache.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		zapLogger.Warn("Failed to connect to Redis, continuing without cache", zap.Error(err))
		redisClient = nil
	}
	
	// Repository 초기화
	tm := mysql.NewTransactionManager(db)
	pointRepo := mysql.NewPointRepository(tm)
	var pointCache *redis.PointCache
	if redisClient != nil {
		pointCache = redis.NewPointCache(redisClient)
	}
	
	// Policy 초기화
	policy := point.NewDefaultPolicy()
	
	// UseCase 초기화
	queryUseCase := pointUseCase.NewQueryPointsUseCase(pointRepo, pointCache)
	useUseCase := pointUseCase.NewUsePointsUseCase(pointRepo, tm, policy)
	earnUseCase := pointUseCase.NewEarnPointsUseCase(pointRepo, tm, policy)
	refundUseCase := pointUseCase.NewRefundPointsUseCase(pointRepo, tm)
	
	// Handler 초기화
	pointHandler := http.NewPointHandler(queryUseCase, useUseCase, earnUseCase)
	orderHandler := http.NewOrderHandler(useUseCase, earnUseCase, refundUseCase)
	
	// Router 설정
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()
	
	// 포인트 관련 엔드포인트
	api.HandleFunc("/points/balance", pointHandler.GetBalance).Methods("GET")
	api.HandleFunc("/points/transactions", pointHandler.GetTransactions).Methods("GET")
	api.HandleFunc("/points/use", pointHandler.UsePoints).Methods("POST")
	api.HandleFunc("/points/earn", pointHandler.EarnPoints).Methods("POST")
	
	// 주문 관련 엔드포인트
	api.HandleFunc("/orders/{id}/confirm", orderHandler.ConfirmOrder).Methods("POST")
	api.HandleFunc("/orders/{id}/refund", orderHandler.RefundOrder).Methods("POST")
	
	// 서버 시작
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// Graceful shutdown
	go func() {
		zapLogger.Info("Server starting", zap.String("port", cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapLogger.Fatal("Server failed to start", zap.Error(err))
		}
	}()
	
	// 시그널 대기
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	zapLogger.Info("Server shutting down...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		zapLogger.Fatal("Server forced to shutdown", zap.Error(err))
	}
	
	zapLogger.Info("Server exited")
}


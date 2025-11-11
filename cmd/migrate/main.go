package main

import (
	"log"
	"os"

	"shopping-mall/config"
	"shopping-mall/internal/infrastructure/database"
)

func main() {
	// 설정 로드
	cfg := config.Load()

	// 마이그레이션 디렉토리
	migrationsDir := "migrations"
	if len(os.Args) > 1 {
		migrationsDir = os.Args[1]
	}

	log.Printf("Initializing database: %s", cfg.MySQL.Database)
	log.Printf("MySQL Host: %s:%d", cfg.MySQL.Host, cfg.MySQL.Port)
	log.Printf("MySQL User: %s", cfg.MySQL.User)
	
	// 비밀번호 확인
	if cfg.MySQL.Password == "" {
		log.Println("⚠ Warning: MYSQL_PASSWORD is not set. Using empty password.")
		log.Println("   Set MYSQL_PASSWORD environment variable if your MySQL requires a password.")
	}
	
	log.Printf("Migrations directory: %s", migrationsDir)

	// 데이터베이스 초기화
	if err := database.InitDatabase(database.Config{
		Host:     cfg.MySQL.Host,
		Port:     cfg.MySQL.Port,
		User:     cfg.MySQL.User,
		Password: cfg.MySQL.Password,
		Database: cfg.MySQL.Database,
	}, migrationsDir); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("✓ Database initialization completed successfully")
}


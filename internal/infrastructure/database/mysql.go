package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Config MySQL 설정
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// NewMySQL MySQL 연결 생성
func NewMySQL(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// 연결 풀 설정
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0) // 연결 재사용 시간 제한 없음

	return db, nil
}

// NewMySQLWithInit MySQL 연결 생성 및 초기화
func NewMySQLWithInit(cfg Config, migrationsDir string) (*sql.DB, error) {
	// 데이터베이스 초기화
	if err := InitDatabase(cfg, migrationsDir); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// 연결 생성
	return NewMySQL(cfg)
}

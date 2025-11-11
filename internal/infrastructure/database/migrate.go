package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Migrate 마이그레이션 실행
func Migrate(db *sql.DB, migrationsDir string) error {
	// 마이그레이션 파일 읽기
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// 파일명 순서대로 정렬하여 실행
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		filePath := filepath.Join(migrationsDir, file.Name())
		sqlBytes, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		// SQL 파일 전체를 하나의 트랜잭션으로 실행
		sqlContent := string(sqlBytes)
		
		// 주석 제거 (간단한 처리)
		lines := strings.Split(sqlContent, "\n")
		var cleanLines []string
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			// 주석 라인과 빈 라인 제외
			if trimmed != "" && !strings.HasPrefix(trimmed, "--") {
				cleanLines = append(cleanLines, line)
			}
		}
		sqlContent = strings.Join(cleanLines, "\n")

		// 세미콜론으로 분리하되, 빈 문장은 제외
		statements := strings.Split(sqlContent, ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			// SQL 실행
			if _, err := db.Exec(stmt); err != nil {
				// 이미 존재하는 테이블은 무시 (IF NOT EXISTS)
				if strings.Contains(err.Error(), "already exists") || 
				   strings.Contains(err.Error(), "Duplicate") {
					fmt.Printf("⚠ Table already exists in %s, skipping\n", file.Name())
					continue
				}
				return fmt.Errorf("failed to execute migration %s: %w\nSQL: %s", file.Name(), err, stmt)
			}
		}

		fmt.Printf("✓ Migration %s executed successfully\n", file.Name())
	}

	return nil
}

// EnsureDatabase 데이터베이스가 없으면 생성
func EnsureDatabase(cfg Config) error {
	// 데이터베이스 이름을 제외한 DSN 생성
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?parseTime=true&charset=utf8mb4",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	defer db.Close()

	// 데이터베이스 생성 (없으면)
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.Database)
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	return nil
}

// InitDatabase 데이터베이스 초기화 (데이터베이스 생성 + 마이그레이션)
func InitDatabase(cfg Config, migrationsDir string) error {
	// 1. 데이터베이스 생성
	if err := EnsureDatabase(cfg); err != nil {
		return fmt.Errorf("failed to ensure database: %w", err)
	}

	// 2. 데이터베이스 연결
	db, err := NewMySQL(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// 3. 마이그레이션 실행
	if migrationsDir != "" {
		if _, err := os.Stat(migrationsDir); err == nil {
			if err := Migrate(db, migrationsDir); err != nil {
				return fmt.Errorf("failed to run migrations: %w", err)
			}
		}
	}

	return nil
}


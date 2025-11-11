package config

import (
	"os"
	"strconv"
)

// Config 애플리케이션 설정
type Config struct {
	Server ServerConfig
	MySQL  MySQLConfig
	Redis  RedisConfig
}

// ServerConfig 서버 설정
type ServerConfig struct {
	Port string
	Env  string
}

// MySQLConfig MySQL 설정
type MySQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// RedisConfig Redis 설정
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Load 설정 로드
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		MySQL: MySQLConfig{
			Host:     getEnv("MYSQL_HOST", "localhost"),
			Port:     getEnvAsInt("MYSQL_PORT", 3306),
			User:     getEnv("MYSQL_USER", "root"),
			Password: getEnv("MYSQL_PASSWORD", ""),
			Database: getEnv("MYSQL_DATABASE", "shopping_mall"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}


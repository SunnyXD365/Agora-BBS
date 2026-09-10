package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	DBDSN          string
	JWTSecret      string
	JWTExpireHours int
}

// 从环境变量加载配置，提供默认备选项
func LoadConfig() *Config {
	port := getEnv("PORT", "8080")
	dbDSN := getEnv("DATABASE_URL", "postgres://agora_user:agora_password@agora-postgres:5432/agora_db?sslmode=disable")
	jwtSecret := getEnv("JWT_SECRET", "agora_dev_secret_key_change_in_prod")
	expireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "72"))

	return &Config{
		Port:           port,
		DBDSN:          dbDSN,
		JWTSecret:      jwtSecret,
		JWTExpireHours: expireHours,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

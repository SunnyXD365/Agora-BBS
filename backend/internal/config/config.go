package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port                        string
	DBDSN                       string
	JWTSecret                   string
	JWTExpireHours              int
	AppEnv                      string
	RedisAddr                   string
	TemporalHost                string
	TemporalTaskQueue           string
	GovernanceCoolingSeconds    int
	GovernanceReplyDwellSeconds int
	GovernanceLongTopicChars    int
}

// 从环境变量加载配置，提供默认备选项
func LoadConfig() *Config {
	port := getEnv("PORT", "8080")
	dbDSN := getEnv("DATABASE_URL", "postgres://agora_user:agora_password@agora-postgres:5432/agora_db?sslmode=disable")
	jwtSecret := getEnv("JWT_SECRET", "agora_dev_secret_key_change_in_prod")
	expireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "72"))
	coolingSeconds := getEnvInt("GOVERNANCE_COOLING_SECONDS", 60)
	replyDwellSeconds := getEnvInt("GOVERNANCE_REPLY_DWELL_SECONDS", 10)
	longTopicChars := getEnvInt("GOVERNANCE_LONG_TOPIC_CHARS", 800)

	return &Config{
		Port:                        port,
		DBDSN:                       dbDSN,
		JWTSecret:                   jwtSecret,
		JWTExpireHours:              expireHours,
		AppEnv:                      getEnv("APP_ENV", "development"),
		RedisAddr:                   getEnv("REDIS_ADDR", "agora-redis:6379"),
		TemporalHost:                getEnv("TEMPORAL_HOST", "agora-temporal:7233"),
		TemporalTaskQueue:           getEnv("TEMPORAL_TASK_QUEUE", "agora-governance"),
		GovernanceCoolingSeconds:    coolingSeconds,
		GovernanceReplyDwellSeconds: replyDwellSeconds,
		GovernanceLongTopicChars:    longTopicChars,
	}
}

func getEnvInt(key string, fallback int) int {
	value, err := strconv.Atoi(getEnv(key, strconv.Itoa(fallback)))
	if err != nil || value < 0 {
		return fallback
	}
	return value
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

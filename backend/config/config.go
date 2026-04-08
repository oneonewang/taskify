package config

import (
	"os"
)

// Config 应用配置
type Config struct {
	Port        string
	DatabaseURL string
	GinMode     string
}

// Load 加载环境变量配置
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		GinMode:     getEnv("GIN_MODE", "debug"),
	}
}

// getEnv 获取环境变量，默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

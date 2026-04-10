package config

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// SessionConfig 会话配置
type SessionConfig struct {
	Name     string
	Secret   string
	MaxAge   int
	Secure   bool
	HTTPOnly bool
	SameSite http.SameSite
}

// DefaultSessionConfig 默认会话配置
func DefaultSessionConfig() SessionConfig {
	return SessionConfig{
		Name:     "session_id",
		Secret:   getEnv("SESSION_SECRET", "taskify-secret-key-change-in-production"),
		MaxAge:   86400, // 24小时
		Secure:   false, // 开发环境为false
		HTTPOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

// InitSession 初始化会话中间件
func InitSession(r *gin.Engine, cfg SessionConfig) {
	store := cookie.NewStore([]byte(cfg.Secret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   cfg.MaxAge,
		Secure:   cfg.Secure,
		HttpOnly: cfg.HTTPOnly,
		SameSite: cfg.SameSite,
	})
	r.Use(sessions.Sessions(cfg.Name, store))
}

// SessionSecret 获取会话密钥
func SessionSecret() string {
	return getEnv("SESSION_SECRET", "taskify-secret-key-change-in-production")
}

// SessionMaxAge 获取会话最大年龄（秒）
func SessionMaxAge() int {
	return 86400 // 24小时
}

// IsProduction 是否生产环境
func IsProduction() bool {
	return getEnv("GIN_MODE", "debug") == "release"
}
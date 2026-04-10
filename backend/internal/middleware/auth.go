package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/pkg/response"
)

// UserIDKey 用户ID在context中的key
const UserIDKey = "user_id"

// AuthRequired 认证中间件 - 验证用户是否已登录
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")

		if userID == nil {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// 将user_id存入gin.Context供后续处理器使用
		c.Set(UserIDKey, userID)
		c.Next()
	}
}

// GetUserID 从context中获取当前用户ID
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}
	switch v := userID.(type) {
	case uint:
		return v, true
	case int:
		return uint(v), true
	case int64:
		return uint(v), true
	default:
		return 0, false
	}
}

// GetClientIP 获取客户端IP地址
func GetClientIP(c *gin.Context) string {
	// 优先从X-Forwarded-For获取
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		return xff
	}
	// 其次从X-Real-IP获取
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}
	return c.ClientIP()
}

// ErrorHandler 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
	}
}
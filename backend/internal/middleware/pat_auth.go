package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/internal/services"
)

const (
	// ContextKeyUserID 用户ID上下文键
	ContextKeyUserID = "user_id"
	// ContextKeyUsername 用户名上下文键
	ContextKeyUsername = "username"
	// ContextKeyTokenScope 令牌作用域上下文键
	ContextKeyTokenScope = "token_scope"
)

// extractBearerToken 从请求头提取 Bearer Token
func ExtractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

// PATAuthMiddleware PAT Bearer 认证中间件
// 验证 PAT 并注入用户信息到 Context
func PATAuthMiddleware() gin.HandlerFunc {
	tokenService := services.NewTokenService()

	return func(c *gin.Context) {
		// 从请求中提取 Bearer Token
		token := ExtractBearerToken(c.Request)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "token_invalid",
				"message": "Authorization header required",
			})
			return
		}

		// 验证令牌
		tokenInfo, err := tokenService.GetTokenInfo(token)
		if err != nil || !tokenInfo.Active {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "token_invalid",
				"message": "Invalid or expired token",
			})
			return
		}

		// 注入用户信息到 Context
		c.Set(ContextKeyUserID, tokenInfo.UserID)
		c.Set(ContextKeyUsername, tokenInfo.Username)
		c.Set(ContextKeyTokenScope, tokenInfo.Scope)

		c.Next()
	}
}

// GetUserIDFromContext 从 Context 获取用户ID
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(ContextKeyUserID)
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}

// GetUsernameFromContext 从 Context 获取用户名
func GetUsernameFromContext(c *gin.Context) (string, bool) {
	username, exists := c.Get(ContextKeyUsername)
	if !exists {
		return "", false
	}
	name, ok := username.(string)
	return name, ok
}

// GetTokenScopeFromContext 从 Context 获取令牌作用域
func GetTokenScopeFromContext(c *gin.Context) (string, bool) {
	scope, exists := c.Get(ContextKeyTokenScope)
	if !exists {
		return "", false
	}
	s, ok := scope.(string)
	return s, ok
}

// RequireScope 检查令牌是否有指定权限
func RequireScope(requiredScope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		scope, exists := GetTokenScopeFromContext(c)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "scope_insufficient",
				"message": "Token scope not found",
			})
			return
		}

		if !services.CheckScope(scope, requiredScope) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "scope_insufficient",
				"message": "Insufficient token scope",
			})
			return
		}

		c.Next()
	}
}
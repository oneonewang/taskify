package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/internal/middleware"
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/services"
)

// TokenHandler 令牌处理
type TokenHandler struct {
	tokenService *services.TokenService
}

// NewTokenHandler 创建令牌处理器
func NewTokenHandler() *TokenHandler {
	return &TokenHandler{
		tokenService: services.NewTokenService(),
	}
}

// GetTokenInfo GET /oauth/token/info - Token Introspection（公开端点）
// MCP 客户端使用此端点验证 PAT
func (h *TokenHandler) GetTokenInfo(c *gin.Context) {
	// 从 Authorization Header 提取 Bearer Token
	token := middleware.ExtractBearerToken(c.Request)
	if token == "" {
		c.JSON(http.StatusOK, models.TokenInfoResponse{
			Active: false,
		})
		return
	}

	// 获取令牌信息
	info, err := h.tokenService.GetTokenInfo(token)
	if err != nil || !info.Active {
		c.JSON(http.StatusOK, models.TokenInfoResponse{
			Active: false,
		})
		return
	}

	c.JSON(http.StatusOK, info)
}

// CreateToken POST /oauth/tokens - 创建新令牌
func (h *TokenHandler) CreateToken(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not found in context",
		})
		return
	}

	var req models.CreateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid request body: " + err.Error(),
		})
		return
	}

	// 验证作用域格式
	if !services.ValidateScopeFormat(req.Scope) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid scope format. Expected: resource:action,resource:action",
		})
		return
	}

	// 默认过期天数
	if req.ExpiresInDays <= 0 {
		req.ExpiresInDays = 90
	}

	// 创建令牌
	resp, rawToken, err := h.tokenService.CreateTokenForUser(userID, req.Name, req.Scope, req.ExpiresInDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to create token: " + err.Error(),
		})
		return
	}

	// 返回带明文令牌的响应
	c.JSON(http.StatusCreated, gin.H{
		"id":           resp.ID,
		"token":        rawToken,
		"name":         resp.Name,
		"scope":        resp.Scope,
		"token_prefix": resp.TokenPrefix,
		"created_at":   resp.CreatedAt,
		"expires_at":   resp.ExpiresAt,
	})
}

// ListTokens GET /oauth/tokens - 列出用户令牌
func (h *TokenHandler) ListTokens(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not found in context",
		})
		return
	}

	resp, err := h.tokenService.ListUserTokens(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to list tokens: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// RevokeToken DELETE /oauth/tokens/:id - 撤销令牌
func (h *TokenHandler) RevokeToken(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unauthorized",
			"message": "User not found in context",
		})
		return
	}

	tokenIDStr := c.Param("id")
	tokenID, err := strconv.ParseUint(tokenIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Invalid token ID",
		})
		return
	}

	err = h.tokenService.RevokeToken(uint(tokenID), userID)
	if err != nil {
		if err.Error() == "token not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "token_not_found",
			})
			return
		}
		if err.Error() == "unauthorized" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "Cannot revoke token owned by another user",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to revoke token: " + err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
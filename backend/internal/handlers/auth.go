package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/config"
	"github.com/taskify/backend/internal/middleware"
	"github.com/taskify/backend/internal/repository"
	"github.com/taskify/backend/internal/services"
	"github.com/taskify/backend/pkg/response"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
	}
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	DisplayName string `json:"display_name" binding:"required,min=1,max=100"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest 更新资料请求结构
type UpdateProfileRequest struct {
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

// ChangePasswordRequest 修改密码请求结构
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// Register 注册新用户
// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	authResp, err := h.authService.Register(&services.RegisterRequest{
		Email:       req.Email,
		Password:    req.Password,
		DisplayName: req.DisplayName,
	})

	if err != nil {
		if err.Error() == "邮箱已被注册" {
			response.Conflict(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, "注册成功", authResp)
}

// Login 用户登录
// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	ipAddress := middleware.GetClientIP(c)
	authResp, err := h.authService.Login(&services.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}, ipAddress)

	if err != nil {
		if err.Error() == "邮箱或密码错误" {
			response.Unauthorized(c, err.Error())
			return
		}
		response.Unauthorized(c, err.Error())
		return
	}

	// 创建会话
	session := sessions.Default(c)
	session.Set("user_id", authResp.User.ID)
	session.Options(sessions.Options{
		Path:     "/",
		MaxAge:   config.SessionMaxAge(),
		HttpOnly: true,
		Secure:   config.IsProduction(),
		SameSite: http.SameSiteLaxMode,
	})
	session.Save()

	isAdmin, _ := repository.NewMembershipRepository().IsAdmin(authResp.User.ID)
	authResp.User.IsAdmin = isAdmin

	response.Success(c, authResp)
}

// Logout 用户登出
// POST /api/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if ok {
		ipAddress := middleware.GetClientIP(c)
		h.authService.Logout(userID, ipAddress)
	}

	// 清除会话
	session := sessions.Default(c)
	session.Options(sessions.Options{
		Path:   "/",
		MaxAge: 0,
	})
	session.Clear()
	session.Save()

	response.OK(c, "登出成功")
}

// GetCurrentUser 获取当前用户
// GET /api/users/me
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	userResp, err := h.authService.GetUserByID(userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	isAdmin, _ := repository.NewMembershipRepository().IsAdmin(userID)
	userResp.IsAdmin = isAdmin

	response.Success(c, userResp)
}

// UpdateProfile 更新当前用户资料
// PUT /api/users/me
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	userResp, err := h.authService.UpdateProfile(userID, req.DisplayName, req.AvatarURL)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, userResp)
}

// ChangePassword 修改密码
// PUT /api/users/me/password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	err := h.authService.ChangePassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		if err.Error() == "当前密码错误" {
			response.Unauthorized(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "密码修改成功")
}
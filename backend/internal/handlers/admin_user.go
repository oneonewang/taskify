package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/internal/middleware"
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
	"github.com/taskify/backend/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

const DefaultResetPassword = "admin123" // 默认重置密码

// AdminUserHandler 管理员用户处理器
type AdminUserHandler struct {
	userRepo *repository.UserRepository
}

// NewAdminUserHandler 创建管理员用户处理器
func NewAdminUserHandler() *AdminUserHandler {
	return &AdminUserHandler{
		userRepo: repository.NewUserRepository(),
	}
}

// UserQueryRequest 用户查询请求
type UserQueryRequest struct {
	Email      string `form:"email"`       // 邮箱精确匹配
	DisplayName string `form:"display_name"` // 显示名称模糊匹配
	IsDisabled string `form:"is_disabled"`  // 禁用状态筛选: "true", "false", ""(不过滤)
	Page       int    `form:"page,default=1"`      // 页码
	PageSize   int    `form:"page_size,default=20"` // 每页数量
}

// ListUsers 获取用户列表（支持筛选和分页）
// GET /api/admin/users
func (h *AdminUserHandler) ListUsers(c *gin.Context) {
	var req UserQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	filter := repository.UserFilter{
		Email:       req.Email,
		DisplayName: req.DisplayName,
		Page:        req.Page,
		PageSize:    req.PageSize,
	}

	// 解析 is_disabled 参数
	if req.IsDisabled == "true" {
		disabled := true
		filter.IsDisabled = &disabled
	} else if req.IsDisabled == "false" {
		disabled := false
		filter.IsDisabled = &disabled
	}

	result, err := h.userRepo.FindUsers(filter)
	if err != nil {
		response.InternalError(c, "查询用户列表失败")
		return
	}

	// 转换为响应格式
	users := make([]models.UserResponse, len(result.Users))
	for i, user := range result.Users {
		users[i] = user.ToResponse()
	}

	response.Success(c, gin.H{
		"users":        users,
		"total":        result.Total,
		"page":         result.Page,
		"page_size":    result.PageSize,
		"total_pages":  result.TotalPages,
	})
}

// ResetPassword 重置用户密码
// POST /api/admin/users/:id/reset-password
func (h *AdminUserHandler) ResetPassword(c *gin.Context) {
	adminID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := h.userRepo.FindByID(uint(userID))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	// 使用 bcrypt 重置密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(DefaultResetPassword), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(c, "密码重置失败")
		return
	}
	user.PasswordHash = string(hashedPassword)
	if err := h.userRepo.Update(user); err != nil {
		response.InternalError(c, "密码重置失败")
		return
	}

	// 记录审计日志
	logAdminAudit(adminID, models.EventPasswordReset, "重置用户 "+user.Email+" 的密码", middleware.GetClientIP(c))

	response.OK(c, "密码已重置为: "+DefaultResetPassword)
}

// DisableUser 禁用用户账号
// POST /api/admin/users/:id/disable
func (h *AdminUserHandler) DisableUser(c *gin.Context) {
	adminID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := h.userRepo.FindByID(uint(userID))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	if user.IsDisabled {
		response.BadRequest(c, "用户已经是禁用状态")
		return
	}

	user.IsDisabled = true
	if err := h.userRepo.Update(user); err != nil {
		response.InternalError(c, "禁用用户失败")
		return
	}

	// 记录审计日志
	logAdminAudit(adminID, models.EventUserDisabled, "禁用用户 "+user.Email, middleware.GetClientIP(c))

	response.OK(c, "用户已禁用")
}

// EnableUser 启用用户账号
// POST /api/admin/users/:id/enable
func (h *AdminUserHandler) EnableUser(c *gin.Context) {
	adminID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := h.userRepo.FindByID(uint(userID))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	if !user.IsDisabled {
		response.BadRequest(c, "用户已经是启用状态")
		return
	}

	user.IsDisabled = false
	if err := h.userRepo.Update(user); err != nil {
		response.InternalError(c, "启用用户失败")
		return
	}

	// 记录审计日志
	logAdminAudit(adminID, models.EventUserEnabled, "启用用户 "+user.Email, middleware.GetClientIP(c))

	response.OK(c, "用户已启用")
}

// logAdminAudit 记录管理员操作审计日志
func logAdminAudit(userID uint, eventType, details, ipAddress string) {
	auditLog := models.AuditLog{
		UserID:    userID,
		EventType: eventType,
		Details:   details,
		IPAddress: ipAddress,
	}
	repository.GetDB().Create(&auditLog)
}
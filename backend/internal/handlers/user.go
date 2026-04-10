package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/internal/middleware"
	"github.com/taskify/backend/internal/services"
	"github.com/taskify/backend/pkg/response"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(),
	}
}

// GetCurrentUser 获取当前用户
// GET /api/users/me
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	userResp, err := h.userService.GetUserByID(userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, userResp)
}

// UpdateProfile 更新当前用户资料
// PUT /api/users/me
func (h *UserHandler) UpdateProfile(c *gin.Context) {
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

	userResp, err := h.userService.UpdateProfile(userID, req.DisplayName, req.AvatarURL)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, userResp)
}

// ChangePassword 修改密码
// PUT /api/users/me/password
func (h *UserHandler) ChangePassword(c *gin.Context) {
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

	err := h.userService.ChangePassword(userID, req.CurrentPassword, req.NewPassword)
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
package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/internal/services"
	"github.com/taskify/backend/pkg/response"
)

// PermissionHandler 权限处理器
type PermissionHandler struct {
	permissionService *services.PermissionService
}

// NewPermissionHandler 创建权限处理器
func NewPermissionHandler() *PermissionHandler {
	return &PermissionHandler{
		permissionService: services.NewPermissionService(),
	}
}

// GetPermissions 获取所有权限
// GET /api/permissions
func (h *PermissionHandler) GetPermissions(c *gin.Context) {
	permissions, err := h.permissionService.GetAllPermissions()
	if err != nil {
		response.InternalError(c, "获取权限列表失败")
		return
	}

	response.Success(c, gin.H{"permissions": permissions})
}
package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/internal/services"
	"github.com/taskify/backend/pkg/response"
)

// RoleHandler 角色处理器
type RoleHandler struct {
	roleService *services.RoleService
}

// NewRoleHandler 创建角色处理器
func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleService: services.NewRoleService(),
	}
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Scope       string `json:"scope" binding:"required"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SetPermissionsRequest 设置权限请求
type SetPermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids"`
}

// GetRoles 获取所有角色
// GET /api/roles
func (h *RoleHandler) GetRoles(c *gin.Context) {
	roles, err := h.roleService.GetAllRoles()
	if err != nil {
		response.InternalError(c, "获取角色列表失败")
		return
	}

	response.Success(c, gin.H{"roles": roles})
}

// GetRole 获取单个角色
// GET /api/roles/:id
func (h *RoleHandler) GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	role, err := h.roleService.GetRoleByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, role)
}

// CreateRole 创建角色
// POST /api/roles
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	role, err := h.roleService.CreateRole(req.Name, req.Description, req.Scope)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, "角色创建成功", role)
}

// UpdateRole 更新角色
// PUT /api/roles/:id
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	err = h.roleService.UpdateRole(uint(id), req.Name, req.Description)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "角色更新成功")
}

// DeleteRole 删除角色
// DELETE /api/roles/:id
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	err = h.roleService.DeleteRole(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "角色删除成功")
}

// SetRolePermissions 设置角色权限
// PUT /api/roles/:id/permissions
func (h *RoleHandler) SetRolePermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	var req SetPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	err = h.roleService.SetRolePermissions(uint(id), req.PermissionIDs)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "权限更新成功")
}
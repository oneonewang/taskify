package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/internal/middleware"
	"github.com/taskify/backend/internal/services"
	"github.com/taskify/backend/pkg/response"
)

// MembershipHandler 成员资格处理器
type MembershipHandler struct {
	membershipService *services.MembershipService
}

// NewMembershipHandler 创建成员资格处理器
func NewMembershipHandler() *MembershipHandler {
	return &MembershipHandler{
		membershipService: services.NewMembershipService(),
	}
}

// AssignSystemRoleRequest 分配系统角色请求
type AssignSystemRoleRequest struct {
	RoleID uint `json:"role_id" binding:"required"`
}

// AddProjectMemberRequest 添加项目成员请求
type AddProjectMemberRequest struct {
	UserEmail string `json:"user_email" binding:"required,email"`
	RoleID    uint   `json:"role_id" binding:"required"`
}

// UpdateMemberRoleRequest 更新成员角色请求
type UpdateMemberRoleRequest struct {
	RoleID uint `json:"role_id" binding:"required"`
}

// GetUserRoles 获取用户的角色
// GET /api/users/:id/roles
func (h *MembershipHandler) GetUserRoles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	memberships, err := h.membershipService.GetUserMemberships(uint(id))
	if err != nil {
		response.InternalError(c, "获取用户角色失败")
		return
	}

	response.Success(c, memberships)
}

// AssignSystemRole 分配系统角色
// POST /api/users/:id/roles
func (h *MembershipHandler) AssignSystemRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req AssignSystemRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	err = h.membershipService.AssignSystemRole(uint(id), req.RoleID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "角色分配成功")
}

// RemoveSystemRole 移除系统角色
// DELETE /api/users/:id/roles/:role_id
func (h *MembershipHandler) RemoveSystemRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	roleIDStr := c.Param("role_id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	err = h.membershipService.RemoveSystemRole(uint(id), uint(roleID))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "角色移除成功")
}

// GetUserProjectMemberships 获取用户的项目成员资格
// GET /api/users/:id/project-memberships
func (h *MembershipHandler) GetUserProjectMemberships(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	memberships, err := h.membershipService.GetUserMemberships(uint(id))
	if err != nil {
		response.InternalError(c, "获取项目成员资格失败")
		return
	}

	response.Success(c, memberships)
}

// GetProjectMembers 获取项目成员
// GET /api/projects/:id/members
func (h *MembershipHandler) GetProjectMembers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	members, err := h.membershipService.GetProjectMembers(uint(id))
	if err != nil {
		response.InternalError(c, "获取项目成员失败")
		return
	}

	response.Success(c, gin.H{"members": members})
}

// AddProjectMember 添加项目成员
// POST /api/projects/:id/members
func (h *MembershipHandler) AddProjectMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	var req AddProjectMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	err = h.membershipService.AddProjectMember(uint(id), req.UserEmail, req.RoleID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, "成员添加成功", nil)
}

// UpdateMemberRole 更新成员角色
// PUT /api/projects/:id/members/:user_id
func (h *MembershipHandler) UpdateMemberRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	err = h.membershipService.UpdateMemberRole(uint(id), uint(userID), req.RoleID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "角色更新成功")
}

// RemoveProjectMember 移除项目成员
// DELETE /api/projects/:id/members/:user_id
func (h *MembershipHandler) RemoveProjectMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	err = h.membershipService.RemoveProjectMember(uint(id), uint(userID))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "成员移除成功")
}

// GetCurrentUserRoles 获取当前用户的角色（临时使用简便方法）
func (h *MembershipHandler) GetCurrentUserRoles(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	memberships, err := h.membershipService.GetUserMemberships(userID)
	if err != nil {
		response.InternalError(c, "获取用户角色失败")
		return
	}

	response.Success(c, memberships)
}

// GetMyProjectMembership 获取当前用户在项目中的成员资格
// GET /api/projects/:id/my-membership
func (h *MembershipHandler) GetMyProjectMembership(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	idStr := c.Param("id")
	projectID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	membership, err := h.membershipService.GetUserMembershipForProject(userID, uint(projectID))
	if err != nil {
		// 用户不是项目成员
		response.Success(c, nil)
		return
	}

	response.Success(c, membership)
}
package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
	"github.com/taskify/backend/pkg/response"
)

// RequirePermission 权限检查中间件
func RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := GetUserID(c)
		if !ok {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// 检查用户是否是系统管理员
		if isAdmin, err := checkIsAdmin(userID); err == nil && isAdmin {
			c.Next()
			return
		}

		// 检查用户是否有指定权限
		hasPermission, err := checkUserPermission(userID, resource, action)
		if err != nil || !hasPermission {
			// 记录访问拒绝事件
			logAccessDenied(userID, resource, action, GetClientIP(c))
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireProjectMember 项目成员检查中间件
func RequireProjectMember() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := GetUserID(c)
		if !ok {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// 从URL参数获取项目ID
		projectIDStr := c.Param("id")
		projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
		if err != nil {
			response.BadRequest(c, "无效的项目ID")
			c.Abort()
			return
		}

		// 先检查是否是项目成员
		isMember, err := checkIsProjectMember(uint(userID), uint(projectID))
		if err == nil && isMember {
			c.Next()
			return
		}

		// 只有系统管理员可以绕过项目成员检查
		if isAdmin, err := checkIsAdmin(userID); err == nil && isAdmin {
			c.Next()
			return
		}

		response.Forbidden(c, "您不是该项目成员")
		c.Abort()
	}
}

// RequireProjectOwner 项目所有者检查中间件
func RequireProjectOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := GetUserID(c)
		if !ok {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// 从URL参数获取项目ID
		projectIDStr := c.Param("id")
		projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
		if err != nil {
			response.BadRequest(c, "无效的项目ID")
			c.Abort()
			return
		}

		// 先检查是否是项目所有者（管理员也是项目所有者，跳过此检查）
		isOwner, err := checkIsProjectOwner(uint(userID), uint(projectID))
		if err == nil && isOwner {
			c.Next()
			return
		}

		// 只有系统管理员可以绕过项目所有者检查
		if isAdmin, err := checkIsAdmin(userID); err == nil && isAdmin {
			c.Next()
			return
		}

		response.Forbidden(c, "只有项目所有者才能执行此操作")
		c.Abort()
	}
}

// RequireAdmin 系统管理员检查中间件
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := GetUserID(c)
		if !ok {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		isAdmin, err := checkIsAdmin(userID)
		if err != nil || !isAdmin {
			response.Forbidden(c, "需要系统管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// CheckIsAdmin 检查用户是否是系统管理员（导出版本）
func CheckIsAdmin(userID uint) (bool, error) {
	return checkIsAdmin(userID)
}

// CheckIsProjectMember 检查用户是否是项目成员（导出版本）
func CheckIsProjectMember(userID, projectID uint) (bool, error) {
	return checkIsProjectMember(userID, projectID)
}

// checkIsAdmin 检查用户是否是系统管理员
func checkIsAdmin(userID uint) (bool, error) {
	var count int64
	err := repository.GetDB().
		Model(&models.ProjectMembership{}).
		Joins("JOIN roles ON project_memberships.role_id = roles.id").
		Where("project_memberships.user_id = ? AND roles.name = ? AND roles.scope = ?", userID, models.RoleAdmin, "system").
		Count(&count).Error
	return count > 0, err
}

// checkIsProjectMember 检查用户是否是项目成员
func checkIsProjectMember(userID, projectID uint) (bool, error) {
	var count int64
	err := repository.GetDB().
		Model(&models.ProjectMembership{}).
		Where("user_id = ? AND project_id = ?", userID, projectID).
		Count(&count).Error
	return count > 0, err
}

// checkIsProjectOwner 检查用户是否是项目所有者
func checkIsProjectOwner(userID, projectID uint) (bool, error) {
	var count int64
	err := repository.GetDB().
		Model(&models.ProjectMembership{}).
		Joins("JOIN roles ON project_memberships.role_id = roles.id").
		Where("project_memberships.user_id = ? AND project_memberships.project_id = ? AND roles.name = ?", userID, projectID, models.RoleOwner).
		Count(&count).Error
	return count > 0, err
}

// checkUserPermission 检查用户是否有指定权限
func checkUserPermission(userID uint, resource, action string) (bool, error) {
	// 获取用户的所有项目成员资格
	var memberships []models.ProjectMembership
	err := repository.GetDB().
		Where("user_id = ?", userID).
		Find(&memberships).Error
	if err != nil {
		return false, err
	}

	if len(memberships) == 0 {
		return false, nil
	}

	// 收集用户的所有角色ID
	roleIDs := make([]uint, len(memberships))
	for i, m := range memberships {
		roleIDs[i] = m.RoleID
	}

	// 检查是否有所需权限
	var count int64
	permKey := resource + "." + action
	err = repository.GetDB().
		Model(&models.Permission{}).
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id IN ? AND CONCAT(permissions.resource, '.', permissions.action) = ?", roleIDs, permKey).
		Count(&count).Error

	return count > 0, err
}

// logAccessDenied 记录访问拒绝事件
func logAccessDenied(userID uint, resource, action, ipAddress string) {
	details := "尝试访问 " + resource + "." + action
	auditLog := models.AuditLog{
		UserID:    userID,
		EventType: models.EventAccessDenied,
		Details:   details,
		IPAddress: ipAddress,
	}
	repository.GetDB().Create(&auditLog)
}
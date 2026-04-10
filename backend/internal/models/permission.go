package models

import (
	"time"
)

// Permission 权限模型
type Permission struct {
	ID          uint      `gorm:"primaryKey"`                     // 唯一标识符
	Resource    string    `gorm:"index;size:50;not null"`        // 资源类型
	Action      string    `gorm:"index;size:50;not null"`        // 操作名称
	Description string    `gorm:"size:255"`                       // 权限描述
	CreatedAt   time.Time                                      // 创建时间戳
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// PermissionResponse 权限API响应
type PermissionResponse struct {
	ID          uint   `json:"id"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

// ToResponse 转换为API响应格式
func (p *Permission) ToResponse() PermissionResponse {
	return PermissionResponse{
		ID:          p.ID,
		Resource:    p.Resource,
		Action:      p.Action,
		Description: p.Description,
	}
}

// FormatPermission 格式化权限字符串 (resource.action)
func FormatPermission(resource, action string) string {
	return resource + "." + action
}

// 权限资源类型
const (
	ResourceUsers    = "users"
	ResourceProjects = "projects"
	ResourceTasks    = "tasks"
	ResourceComments = "comments"
	ResourceRoles    = "roles"
)

// 权限操作类型
const (
	ActionView    = "view"
	ActionCreate  = "create"
	ActionEdit    = "edit"
	ActionDelete  = "delete"
	ActionManage  = "manage"
	ActionAssign  = "assign"
	ActionMove    = "move"
)

// 预定义权限列表
var SystemPermissions = []Permission{
	// 用户权限
	{Resource: ResourceUsers, Action: ActionView, Description: "查看用户"},
	{Resource: ResourceUsers, Action: ActionManage, Description: "管理用户"},
	// 项目权限
	{Resource: ResourceProjects, Action: ActionView, Description: "查看项目"},
	{Resource: ResourceProjects, Action: ActionCreate, Description: "创建项目"},
	{Resource: ResourceProjects, Action: ActionEdit, Description: "编辑项目"},
	{Resource: ResourceProjects, Action: ActionDelete, Description: "删除项目"},
	{Resource: ResourceProjects, Action: ActionManage, Description: "管理项目"},
	// 任务权限
	{Resource: ResourceTasks, Action: ActionView, Description: "查看任务"},
	{Resource: ResourceTasks, Action: ActionCreate, Description: "创建任务"},
	{Resource: ResourceTasks, Action: ActionEdit, Description: "编辑任务"},
	{Resource: ResourceTasks, Action: ActionMove, Description: "移动任务"},
	{Resource: ResourceTasks, Action: ActionDelete, Description: "删除任务"},
	// 评论权限
	{Resource: ResourceComments, Action: ActionView, Description: "查看评论"},
	{Resource: ResourceComments, Action: ActionCreate, Description: "创建评论"},
	{Resource: ResourceComments, Action: ActionEdit, Description: "编辑评论"},
	{Resource: ResourceComments, Action: ActionDelete, Description: "删除评论"},
	// 角色权限
	{Resource: ResourceRoles, Action: ActionView, Description: "查看角色"},
	{Resource: ResourceRoles, Action: ActionCreate, Description: "创建角色"},
	{Resource: ResourceRoles, Action: ActionEdit, Description: "编辑角色"},
	{Resource: ResourceRoles, Action: ActionDelete, Description: "删除角色"},
	{Resource: ResourceRoles, Action: ActionAssign, Description: "分配角色"},
}
package models

import (
	"time"
)

// Role 角色模型
type Role struct {
	ID          uint      `gorm:"primaryKey"`                               // 唯一标识符
	Name        string    `gorm:"uniqueIndex;size:50;not null"`              // 角色名称
	DisplayName string    `gorm:"size:50"`                                   // 显示名称
	Description string    `gorm:"size:255"`                                 // 角色描述
	IsSystem    bool      `gorm:"default:false"`                            // 是否系统预定义
	Scope       string    `gorm:"size:20;default:'system'"`                  // 作用域: system/project
	CreatedAt   time.Time                                              // 创建时间戳
	UpdatedAt   time.Time                                              // 更新时间戳
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// RoleResponse 角色API响应
type RoleResponse struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Description string   `json:"description"`
	IsSystem    bool     `json:"is_system"`
	Scope       string   `json:"scope"`
	Permissions []string `json:"permissions,omitempty"`
}

// ToResponse 转换为API响应格式
func (r *Role) ToResponse() RoleResponse {
	return RoleResponse{
		ID:          r.ID,
		Name:        r.Name,
		DisplayName: r.DisplayName,
		Description: r.Description,
		IsSystem:    r.IsSystem,
		Scope:       r.Scope,
	}
}

// 预定义角色常量
const (
	RoleAdmin   = "admin"   // 系统管理员
	RoleOwner   = "owner"   // 项目所有者
	RoleMember  = "member"  // 项目成员
	RoleGuest   = "guest"   // 项目访客
)

// 预定义角色列表
var SystemRoles = []Role{
	{Name: RoleAdmin, DisplayName: "系统管理员", Description: "系统管理员", IsSystem: true, Scope: "system"},
	{Name: RoleOwner, DisplayName: "项目所有者", Description: "项目所有者", IsSystem: true, Scope: "project"},
	{Name: RoleMember, DisplayName: "项目成员", Description: "项目成员", IsSystem: true, Scope: "project"},
	{Name: RoleGuest, DisplayName: "项目访客", Description: "项目访客", IsSystem: true, Scope: "project"},
}
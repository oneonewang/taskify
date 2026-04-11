package models

import (
	"time"
)

// ProjectMembership 项目成员资格
type ProjectMembership struct {
	ID        uint      `gorm:"primaryKey"`                                                      // 唯一标识符
	UserID    uint      `gorm:"uniqueIndex:idx_pm_user_project;index:idx_pm_project_user_role"` // 用户ID
	ProjectID uint      `gorm:"uniqueIndex:idx_pm_user_project;index:idx_pm_project_user_role"` // 项目ID
	RoleID    uint      `gorm:"index:idx_pm_project_user_role"`                                // 角色ID
	JoinedAt  time.Time                                                               // 加入时间
}

// TableName 指定表名
func (ProjectMembership) TableName() string {
	return "project_memberships"
}

// ProjectMembershipResponse 项目成员API响应
type ProjectMembershipResponse struct {
	UserID      uint      `json:"user_id"`
	DisplayName string    `json:"display_name"`
	AvatarURL   string    `json:"avatar_url"`
	Role        string    `json:"role"`
	JoinedAt    time.Time `json:"joined_at"`
}
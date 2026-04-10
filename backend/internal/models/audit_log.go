package models

import (
	"time"
)

// AuditLog 审计日志
type AuditLog struct {
	ID        uint      `gorm:"primaryKey"`                     // 唯一标识符
	UserID    uint      `gorm:"index"`                          // 用户ID
	EventType string    `gorm:"index;size:50"`                  // 事件类型
	Details   string    `gorm:"type:text"`                     // 事件详情
	IPAddress string    `gorm:"size:45"`                        // IP地址 (IPv6兼容)
	CreatedAt time.Time                                      // 创建时间
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}

// AuditLogResponse 审计日志API响应
type AuditLogResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	UserEmail string    `json:"user_email,omitempty"`
	EventType string    `json:"event_type"`
	Details   string    `json:"details"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse 转换为API响应格式
func (a *AuditLog) ToResponse() AuditLogResponse {
	return AuditLogResponse{
		ID:        a.ID,
		UserID:    a.UserID,
		EventType: a.EventType,
		Details:   a.Details,
		IPAddress: a.IPAddress,
		CreatedAt: a.CreatedAt,
	}
}

// 审计日志事件类型
const (
	EventLogin            = "login"              // 用户登录
	EventLogout          = "logout"             // 用户登出
	EventPermissionChange = "permission_change" // 权限变更
	EventAccessDenied    = "access_denied"      // 访问被拒绝
	EventProjectCreated  = "project_created"    // 项目创建
	EventProjectDeleted  = "project_deleted"    // 项目删除
	EventProjectArchived = "project_archived"   // 项目归档
	EventPasswordReset   = "password_reset"    // 密码重置
	EventUserDisabled    = "user_disabled"      // 用户禁用
	EventUserEnabled     = "user_enabled"      // 用户启用
	EventUserImport     = "user_import"       // 用户批量导入
)
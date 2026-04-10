package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID            uint       `gorm:"primaryKey"`                     // 唯一标识符
	Email         string     `gorm:"uniqueIndex;size:255;not null"` // 邮箱地址
	DisplayName   string     `gorm:"size:100"`                       // 显示名称
	AvatarURL     string     `gorm:"size:500"`                       // 头像 URL
	PasswordHash  string     `gorm:"size:255;not null"`             // bcrypt 哈希密码
	EmailVerified bool       `gorm:"default:false"`                 // 邮箱已验证
	IsDisabled    bool       `gorm:"default:false"`                 // 账号是否被禁用
	LastLoginAt   *time.Time                                // 最近登录时间戳
	CreatedAt     time.Time                                  // 创建时间戳
	UpdatedAt     time.Time                                  // 更新时间戳
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserResponse 用户API响应
type UserResponse struct {
	ID            uint      `json:"id"`
	Email         string    `json:"email"`
	DisplayName   string    `json:"display_name"`
	AvatarURL     string    `json:"avatar_url"`
	EmailVerified bool      `json:"email_verified"`
	IsDisabled    bool      `json:"is_disabled"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// ToResponse 转换为API响应格式
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		DisplayName:   u.DisplayName,
		AvatarURL:     u.AvatarURL,
		EmailVerified: u.EmailVerified,
		IsDisabled:    u.IsDisabled,
		LastLoginAt:   u.LastLoginAt,
		CreatedAt:     u.CreatedAt,
	}
}
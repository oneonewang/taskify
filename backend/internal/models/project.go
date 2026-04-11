package models

import (
	"time"
)

// Project 项目模型
type Project struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	Description string    `json:"description" gorm:"type:varchar(500)"`
	IsArchived  bool      `json:"is_archived" gorm:"default:false;index"` // 是否已归档
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Project) TableName() string {
	return "projects"
}

// ProjectResponse 项目API响应
type ProjectResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsArchived  bool      `json:"is_archived"`
	CreatedAt   time.Time `json:"created_at"`
	OwnerID     uint      `json:"owner_id,omitempty"`
	OwnerName   string    `json:"owner_name,omitempty"`
	OwnerAvatar string    `json:"owner_avatar,omitempty"`
}

// ToResponse 转换为API响应格式
func (p *Project) ToResponse() ProjectResponse {
	return ProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		IsArchived:  p.IsArchived,
		CreatedAt:   p.CreatedAt,
	}
}
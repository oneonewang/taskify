package models

import (
	"time"
)

// Project 项目模型
type Project struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	Description string    `json:"description" gorm:"type:varchar(500)"`
	IsArchived  bool      `json:"is_archived" gorm:"default:false"` // 是否已归档
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

// 预定义项目数据（保留用于迁移参考）
var SeedProjects = []Project{
	{ID: 1, Name: "任务管理重构", Description: "重构现有任务管理系统，提升性能和可维护性"},
	{ID: 2, Name: "移动端开发", Description: "开发iOS和Android移动应用"},
	{ID: 3, Name: "API网关升级", Description: "升级API网关，支持更多协议和认证方式"},
}
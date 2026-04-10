package models

import (
	"time"
)

// TaskStatus 任务状态枚举
type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusReview     TaskStatus = "review"
	StatusDone       TaskStatus = "done"
)

// Task 任务模型
type Task struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string    `json:"title" gorm:"type:varchar(200);not null"`
	Description string    `json:"description" gorm:"type:varchar(2000)"`
	Status     TaskStatus `json:"status" gorm:"type:varchar(50);not null;default:'todo'"`
	Position    int       `json:"position" gorm:"default:0"`
	AssigneeID  uint      `json:"assignee_id" gorm:"not null"`
	ProjectID   uint      `json:"project_id" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// 关联关系
	Assignee *User `json:"assignee,omitempty" gorm:"foreignKey:AssigneeID"`
	Project  *Project `json:"-" gorm:"foreignKey:ProjectID"`
}

// TableName 指定表名
func (Task) TableName() string {
	return "tasks"
}

// AssigneeInfo 负责人信息
type AssigneeInfo struct {
	ID          uint   `json:"id"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

// TaskResponse 任务API响应
type TaskResponse struct {
	ID          uint        `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	Status     TaskStatus  `json:"status"`
	Position    int         `json:"position"`
	Assignee    AssigneeInfo `json:"assignee"`
	ProjectID   uint        `json:"project_id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// ToResponse 转换为API响应格式
func (t *Task) ToResponse() TaskResponse {
	resp := TaskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Position:    t.Position,
		ProjectID:   t.ProjectID,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
	if t.Assignee != nil {
		resp.Assignee = AssigneeInfo{
			ID:          t.Assignee.ID,
			DisplayName: t.Assignee.DisplayName,
			AvatarURL:   t.Assignee.AvatarURL,
		}
	}
	return resp
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required,max=200"`
	Description string `json:"description" binding:"max=2000"`
	AssigneeID  uint   `json:"assignee_id" binding:"required"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	Title       *string `json:"title" binding:"omitempty,max=200"`
	Description *string `json:"description" binding:"omitempty,max=2000"`
	AssigneeID  *uint   `json:"assignee_id"`
}

// UpdateStatusRequest 更新状态请求
type UpdateStatusRequest struct {
	Status   TaskStatus `json:"status" binding:"required"`
	Position int        `json:"position" binding:"gte=0"`
}

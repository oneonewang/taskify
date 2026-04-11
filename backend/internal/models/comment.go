package models

import (
	"time"
)

// Comment 评论模型
type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Content   string    `json:"content" gorm:"type:varchar(2000);not null"`
	UserID    uint      `json:"user_id" gorm:"index"`     // 用户ID
	TaskID    uint      `json:"task_id" gorm:"index"`     // 任务ID
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// 关联关系
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Task *Task `json:"-" gorm:"foreignKey:TaskID"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}

// UserInfo 用户信息
type UserInfo struct {
	ID          uint   `json:"id"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
}

// CommentResponse 评论API响应
type CommentResponse struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	User      UserInfo  `json:"user"`
	TaskID    uint      `json:"task_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse 转换为API响应格式
func (c *Comment) ToResponse() CommentResponse {
	resp := CommentResponse{
		ID:        c.ID,
		Content:   c.Content,
		TaskID:    c.TaskID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
	if c.User != nil {
		resp.User = UserInfo{
			ID:          c.User.ID,
			DisplayName: c.User.DisplayName,
			AvatarURL:   c.User.AvatarURL,
		}
	}
	return resp
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,max=2000"`
}

// UpdateCommentRequest 更新评论请求
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,max=2000"`
}

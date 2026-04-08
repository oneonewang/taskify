package models

import (
	"time"
)

// UserRole 用户角色枚举
type UserRole string

const (
	RoleProductManager UserRole = "product_manager"
	RoleEngineer       UserRole = "engineer"
)

// User 用户模型
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(50);not null"`
	Role      UserRole  `json:"role" gorm:"type:varchar(50);not null"`
	Avatar    string    `json:"avatar" gorm:"type:varchar(100)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserResponse 用户API响应
type UserResponse struct {
	ID     uint     `json:"id"`
	Name   string   `json:"name"`
	Role   UserRole `json:"role"`
	Avatar string   `json:"avatar"`
}

// ToResponse 转换为API响应格式
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:     u.ID,
		Name:   u.Name,
		Role:   u.Role,
		Avatar: u.Avatar,
	}
}

//预定义用户数据
var SeedUsers = []User{
	{ID: 1, Name: "张明", Role: RoleProductManager, Avatar: "#FF6B6B"},
	{ID: 2, Name: "李伟", Role: RoleEngineer, Avatar: "#4ECDC4"},
	{ID: 3, Name: "王芳", Role: RoleEngineer, Avatar: "#45B7D1"},
	{ID: 4, Name: "刘强", Role: RoleEngineer, Avatar: "#96CEB4"},
	{ID: 5, Name: "陈静", Role: RoleEngineer, Avatar: "#FFEAA7"},
}

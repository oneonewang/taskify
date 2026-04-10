package models

// RolePermission 角色权限关联表
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey"` // 角色ID
	PermissionID uint `gorm:"primaryKey"` // 权限ID
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}
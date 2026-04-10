package repository

import (
	"github.com/taskify/backend/internal/models"
	"gorm.io/gorm"
)

// RoleRepository 角色数据访问层
type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓库
func NewRoleRepository() *RoleRepository {
	return &RoleRepository{db: GetDB()}
}

// FindByID 根据ID查找角色
func (r *RoleRepository) FindByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// FindByName 根据名称查找角色
func (r *RoleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetAll 获取所有角色
func (r *RoleRepository) GetAll() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Find(&roles).Error
	return roles, err
}

// Create 创建角色
func (r *RoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

// Update 更新角色
func (r *RoleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

// Delete 删除角色
func (r *RoleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

// GetPermissions 获取角色的权限
func (r *RoleRepository) GetPermissions(roleID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.
		Model(&models.Permission{}).
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}

// AddPermission 添加角色权限
func (r *RoleRepository) AddPermission(roleID, permissionID uint) error {
	rp := models.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	return r.db.Create(&rp).Error
}

// RemovePermission 移除角色权限
func (r *RoleRepository) RemovePermission(roleID, permissionID uint) error {
	return r.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&models.RolePermission{}).Error
}

// SetPermissions 设置角色权限（先删除后添加）
func (r *RoleRepository) SetPermissions(roleID uint, permissionIDs []uint) error {
	// 开启事务
	tx := r.db.Begin()

	// 删除现有权限
	if err := tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 添加新权限
	for _, permID := range permissionIDs {
		rp := models.RolePermission{
			RoleID:       roleID,
			PermissionID: permID,
		}
		if err := tx.Create(&rp).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

// IsAssigned 检查角色是否已被分配
func (r *RoleRepository) IsAssigned(roleID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.ProjectMembership{}).Where("role_id = ?", roleID).Count(&count).Error
	return count > 0, err
}

// CountAdmins 统计系统管理员数量
func (r *RoleRepository) CountAdmins() (int64, error) {
	var count int64
	adminRole := models.Role{}
	if err := r.db.Where("name = ? AND scope = ?", models.RoleAdmin, "system").First(&adminRole).Error; err != nil {
		return 0, err
	}
	err := r.db.Model(&models.ProjectMembership{}).Where("role_id = ?", adminRole.ID).Count(&count).Error
	return count, err
}
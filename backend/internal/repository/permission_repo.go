package repository

import (
	"github.com/taskify/backend/internal/models"
	"gorm.io/gorm"
)

// PermissionRepository 权限数据访问层
type PermissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository 创建权限仓库
func NewPermissionRepository() *PermissionRepository {
	return &PermissionRepository{db: GetDB()}
}

// FindByID 根据ID查找权限
func (r *PermissionRepository) FindByID(id uint) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetAll 获取所有权限
func (r *PermissionRepository) GetAll() ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Find(&permissions).Error
	return permissions, err
}

// GetByResource 按资源类型获取权限
func (r *PermissionRepository) GetByResource(resource string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Where("resource = ?", resource).Find(&permissions).Error
	return permissions, err
}
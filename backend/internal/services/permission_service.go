package services

import (
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
)

// PermissionService 权限服务
type PermissionService struct {
	permissionRepo *repository.PermissionRepository
}

// NewPermissionService 创建权限服务
func NewPermissionService() *PermissionService {
	return &PermissionService{
		permissionRepo: repository.NewPermissionRepository(),
	}
}

// GetAllPermissions 获取所有权限
func (s *PermissionService) GetAllPermissions() ([]models.PermissionResponse, error) {
	permissions, err := s.permissionRepo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]models.PermissionResponse, len(permissions))
	for i, p := range permissions {
		responses[i] = p.ToResponse()
	}
	return responses, nil
}

// GetPermissionsByResource 按资源获取权限
func (s *PermissionService) GetPermissionsByResource(resource string) ([]models.PermissionResponse, error) {
	permissions, err := s.permissionRepo.GetByResource(resource)
	if err != nil {
		return nil, err
	}

	responses := make([]models.PermissionResponse, len(permissions))
	for i, p := range permissions {
		responses[i] = p.ToResponse()
	}
	return responses, nil
}
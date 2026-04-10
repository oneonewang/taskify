package services

import (
	"errors"

	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
)

// RoleService 角色服务
type RoleService struct {
	roleRepo *repository.RoleRepository
}

// NewRoleService 创建角色服务
func NewRoleService() *RoleService {
	return &RoleService{
		roleRepo: repository.NewRoleRepository(),
	}
}

// GetAllRoles 获取所有角色
func (s *RoleService) GetAllRoles() ([]models.RoleResponse, error) {
	roles, err := s.roleRepo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]models.RoleResponse, len(roles))
	for i, role := range roles {
		perms, _ := s.roleRepo.GetPermissions(role.ID)
		permStrings := make([]string, len(perms))
		for j, p := range perms {
			permStrings[j] = models.FormatPermission(p.Resource, p.Action)
		}
		responses[i] = role.ToResponse()
		responses[i].Permissions = permStrings
	}
	return responses, nil
}

// GetRoleByID 根据ID获取角色
func (s *RoleService) GetRoleByID(id uint) (*models.RoleResponse, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("角色不存在")
	}

	resp := role.ToResponse()
	perms, _ := s.roleRepo.GetPermissions(role.ID)
	permStrings := make([]string, len(perms))
	for j, p := range perms {
		permStrings[j] = models.FormatPermission(p.Resource, p.Action)
	}
	resp.Permissions = permStrings
	return &resp, nil
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(name, description, scope string) (*models.RoleResponse, error) {
	// 检查名称是否已存在
	existing, _ := s.roleRepo.FindByName(name)
	if existing != nil {
		return nil, errors.New("角色名称已存在")
	}

	role := &models.Role{
		Name:        name,
		Description: description,
		Scope:       scope,
		IsSystem:    false,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}

	resp := role.ToResponse()
	return &resp, nil
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(id uint, name, description string) error {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return errors.New("角色不存在")
	}

	if role.IsSystem {
		return errors.New("无法更新系统预定义角色")
	}

	if name != "" {
		role.Name = name
	}
	if description != "" {
		role.Description = description
	}

	return s.roleRepo.Update(role)
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(id uint) error {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 禁止删除系统预定义角色
	if role.IsSystem {
		return errors.New("无法删除系统预定义角色")
	}

	// 检查是否已被分配
	assigned, err := s.roleRepo.IsAssigned(id)
	if err != nil {
		return err
	}
	if assigned {
		return errors.New("角色已被分配给用户，无法删除")
	}

	return s.roleRepo.Delete(id)
}

// SetRolePermissions 设置角色权限
func (s *RoleService) SetRolePermissions(roleID uint, permissionIDs []uint) error {
	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("角色不存在")
	}

	if role.IsSystem {
		return errors.New("无法修改系统预定义角色的权限")
	}

	return s.roleRepo.SetPermissions(roleID, permissionIDs)
}
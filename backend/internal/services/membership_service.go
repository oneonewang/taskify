package services

import (
	"errors"

	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
)

// MembershipService 成员资格服务
type MembershipService struct {
	membershipRepo *repository.MembershipRepository
	userRepo       *repository.UserRepository
	roleRepo       *repository.RoleRepository
}

// NewMembershipService 创建成员资格服务
func NewMembershipService() *MembershipService {
	return &MembershipService{
		membershipRepo: repository.NewMembershipRepository(),
		userRepo:       repository.NewUserRepository(),
		roleRepo:       repository.NewRoleRepository(),
	}
}

// AssignSystemRole 分配系统角色
func (s *MembershipService) AssignSystemRole(userID uint, roleID uint) error {
	// 检查用户是否存在
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 检查角色是否存在
	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 检查是否是系统角色
	if role.Scope != "system" {
		return errors.New("只能分配系统角色")
	}

	// 为用户分配系统角色（使用特殊项目ID 0 表示系统角色）
	return s.membershipRepo.AssignRole(userID, 0, roleID)
}

// RemoveSystemRole 移除系统角色
func (s *MembershipService) RemoveSystemRole(userID uint, roleID uint) error {
	role, err := s.roleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("角色不存在")
	}

	if role.Scope != "system" {
		return errors.New("只能移除系统角色")
	}

	// 检查是否是最后一个管理员
	if role.Name == models.RoleAdmin {
		count, err := s.membershipRepo.CountByRole(roleID)
		if err != nil {
			return err
		}
		if count <= 1 {
			return errors.New("无法移除最后一个系统管理员")
		}
	}

	return s.membershipRepo.Delete(userID, 0)
}

// AddProjectMember 添加项目成员
func (s *MembershipService) AddProjectMember(projectID uint, userEmail string, roleID uint) error {
	// 查找用户
	user, err := s.userRepo.FindByEmail(userEmail)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 检查角色是否存在
	_, err = s.roleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 检查是否已是项目成员
	existing, _ := s.membershipRepo.FindByUserAndProject(user.ID, projectID)
	if existing != nil {
		return errors.New("用户已是项目成员")
	}

	return s.membershipRepo.AssignRole(user.ID, projectID, roleID)
}

// RemoveProjectMember 移除项目成员
func (s *MembershipService) RemoveProjectMember(projectID, userID uint) error {
	return s.membershipRepo.Delete(userID, projectID)
}

// UpdateMemberRole 更新成员角色
func (s *MembershipService) UpdateMemberRole(projectID, userID, roleID uint) error {
	// 检查成员资格是否存在
	membership, err := s.membershipRepo.FindByUserAndProject(userID, projectID)
	if err != nil {
		return errors.New("用户不是项目成员")
	}

	// 检查新角色是否存在
	_, err = s.roleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("角色不存在")
	}

	membership.RoleID = roleID
	return s.membershipRepo.Update(membership)
}

// GetUserMemberships 获取用户的成员资格列表
func (s *MembershipService) GetUserMemberships(userID uint) (map[string]interface{}, error) {
	memberships, err := s.membershipRepo.GetUserMembershipsWithDetails(userID)
	if err != nil {
		return nil, err
	}

	// 分离系统角色和项目成员资格
	var systemRoles []string
	var projectMemberships []map[string]interface{}

	for _, m := range memberships {
		if projectID, ok := m["project_id"].(uint64); ok && projectID == 0 {
			if roleName, ok := m["role_name"].(string); ok {
				systemRoles = append(systemRoles, roleName)
			}
		} else {
			projectMemberships = append(projectMemberships, m)
		}
	}

	return map[string]interface{}{
		"system_roles":       systemRoles,
		"project_memberships": projectMemberships,
	}, nil
}

// GetProjectMembers 获取项目的成员列表
func (s *MembershipService) GetProjectMembers(projectID uint) ([]map[string]interface{}, error) {
	return s.membershipRepo.GetProjectMembershipsWithDetails(projectID)
}
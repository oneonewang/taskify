package repository

import (
	"time"

	"github.com/taskify/backend/internal/models"
	"gorm.io/gorm"
)

// MembershipRepository 项目成员资格数据访问层
type MembershipRepository struct {
	db *gorm.DB
}

// NewMembershipRepository 创建成员资格仓库
func NewMembershipRepository() *MembershipRepository {
	return &MembershipRepository{db: GetDB()}
}

// Create 创建成员资格
func (r *MembershipRepository) Create(membership *models.ProjectMembership) error {
	return r.db.Create(membership).Error
}

// Delete 删除成员资格
func (r *MembershipRepository) Delete(userID, projectID uint) error {
	return r.db.Where("user_id = ? AND project_id = ?", userID, projectID).Delete(&models.ProjectMembership{}).Error
}

// FindByUserAndProject 查找特定用户和项目的成员资格
func (r *MembershipRepository) FindByUserAndProject(userID, projectID uint) (*models.ProjectMembership, error) {
	var membership models.ProjectMembership
	err := r.db.Where("user_id = ? AND project_id = ?", userID, projectID).First(&membership).Error
	if err != nil {
		return nil, err
	}
	return &membership, nil
}

// GetByUser 获取用户的所有成员资格
func (r *MembershipRepository) GetByUser(userID uint) ([]models.ProjectMembership, error) {
	var memberships []models.ProjectMembership
	err := r.db.Where("user_id = ?", userID).Find(&memberships).Error
	return memberships, err
}

// GetByProject 获取项目的所有成员资格
func (r *MembershipRepository) GetByProject(projectID uint) ([]models.ProjectMembership, error) {
	var memberships []models.ProjectMembership
	err := r.db.Where("project_id = ?", projectID).Find(&memberships).Error
	return memberships, err
}

// Update 更新成员资格
func (r *MembershipRepository) Update(membership *models.ProjectMembership) error {
	return r.db.Save(membership).Error
}

// CountByRole 统计具有特定角色的成员数量
func (r *MembershipRepository) CountByRole(roleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.ProjectMembership{}).Where("role_id = ?", roleID).Count(&count).Error
	return count, err
}

// IsAdmin 检查用户是否是系统管理员
func (r *MembershipRepository) IsAdmin(userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.ProjectMembership{}).
		Joins("JOIN roles ON project_memberships.role_id = roles.id").
		Where("project_memberships.user_id = ? AND roles.name = ? AND roles.scope = ?", userID, models.RoleAdmin, "system").
		Count(&count).Error
	return count > 0, err
}

// IsMember 检查用户是否是项目成员
func (r *MembershipRepository) IsMember(userID, projectID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.ProjectMembership{}).
		Where("user_id = ? AND project_id = ?", userID, projectID).
		Count(&count).Error
	return count > 0, err
}

// GetUserMembershipsWithDetails 获取用户的成员资格详情
func (r *MembershipRepository) GetUserMembershipsWithDetails(userID uint) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := r.db.Table("project_memberships").
		Select("project_memberships.*, projects.name as project_name, roles.name as role_name, roles.display_name as role_display_name").
		Joins("JOIN projects ON project_memberships.project_id = projects.id").
		Joins("JOIN roles ON project_memberships.role_id = roles.id").
		Where("project_memberships.user_id = ?", userID).
		Scan(&results).Error
	return results, err
}

// GetProjectMembershipsWithDetails 获取项目的成员资格详情
func (r *MembershipRepository) GetProjectMembershipsWithDetails(projectID uint) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := r.db.Table("project_memberships").
		Select("project_memberships.*, users.email, users.display_name, users.avatar_url, roles.name as role_name, roles.display_name as role_display_name").
		Joins("JOIN users ON project_memberships.user_id = users.id").
		Joins("JOIN roles ON project_memberships.role_id = roles.id").
		Where("project_memberships.project_id = ?", projectID).
		Scan(&results).Error
	return results, err
}

// GetUserMembershipForProject 获取用户在特定项目中的成员资格详情
func (r *MembershipRepository) GetUserMembershipForProject(userID, projectID uint) (*models.ProjectMembership, map[string]interface{}, error) {
	var membership models.ProjectMembership
	err := r.db.Where("user_id = ? AND project_id = ?", userID, projectID).First(&membership).Error
	if err != nil {
		return nil, nil, err
	}

	var details map[string]interface{}
	err = r.db.Table("project_memberships").
		Select("project_memberships.*, projects.name as project_name, roles.name as role_name, roles.display_name as role_display_name").
		Joins("JOIN projects ON project_memberships.project_id = projects.id").
		Joins("JOIN roles ON project_memberships.role_id = roles.id").
		Where("project_memberships.user_id = ? AND project_memberships.project_id = ?", userID, projectID).
		Scan(&details).Error
	if err != nil {
		return nil, nil, err
	}

	return &membership, details, nil
}

// AssignRole 分配角色
func (r *MembershipRepository) AssignRole(userID, projectID, roleID uint) error {
	// 检查是否已存在成员资格
	existing, _ := r.FindByUserAndProject(userID, projectID)
	if existing != nil {
		existing.RoleID = roleID
		existing.JoinedAt = time.Now()
		return r.Update(existing)
	}

	membership := &models.ProjectMembership{
		UserID:    userID,
		ProjectID: projectID,
		RoleID:    roleID,
		JoinedAt:  time.Now(),
	}
	return r.Create(membership)
}
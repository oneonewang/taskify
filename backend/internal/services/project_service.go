package services

import (
	"errors"
	"fmt"

	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
	"gorm.io/gorm"
)

// ProjectListFilter 项目列表过滤条件
type ProjectListFilter struct {
	Keyword    string // 项目名称关键字（模糊搜索）
	IsArchived *bool  // 是否归档（nil 表示不过滤）
	OwnerName  string // 所有者名称（模糊搜索）
	UserID     uint   // 用户ID（大于0时过滤该用户参与的项目）
	Page       int    // 页码（从1开始）
	PageSize   int    // 每页数量
}

// ProjectService 项目服务
type ProjectService struct {
	projectRepo    *repository.ProjectRepository
	userRepo       *repository.UserRepository
	membershipRepo *repository.MembershipRepository
}

// NewProjectService 创建项目服务
func NewProjectService() *ProjectService {
	return &ProjectService{
		projectRepo:    repository.NewProjectRepository(),
		userRepo:       repository.NewUserRepository(),
		membershipRepo: repository.NewMembershipRepository(),
	}
}

// GetAllUsers 获取所有用户
func (s *ProjectService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}

// GetAllProjects 获取所有未归档的项目
func (s *ProjectService) GetAllProjects() ([]models.Project, error) {
	return s.projectRepo.GetAll()
}

// ListProjects 分页获取项目列表（包含所有者信息）
func (s *ProjectService) ListProjects(filter ProjectListFilter) (*repository.PaginatedProjects, error) {
	repoFilter := repository.ProjectFilter{
		Keyword:    filter.Keyword,
		IsArchived: filter.IsArchived,
		OwnerName:  filter.OwnerName,
		UserID:     filter.UserID,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
	}
	return s.projectRepo.ListProjects(repoFilter)
}

// GetProjectWithTaskCounts 获取项目详情及任务统计
func (s *ProjectService) GetProjectWithTaskCounts(projectID uint) (*models.Project, map[string]int64, error) {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, nil, err
	}

	counts, err := s.projectRepo.GetTaskCounts(projectID)
	if err != nil {
		return project, nil, nil // 返回项目但忽略统计错误
	}

	return project, counts, nil
}

// CreateAuditLog 记录审计日志
func (s *ProjectService) CreateAuditLog(userID uint, eventType, details, ipAddress string) {
	auditLog := models.AuditLog{
		UserID:    userID,
		EventType: eventType,
		Details:   details,
		IPAddress: ipAddress,
	}
	s.projectRepo.CreateAuditLog(&auditLog)
}

// CreateProject 创建项目
func (s *ProjectService) CreateProject(name, description string, creatorID uint) (*models.Project, error) {
	// 创建项目
	project := &models.Project{
		Name:        name,
		Description: description,
		IsArchived:  false,
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}

	// 自动为创建者分配 Owner 角色
	ownerRole, err := s.findProjectOwnerRole()
	if err != nil {
		// 回滚项目创建
		s.projectRepo.Delete(project.ID)
		return nil, errors.New("系统错误：无法找到项目所有者角色")
	}

	if err := s.membershipRepo.AssignRole(creatorID, project.ID, ownerRole.ID); err != nil {
		// 回滚项目创建
		s.projectRepo.Delete(project.ID)
		return nil, err
	}

	return project, nil
}

// UpdateProject 更新项目
func (s *ProjectService) UpdateProject(projectID uint, name, description string) (*models.Project, error) {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("项目不存在")
	}

	project.Name = name
	project.Description = description

	if err := s.projectRepo.Update(project); err != nil {
		return nil, err
	}

	return project, nil
}

// DeleteProject 删除项目（级联删除任务和评论）
func (s *ProjectService) DeleteProject(projectID uint) error {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return errors.New("项目不存在")
	}

	// 使用事务执行级联删除
	return repository.GetDB().Transaction(func(tx *gorm.DB) error {
		// 删除评论（通过任务）
		if err := tx.Where("task_id IN (SELECT id FROM tasks WHERE project_id = ?)", projectID).Delete(&models.Comment{}).Error; err != nil {
			return fmt.Errorf("删除评论失败: %v", err)
		}

		// 删除任务
		if err := tx.Where("project_id = ?", projectID).Delete(&models.Task{}).Error; err != nil {
			return fmt.Errorf("删除任务失败: %v", err)
		}

		// 删除项目成员资格
		if err := tx.Where("project_id = ?", projectID).Delete(&models.ProjectMembership{}).Error; err != nil {
			return fmt.Errorf("删除成员资格失败: %v", err)
		}

		// 删除项目
		if err := tx.Delete(project).Error; err != nil {
			return fmt.Errorf("删除项目失败: %v", err)
		}

		return nil
	})
}

// ArchiveProject 归档/取消归档项目
func (s *ProjectService) ArchiveProject(projectID uint, archive bool) (*models.Project, error) {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("项目不存在")
	}

	project.IsArchived = archive

	if err := s.projectRepo.Update(project); err != nil {
		return nil, err
	}

	return project, nil
}

// GetProject 获取项目详情
func (s *ProjectService) GetProject(projectID uint) (*models.Project, error) {
	return s.projectRepo.FindByID(projectID)
}

// GetUserProjects 获取用户参与的项目列表
func (s *ProjectService) GetUserProjects(userID uint) ([]models.Project, error) {
	var projects []models.Project
	err := repository.GetDB().
		Table("projects").
		Select("projects.*").
		Joins("JOIN project_memberships ON projects.id = project_memberships.project_id").
		Where("project_memberships.user_id = ? AND projects.is_archived = ?", userID, false).
		Find(&projects).Error
	return projects, err
}

// findProjectOwnerRole 查找项目所有者角色
func (s *ProjectService) findProjectOwnerRole() (*models.Role, error) {
	var role models.Role
	err := repository.GetDB().Where("name = ? AND scope = ?", models.RoleOwner, "project").First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

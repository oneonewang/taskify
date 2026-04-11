package repository

import (
	"github.com/taskify/backend/internal/models"
	"gorm.io/gorm"
)

// ProjectRepository 项目数据访问层
type ProjectRepository struct {
	db *gorm.DB
}

// NewProjectRepository 创建项目仓库
func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{db: GetDB()}
}

// Create 创建项目
func (r *ProjectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

// Update 更新项目
func (r *ProjectRepository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

// Delete 删除项目
func (r *ProjectRepository) Delete(id uint) error {
	return r.db.Delete(&models.Project{}, id).Error
}

// FindByID 根据ID查找项目
func (r *ProjectRepository) FindByID(id uint) (*models.Project, error) {
	var project models.Project
	err := r.db.First(&project, id).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// GetAll 获取所有未归档的项目
func (r *ProjectRepository) GetAll() ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Where("is_archived = ?", false).Order("created_at DESC").Find(&projects).Error
	return projects, err
}

// GetByUser 获取用户参与的项目
func (r *ProjectRepository) GetByUser(userID uint) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.
		Joins("JOIN project_memberships ON projects.id = project_memberships.project_id").
		Where("project_memberships.user_id = ?", userID).
		Find(&projects).Error
	return projects, err
}

// GetTaskCounts 获取项目的任务统计
func (r *ProjectRepository) GetTaskCounts(projectID uint) (map[string]int64, error) {
	counts := make(map[string]int64)
	var todoCount, inProgressCount, doneCount int64

	r.db.Model(&models.Task{}).Where("project_id = ? AND status = ?", projectID, models.StatusTodo).Count(&todoCount)
	r.db.Model(&models.Task{}).Where("project_id = ? AND status = ?", projectID, models.StatusInProgress).Count(&inProgressCount)
	r.db.Model(&models.Task{}).Where("project_id = ? AND status = ?", projectID, models.StatusDone).Count(&doneCount)

	counts["todo"] = todoCount
	counts["in_progress"] = inProgressCount
	counts["done"] = doneCount

	return counts, nil
}

// CreateAuditLog 创建审计日志
func (r *ProjectRepository) CreateAuditLog(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

// ProjectFilter 项目列表过滤条件
type ProjectFilter struct {
	Keyword    string // 项目名称关键字（模糊搜索）
	IsArchived *bool  // 是否归档（nil 表示不过滤）
	OwnerName  string // 所有者名称（模糊搜索）
	UserID     uint   // 用户ID（大于0时过滤该用户参与的项目）
	Page       int    // 页码（从1开始）
	PageSize   int    // 每页数量
}

// PaginatedProjects 分页项目结果
type PaginatedProjects struct {
	Projects   []models.Project `json:"projects"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

// ListProjects 分页查询项目列表
func (r *ProjectRepository) ListProjects(filter ProjectFilter) (*PaginatedProjects, error) {
	var projects []models.Project
	var total int64

	// 默认分页参数
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	// 构建基础查询
	query := r.db.Model(&models.Project{})

	// 应用过滤条件
	if filter.Keyword != "" {
		query = query.Where("name LIKE ?", "%"+filter.Keyword+"%")
	}
	if filter.IsArchived != nil {
		query = query.Where("is_archived = ?", *filter.IsArchived)
	}

	// 如果指定了所有者名称，需要JOIN用户表
	if filter.OwnerName != "" {
		// 查找所有者角色的ID
		var ownerRole models.Role
		if err := r.db.Where("name = ? AND scope = ?", models.RoleOwner, "project").First(&ownerRole).Error; err == nil {
			// 通过项目成员表关联用户表，查找所有者名称匹配的项目
			query = query.
				Joins("JOIN project_memberships ON projects.id = project_memberships.project_id").
				Joins("JOIN users ON project_memberships.user_id = users.id").
				Where("project_memberships.role_id = ?", ownerRole.ID).
				Where("users.display_name LIKE ?", "%"+filter.OwnerName+"%")
		}
	}

	// 如果指定了用户ID，过滤该用户参与的项目
	if filter.UserID > 0 {
		query = query.
			Joins("JOIN project_memberships ON projects.id = project_memberships.project_id").
			Where("project_memberships.user_id = ?", filter.UserID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.PageSize).Find(&projects).Error; err != nil {
		return nil, err
	}

	// 计算总页数
	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &PaginatedProjects{
		Projects:   projects,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

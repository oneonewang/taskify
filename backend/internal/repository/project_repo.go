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
	var todoCount, inProgressCount, reviewCount, doneCount int64

	r.db.Model(&models.Task{}).Where("project_id = ? AND status = ?", projectID, models.StatusTodo).Count(&todoCount)
	r.db.Model(&models.Task{}).Where("project_id = ? AND status = ?", projectID, models.StatusInProgress).Count(&inProgressCount)
	r.db.Model(&models.Task{}).Where("project_id = ? AND status = ?", projectID, models.StatusReview).Count(&reviewCount)
	r.db.Model(&models.Task{}).Where("project_id = ? AND status = ?", projectID, models.StatusDone).Count(&doneCount)

	counts["todo"] = todoCount
	counts["in_progress"] = inProgressCount
	counts["review"] = reviewCount
	counts["done"] = doneCount

	return counts, nil
}

// CreateAuditLog 创建审计日志
func (r *ProjectRepository) CreateAuditLog(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

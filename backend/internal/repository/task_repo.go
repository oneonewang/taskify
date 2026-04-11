package repository

import (
	"github.com/taskify/backend/internal/models"
	"gorm.io/gorm"
)

// TaskRepository 任务数据访问层
type TaskRepository struct {
	db *gorm.DB
}

// NewTaskRepository 创建任务仓库
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{db: GetDB()}
}

// Create 创建任务
func (r *TaskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

// Update 更新任务
func (r *TaskRepository) Update(task *models.Task) error {
	// 只更新任务本身字段，不保存关联的 Assignee 关系
	return r.db.Model(task).Select("title", "description", "status", "position", "assignee_id", "updated_at").Updates(task).Error
}

// Delete 删除任务
func (r *TaskRepository) Delete(id uint) error {
	return r.db.Delete(&models.Task{}, id).Error
}

// FindByID 根据ID查找任务
func (r *TaskRepository) FindByID(id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.Preload("Assignee").First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetByProject 获取项目的所有任务
func (r *TaskRepository) GetByProject(projectID uint, status string) ([]models.Task, error) {
	var tasks []models.Task
	query := r.db.Preload("Assignee").Where("project_id = ?", projectID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("position ASC").Find(&tasks).Error
	return tasks, err
}

// GetMaxPosition 获取某状态的最大position
func (r *TaskRepository) GetMaxPosition(projectID uint, status models.TaskStatus) (int, error) {
	var maxPosition int
	err := r.db.Model(&models.Task{}).
		Where("project_id = ? AND status = ?", projectID, status).
		Select("COALESCE(MAX(position), 0)").
		Scan(&maxPosition).Error
	return maxPosition, err
}

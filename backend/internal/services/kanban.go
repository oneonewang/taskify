package services

import (
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
	"gorm.io/gorm"
)

// KanbanService 看板业务逻辑服务
// 注意：此服务当前未被使用，任务操作直接在 handlers 中实现
// 如需重构为标准三层架构，可调用此服务并由 handler 层处理事件广播
type KanbanService struct {
	db *gorm.DB
}

// NewKanbanService 创建看板服务实例
func NewKanbanService() *KanbanService {
	return &KanbanService{
		db: repository.GetDB(),
	}
}

// GetProjectTasks 获取项目的所有任务(按状态分组)
func (s *KanbanService) GetProjectTasks(projectID uint) (map[string][]models.TaskResponse, error) {
	var tasks []models.Task
	result := s.db.Preload("Assignee").Where("project_id = ?", projectID).Order("position ASC").Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}

	// 按状态分组
	tasksByStatus := map[string][]models.TaskResponse{
		"todo":        {},
		"in_progress": {},
		"done":        {},
	}

	for _, task := range tasks {
		tasksByStatus[string(task.Status)] = append(tasksByStatus[string(task.Status)], task.ToResponse())
	}

	return tasksByStatus, nil
}

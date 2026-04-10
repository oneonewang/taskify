package services

import (
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
	"github.com/taskify/backend/pkg/broadcaster"
	"gorm.io/gorm"
)

// KanbanService 看板业务逻辑服务
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
		"review":      {},
		"done":        {},
	}

	for _, task := range tasks {
		tasksByStatus[string(task.Status)] = append(tasksByStatus[string(task.Status)], task.ToResponse())
	}

	return tasksByStatus, nil
}

// MoveTask 移动任务到新状态/位置
func (s *KanbanService) MoveTask(taskID uint, newStatus models.TaskStatus, newPosition int) error {
	var task models.Task
	if result := s.db.First(&task, taskID); result.Error != nil {
		return result.Error
	}

	oldStatus := task.Status

	// 更新任务状态和位置
	task.Status = newStatus
	task.Position = newPosition
	if err := s.db.Save(&task).Error; err != nil {
		return err
	}

	// 广播任务移动事件
	broadcaster.Broadcast("task_moved", map[string]interface{}{
		"id":       task.ID,
		"from":     oldStatus,
		"to":       newStatus,
		"position": newPosition,
	})

	return nil
}

// CreateTask 创建任务并广播事件
func (s *KanbanService) CreateTask(projectID uint, title, description string, assigneeID uint) (*models.Task, error) {
	// 获取该状态任务的最大position
	var maxPosition int
	s.db.Model(&models.Task{}).
		Where("project_id = ? AND status = ?", projectID, models.StatusTodo).
		Select("COALESCE(MAX(position), -1)").
		Scan(&maxPosition)

	task := &models.Task{
		Title:       title,
		Description: description,
		Status:      models.StatusTodo,
		Position:    maxPosition + 1,
		AssigneeID:  assigneeID,
		ProjectID:   projectID,
	}

	if err := s.db.Create(task).Error; err != nil {
		return nil, err
	}

	// 加载关联数据
	s.db.Preload("Assignee").First(task, task.ID)

	// 广播任务创建事件
	broadcaster.Broadcast("task_created", task.ToResponse())

	return task, nil
}

// UpdateTask 更新任务并广播事件
func (s *KanbanService) UpdateTask(taskID uint, title, description string, assigneeID uint) (*models.Task, error) {
	var task models.Task
	if err := s.db.First(&task, taskID).Error; err != nil {
		return nil, err
	}

	if title != "" {
		task.Title = title
	}
	if description != "" {
		task.Description = description
	}
	if assigneeID != 0 {
		task.AssigneeID = assigneeID
	}

	if err := s.db.Save(&task).Error; err != nil {
		return nil, err
	}

	// 加载关联数据
	s.db.Preload("Assignee").First(&task, taskID)

	// 广播任务更新事件
	broadcaster.Broadcast("task_updated", task.ToResponse())

	return &task, nil
}

// DeleteTask 删除任务并广播事件
func (s *KanbanService) DeleteTask(taskID uint) error {
	var task models.Task
	if err := s.db.First(&task, taskID).Error; err != nil {
		return err
	}

	if err := s.db.Delete(&task).Error; err != nil {
		return err
	}

	// 广播任务删除事件
	broadcaster.Broadcast("task_deleted", map[string]interface{}{"id": taskID})

	return nil
}

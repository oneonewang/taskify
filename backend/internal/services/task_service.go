package services

import (
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
)

// TaskService 任务服务
type TaskService struct {
	taskRepo       *repository.TaskRepository
	membershipRepo *repository.MembershipRepository
}

// NewTaskService 创建任务服务
func NewTaskService() *TaskService {
	return &TaskService{
		taskRepo:       repository.NewTaskRepository(),
		membershipRepo: repository.NewMembershipRepository(),
	}
}

// CheckProjectAccess 检查用户是否有项目访问权限
func (s *TaskService) CheckProjectAccess(userID, projectID uint) (bool, error) {
	// 系统管理员有权限
	isAdmin, err := s.membershipRepo.IsAdmin(userID)
	if err == nil && isAdmin {
		return true, nil
	}
	// 检查是否是项目成员
	return s.membershipRepo.IsMember(userID, projectID)
}

// GetTasksByProject 获取项目的所有任务
func (s *TaskService) GetTasksByProject(projectID uint, status string) ([]models.Task, error) {
	return s.taskRepo.GetByProject(projectID, status)
}

// GetTask 获取任务详情
func (s *TaskService) GetTask(taskID uint) (*models.Task, error) {
	return s.taskRepo.FindByID(taskID)
}

// CreateTask 创建任务
func (s *TaskService) CreateTask(projectID uint, title, description string, assigneeID uint) (*models.Task, error) {
	// 获取该状态任务的最大position
	maxPosition, err := s.taskRepo.GetMaxPosition(projectID, models.StatusTodo)
	if err != nil {
		maxPosition = 0
	}

	task := &models.Task{
		Title:       title,
		Description: description,
		Status:      models.StatusTodo,
		Position:    maxPosition + 1,
		AssigneeID:  assigneeID,
		ProjectID:   projectID,
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}

	return s.taskRepo.FindByID(task.ID)
}

// UpdateTask 更新任务
func (s *TaskService) UpdateTask(taskID uint, title *string, description *string, assigneeID *uint) (*models.Task, error) {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, err
	}

	if title != nil {
		task.Title = *title
	}
	if description != nil {
		task.Description = *description
	}
	if assigneeID != nil {
		task.AssigneeID = *assigneeID
	}

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	return s.taskRepo.FindByID(task.ID)
}

// UpdateTaskStatus 更新任务状态
func (s *TaskService) UpdateTaskStatus(taskID uint, status models.TaskStatus, position int) (*models.Task, error) {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, err
	}

	task.Status = status
	task.Position = position

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

// DeleteTask 删除任务
func (s *TaskService) DeleteTask(taskID uint) error {
	return s.taskRepo.Delete(taskID)
}

package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/internal/middleware"
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/services"
	"github.com/taskify/backend/pkg/response"
)

// TaskHandler 任务处理器
type TaskHandler struct {
	taskService *services.TaskService
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		taskService: services.NewTaskService(),
	}
}

// GetTasks 获取项目的所有任务
func (h *TaskHandler) GetTasks(c *gin.Context) {
	projectID := c.Param("id")
	status := c.Query("status")

	tasks, err := h.taskService.GetTasksByProject(uint(parseUint(projectID)), status)
	if err != nil {
		response.InternalError(c, "获取任务列表失败")
		return
	}

	// 转换为响应格式
	taskResponses := make([]models.TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = task.ToResponse()
	}

	response.Success(c, taskResponses)
}

// CreateTask 创建新任务
func (h *TaskHandler) CreateTask(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID := uint(parseUint(projectIDStr))

	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 检查项目访问权限
	hasAccess, err := h.taskService.CheckProjectAccess(userID, projectID)
	if err != nil || !hasAccess {
		response.Forbidden(c, "您不是该项目成员")
		return
	}

	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	task, err := h.taskService.CreateTask(projectID, req.Title, req.Description, req.AssigneeID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, task.ToResponse())
}

// UpdateTask 更新任务
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	taskID := uint(parseUint(idStr))

	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 获取任务信息用于权限检查
	task, err := h.taskService.GetTask(taskID)
	if err != nil {
		response.NotFound(c, "任务不存在")
		return
	}

	// 检查项目访问权限
	hasAccess, err := h.taskService.CheckProjectAccess(userID, task.ProjectID)
	if err != nil || !hasAccess {
		response.Forbidden(c, "您不是该项目成员")
		return
	}

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	updatedTask, err := h.taskService.UpdateTask(taskID, req.Title, req.Description, req.AssigneeID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, updatedTask.ToResponse())
}

// UpdateTaskStatus 更新任务状态(看板拖放)
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	idStr := c.Param("id")
	taskID := uint(parseUint(idStr))

	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 获取任务信息用于权限检查
	task, err := h.taskService.GetTask(taskID)
	if err != nil {
		response.NotFound(c, "任务不存在")
		return
	}

	// 检查项目访问权限
	hasAccess, err := h.taskService.CheckProjectAccess(userID, task.ProjectID)
	if err != nil || !hasAccess {
		response.Forbidden(c, "您不是该项目成员")
		return
	}

	var req models.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 验证状态值
	validStatuses := []models.TaskStatus{models.StatusTodo, models.StatusInProgress, models.StatusDone}
	isValid := false
	for _, s := range validStatuses {
		if req.Status == s {
			isValid = true
			break
		}
	}
	if !isValid {
		response.BadRequest(c, "无效的状态值")
		return
	}

	updatedTask, err := h.taskService.UpdateTaskStatus(taskID, req.Status, req.Position)
	if err != nil {
		response.InternalError(c, "更新任务状态失败")
		return
	}

	response.Success(c, gin.H{
		"id":         updatedTask.ID,
		"status":     updatedTask.Status,
		"position":   updatedTask.Position,
		"updated_at": updatedTask.UpdatedAt,
	})
}

// DeleteTask 删除任务
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	taskID := uint(parseUint(idStr))

	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 获取任务信息用于权限检查
	task, err := h.taskService.GetTask(taskID)
	if err != nil {
		response.NotFound(c, "任务不存在")
		return
	}

	// 检查项目访问权限
	hasAccess, err := h.taskService.CheckProjectAccess(userID, task.ProjectID)
	if err != nil || !hasAccess {
		response.Forbidden(c, "您不是该项目成员")
		return
	}

	if err := h.taskService.DeleteTask(taskID); err != nil {
		response.InternalError(c, "删除任务失败")
		return
	}

	response.OK(c, "任务已删除")
}

// Helper function
func parseUint(id string) uint64 {
	uintID, _ := strconv.ParseUint(id, 10, 32)
	return uintID
}

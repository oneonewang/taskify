package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
	"github.com/taskify/backend/pkg/response"
)

// GetTasks 获取项目的所有任务
func GetTasks(c *gin.Context) {
	projectID := c.Param("id")
	status := c.Query("status")

	db := repository.GetDB().Preload("Assignee").Where("project_id = ?", projectID)
	if status != "" {
		db = db.Where("status = ?", status)
	}

	var tasks []models.Task
	result := db.Order("position ASC").Find(&tasks)
	if result.Error != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, "获取任务列表失败")
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
func CreateTask(c *gin.Context) {
	projectID := c.Param("id")

	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "请求参数错误: "+err.Error())
		return
	}

	// 验证项目存在
	var project models.Project
	if result := repository.GetDB().First(&project, projectID); result.Error != nil {
		response.Error(c, http.StatusNotFound, response.CodeNotFound, "项目不存在")
		return
	}

	// 验证用户存在
	var assignee models.User
	if result := repository.GetDB().First(&assignee, req.AssigneeID); result.Error != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "指定的负责人不存在")
		return
	}

	// 获取该状态任务的最大position
	var maxPosition int
	repository.GetDB().Model(&models.Task{}).
		Where("project_id = ? AND status = ?", projectID, models.StatusTodo).
		Select("COALESCE(MAX(position), -1)").
		Scan(&maxPosition)

	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      models.StatusTodo,
		Position:    maxPosition + 1,
		AssigneeID:  req.AssigneeID,
		ProjectID:   uint(projectIDUint(projectID)),
	}

	result := repository.GetDB().Create(&task)
	if result.Error != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternalError, "创建任务失败")
		return
	}

	// 重新加载关联数据
	repository.GetDB().Preload("Assignee").First(&task, task.ID)

	response.Success(c, task.ToResponse())
}

// UpdateTask 更新任务
func UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "请求参数错误: "+err.Error())
		return
	}

	var task models.Task
	result := repository.GetDB().Preload("Assignee").First(&task, id)
	if result.Error != nil {
		response.Error(c, http.StatusNotFound, response.CodeNotFound, "任务不存在")
		return
	}

	// 更新字段
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.AssigneeID != nil {
		// 验证新负责人存在
		var assignee models.User
		if result := repository.GetDB().First(&assignee, *req.AssigneeID); result.Error != nil {
			response.Error(c, http.StatusBadRequest, response.CodeValidationError, "指定的负责人不存在")
			return
		}
		task.AssigneeID = *req.AssigneeID
		task.Assignee = &assignee
	}

	repository.GetDB().Save(&task)

	response.Success(c, task.ToResponse())
}

// UpdateTaskStatus 更新任务状态(看板拖放)
func UpdateTaskStatus(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "请求参数错误: "+err.Error())
		return
	}

	// 验证状态值
	validStatuses := []models.TaskStatus{models.StatusTodo, models.StatusInProgress, models.StatusReview, models.StatusDone}
	isValid := false
	for _, s := range validStatuses {
		if req.Status == s {
			isValid = true
			break
		}
	}
	if !isValid {
		response.Error(c, http.StatusBadRequest, response.CodeValidationError, "无效的状态值")
		return
	}

	var task models.Task
	result := repository.GetDB().First(&task, id)
	if result.Error != nil {
		response.Error(c, http.StatusNotFound, response.CodeNotFound, "任务不存在")
		return
	}

	task.Status = req.Status
	task.Position = req.Position
	repository.GetDB().Save(&task)

	response.Success(c, gin.H{
		"id":         task.ID,
		"status":     task.Status,
		"position":   task.Position,
		"updated_at": task.UpdatedAt,
	})
}

// DeleteTask 删除任务
func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	var task models.Task
	result := repository.GetDB().First(&task, id)
	if result.Error != nil {
		response.Error(c, http.StatusNotFound, response.CodeNotFound, "任务不存在")
		return
	}

	repository.GetDB().Delete(&task)

	response.SuccessWithMessage(c, nil, "任务已删除")
}

// Helper function
func projectIDUint(id string) uint {
	uintID, _ := strconv.ParseUint(id, 10, 32)
	return uint(uintID)
}

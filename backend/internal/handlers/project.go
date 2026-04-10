package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
	"github.com/taskify/backend/pkg/response"
)

// GetUsers 获取所有用户
func GetUsers(c *gin.Context) {
	var users []models.User
	result := repository.GetDB().Find(&users)
	if result.Error != nil {
		response.InternalError(c, "获取用户列表失败")
		return
	}

	// 转换为响应格式
	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	response.Success(c, userResponses)
}

// GetProjects 获取所有项目
func GetProjects(c *gin.Context) {
	var projects []models.Project
	result := repository.GetDB().Find(&projects)
	if result.Error != nil {
		response.InternalError(c, "获取项目列表失败")
		return
	}

	// 转换为响应格式
	projectResponses := make([]models.ProjectResponse, len(projects))
	for i, project := range projects {
		projectResponses[i] = project.ToResponse()
	}

	response.Success(c, projectResponses)
}

// GetProject 获取项目详情
func GetProject(c *gin.Context) {
	id := c.Param("id")

	var project models.Project
	result := repository.GetDB().First(&project, id)
	if result.Error != nil {
		response.NotFound(c, "项目不存在")
		return
	}

	// 统计任务数量
	var todoCount, inProgressCount, reviewCount, doneCount int64
	db := repository.GetDB()
	db.Model(&models.Task{}).Where("project_id = ? AND status = ?", id, models.StatusTodo).Count(&todoCount)
	db.Model(&models.Task{}).Where("project_id = ? AND status = ?", id, models.StatusInProgress).Count(&inProgressCount)
	db.Model(&models.Task{}).Where("project_id = ? AND status = ?", id, models.StatusReview).Count(&reviewCount)
	db.Model(&models.Task{}).Where("project_id = ? AND status = ?", id, models.StatusDone).Count(&doneCount)

	// 构建响应
	detailResp := struct {
		models.ProjectResponse
		TaskCounts struct {
			Todo       int64 `json:"todo"`
			InProgress int64 `json:"in_progress"`
			Review     int64 `json:"review"`
			Done       int64 `json:"done"`
		} `json:"task_counts"`
	}{
		ProjectResponse: project.ToResponse(),
	}
	detailResp.TaskCounts.Todo = todoCount
	detailResp.TaskCounts.InProgress = inProgressCount
	detailResp.TaskCounts.Review = reviewCount
	detailResp.TaskCounts.Done = doneCount

	response.Success(c, detailResp)
}
package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/internal/middleware"
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
	"github.com/taskify/backend/internal/services"
	"github.com/taskify/backend/pkg/response"
)

// ProjectHandler 项目处理器
type ProjectHandler struct {
	projectService *services.ProjectService
}

// NewProjectHandler 创建项目处理器
func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{
		projectService: services.NewProjectService(),
	}
}

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
}

// GetUsers 获取所有用户（管理员用）
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
	result := repository.GetDB().Where("is_archived = ?", false).Order("created_at DESC").Find(&projects)
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

// CreateProject 创建项目
// POST /api/projects
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	project, err := h.projectService.CreateProject(req.Name, req.Description, userID)
	if err != nil {
		response.InternalError(c, "创建项目失败")
		return
	}

	// 记录审计日志
	auditLog := models.AuditLog{
		UserID:    userID,
		EventType: models.EventProjectCreated,
		Details:   "创建项目: " + project.Name,
		IPAddress: middleware.GetClientIP(c),
	}
	repository.GetDB().Create(&auditLog)

	response.Created(c, "项目创建成功", project.ToResponse())
}

// UpdateProject 更新项目
// PUT /api/projects/:id
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	project, err := h.projectService.UpdateProject(uint(id), req.Name, req.Description)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, project.ToResponse())
}

// DeleteProject 删除项目
// DELETE /api/projects/:id
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	// 获取项目信息用于日志
	project, _ := h.projectService.GetProject(uint(id))
	projectName := ""
	if project != nil {
		projectName = project.Name
	}

	if err := h.projectService.DeleteProject(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 记录审计日志
	auditLog := models.AuditLog{
		UserID:    userID,
		EventType: models.EventProjectDeleted,
		Details:   "删除项目: " + projectName,
		IPAddress: middleware.GetClientIP(c),
	}
	repository.GetDB().Create(&auditLog)

	response.OK(c, "项目删除成功")
}

// ArchiveProject 归档项目
// POST /api/projects/:id/archive
func (h *ProjectHandler) ArchiveProject(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	// 获取当前项目状态
	project, err := h.projectService.GetProject(uint(id))
	if err != nil {
		response.NotFound(c, "项目不存在")
		return
	}

	// 切换归档状态
	newArchiveStatus := !project.IsArchived
	updatedProject, err := h.projectService.ArchiveProject(uint(id), newArchiveStatus)
	if err != nil {
		response.InternalError(c, "操作失败")
		return
	}

	// 记录审计日志
	action := "归档"
	if !newArchiveStatus {
		action = "取消归档"
	}
	auditLog := models.AuditLog{
		UserID:    userID,
		EventType: models.EventProjectArchived,
		Details:   action + "项目: " + updatedProject.Name,
		IPAddress: middleware.GetClientIP(c),
	}
	repository.GetDB().Create(&auditLog)

	response.Success(c, updatedProject.ToResponse())
}

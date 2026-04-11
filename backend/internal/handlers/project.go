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
func (h *ProjectHandler) GetUsers(c *gin.Context) {
	users, err := h.projectService.GetAllUsers()
	if err != nil {
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

// GetProjects 获取所有项目（支持分页和过滤）
// 支持 my=true 参数获取当前用户参与的项目
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 解析过滤参数
	keyword := c.Query("keyword")
	ownerName := c.Query("owner_name")
	isArchivedStr := c.Query("is_archived")

	var isArchived *bool
	if isArchivedStr != "" {
		archived := isArchivedStr == "true"
		isArchived = &archived
	}

	// 解析 my 参数
	myStr := c.Query("my")
	var userID uint
	if myStr == "true" {
		uid, ok := middleware.GetUserID(c)
		if !ok {
			response.Unauthorized(c, "请先登录")
			return
		}
		userID = uid
	}

	// 调用服务获取分页结果
	result, err := h.projectService.ListProjects(services.ProjectListFilter{
		Keyword:    keyword,
		IsArchived: isArchived,
		OwnerName:  ownerName,
		UserID:     userID,
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		response.InternalError(c, "获取项目列表失败")
		return
	}

	// 获取每个项目的所有者信息
	db := repository.GetDB()
	projectResponses := make([]models.ProjectResponse, len(result.Projects))
	for i, project := range result.Projects {
		resp := project.ToResponse()
		// 查询项目所有者
		var membership models.ProjectMembership
		if err := db.
			Joins("JOIN roles ON project_memberships.role_id = roles.id").
			Where("project_memberships.project_id = ? AND roles.name = ? AND roles.scope = ?", project.ID, models.RoleOwner, "project").
			First(&membership).Error; err == nil {
			var owner models.User
			if err := db.First(&owner, membership.UserID).Error; err == nil {
				resp.OwnerID = owner.ID
				resp.OwnerName = owner.DisplayName
				resp.OwnerAvatar = owner.AvatarURL
			}
		}
		projectResponses[i] = resp
	}

	// 返回分页结果
	response.Success(c, gin.H{
		"projects":    projectResponses,
		"total":       result.Total,
		"page":        result.Page,
		"page_size":   result.PageSize,
		"total_pages": result.TotalPages,
	})
}

// GetMyProjects 获取当前用户参与的项目
func (h *ProjectHandler) GetMyProjects(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	projects, err := h.projectService.GetUserProjects(userID)
	if err != nil {
		response.InternalError(c, "获取项目列表失败")
		return
	}

	projectResponses := make([]models.ProjectResponse, len(projects))
	for i, project := range projects {
		projectResponses[i] = project.ToResponse()
	}

	response.Success(c, projectResponses)
}

// GetProject 获取项目详情
func (h *ProjectHandler) GetProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的项目ID")
		return
	}

	project, counts, err := h.projectService.GetProjectWithTaskCounts(uint(id))
	if err != nil {
		response.NotFound(c, "项目不存在")
		return
	}

	// 构建响应
	db := repository.GetDB()
	resp := project.ToResponse()
	// 查询项目所有者
	var membership models.ProjectMembership
	if err := db.
		Joins("JOIN roles ON project_memberships.role_id = roles.id").
		Where("project_memberships.project_id = ? AND roles.name = ? AND roles.scope = ?", project.ID, models.RoleOwner, "project").
		First(&membership).Error; err == nil {
		var owner models.User
		if err := db.First(&owner, membership.UserID).Error; err == nil {
			resp.OwnerID = owner.ID
			resp.OwnerName = owner.DisplayName
			resp.OwnerAvatar = owner.AvatarURL
		}
	}

	detailResp := struct {
		models.ProjectResponse
		TaskCounts map[string]int64 `json:"task_counts"`
	}{
		ProjectResponse: resp,
		TaskCounts:      counts,
	}

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
	h.projectService.CreateAuditLog(userID, models.EventProjectCreated, "创建项目: "+project.Name, middleware.GetClientIP(c))

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
	h.projectService.CreateAuditLog(userID, models.EventProjectDeleted, "删除项目: "+projectName, middleware.GetClientIP(c))

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
	h.projectService.CreateAuditLog(userID, models.EventProjectArchived, action+"项目: "+updatedProject.Name, middleware.GetClientIP(c))

	response.Success(c, updatedProject.ToResponse())
}

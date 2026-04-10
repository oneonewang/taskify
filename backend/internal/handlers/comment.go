package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/internal/middleware"
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
	"github.com/taskify/backend/pkg/broadcaster"
	"github.com/taskify/backend/pkg/response"
)

// GetComments 获取任务的评论列表
func GetComments(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}

	// 验证任务存在并获取项目ID
	var task models.Task
	if result := repository.GetDB().First(&task, taskID); result.Error != nil {
		response.NotFound(c, "任务不存在")
		return
	}

	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 检查项目访问权限
	if !checkTaskProjectAccess(userID, task.ProjectID) {
		response.Forbidden(c, "您不是该项目成员")
		return
	}

	var comments []models.Comment
	result := repository.GetDB().Preload("User").Where("task_id = ?", taskID).Order("created_at ASC").Find(&comments)
	if result.Error != nil {
		response.InternalError(c, "获取评论列表失败")
		return
	}

	// 转换为响应格式
	commentResponses := make([]models.CommentResponse, len(comments))
	for i, comment := range comments {
		commentResponses[i] = comment.ToResponse()
	}

	response.Success(c, commentResponses)
}

// CreateComment 添加评论
func CreateComment(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}

	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 验证任务存在
	var task models.Task
	if result := repository.GetDB().First(&task, taskID); result.Error != nil {
		response.NotFound(c, "任务不存在")
		return
	}

	// 检查项目访问权限
	if !checkTaskProjectAccess(userID, task.ProjectID) {
		response.Forbidden(c, "您不是该项目成员")
		return
	}

	// 验证用户存在
	var user models.User
	if result := repository.GetDB().First(&user, userID); result.Error != nil {
		response.BadRequest(c, "用户不存在")
		return
	}

	comment := models.Comment{
		Content: req.Content,
		UserID:  userID,
		TaskID:  uint(taskID),
	}

	result := repository.GetDB().Create(&comment)
	if result.Error != nil {
		response.InternalError(c, "创建评论失败")
		return
	}

	// 重新加载关联数据
	repository.GetDB().Preload("User").First(&comment, comment.ID)

	// 广播评论添加事件
	broadcaster.Broadcast("comment_added", comment.ToResponse())

	response.Success(c, comment.ToResponse())
}

// UpdateComment 编辑评论
func UpdateComment(c *gin.Context) {
	taskIDStr := c.Param("id")
	commentIDStr := c.Param("cid")

	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req models.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	var comment models.Comment
	result := repository.GetDB().Preload("User").First(&comment, commentIDStr)
	if result.Error != nil {
		response.NotFound(c, "评论不存在")
		return
	}

	// 验证评论是否属于该任务
	if strconv.FormatUint(uint64(comment.TaskID), 10) != taskIDStr {
		response.BadRequest(c, "评论不属于该任务")
		return
	}

	// 验证任务存在并获取项目ID用于权限检查
	var task models.Task
	if result := repository.GetDB().First(&task, taskIDStr); result.Error != nil {
		response.NotFound(c, "任务不存在")
		return
	}

	// 检查项目访问权限
	if !checkTaskProjectAccess(userID, task.ProjectID) {
		response.Forbidden(c, "您不是该项目成员")
		return
	}

	// 验证是否是评论作者
	if comment.UserID != userID {
		response.Forbidden(c, "只能编辑自己的评论")
		return
	}

	comment.Content = req.Content
	repository.GetDB().Save(&comment)

	response.Success(c, comment.ToResponse())
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	taskIDStr := c.Param("id")
	commentIDStr := c.Param("cid")

	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	var comment models.Comment
	result := repository.GetDB().First(&comment, commentIDStr)
	if result.Error != nil {
		response.NotFound(c, "评论不存在")
		return
	}

	// 验证评论是否属于该任务
	if strconv.FormatUint(uint64(comment.TaskID), 10) != taskIDStr {
		response.BadRequest(c, "评论不属于该任务")
		return
	}

	// 验证任务存在并获取项目ID用于权限检查
	var task models.Task
	if result := repository.GetDB().First(&task, taskIDStr); result.Error != nil {
		response.NotFound(c, "任务不存在")
		return
	}

	// 检查项目访问权限
	if !checkTaskProjectAccess(userID, task.ProjectID) {
		response.Forbidden(c, "您不是该项目成员")
		return
	}

	// 验证是否是评论作者
	if comment.UserID != userID {
		response.Forbidden(c, "只能删除自己的评论")
		return
	}

	repository.GetDB().Delete(&comment)

	response.OK(c, "评论已删除")
}

// Helper function
func taskIDUint(id string) uint {
	uintID, _ := strconv.ParseUint(id, 10, 32)
	return uint(uintID)
}
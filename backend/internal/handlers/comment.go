package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
	"github.com/taskify/backend/pkg/broadcaster"
	"github.com/taskify/backend/pkg/response"
)

// GetComments 获取任务的评论列表
func GetComments(c *gin.Context) {
	taskID := c.Param("id")

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
	taskID := c.Param("id")
	userIDHeader := c.GetHeader("X-User-ID")

	if userIDHeader == "" {
		response.Unauthorized(c, "缺少用户认证信息")
		return
	}

	userID, err := strconv.ParseUint(userIDHeader, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
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

	// 验证用户存在
	var user models.User
	if result := repository.GetDB().First(&user, userID); result.Error != nil {
		response.BadRequest(c, "用户不存在")
		return
	}

	comment := models.Comment{
		Content: req.Content,
		UserID:  uint(userID),
		TaskID:  uint(taskIDUint(taskID)),
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
	taskID := c.Param("id")
	commentID := c.Param("cid")
	userIDHeader := c.GetHeader("X-User-ID")

	if userIDHeader == "" {
		response.Unauthorized(c, "缺少用户认证信息")
		return
	}

	userID, err := strconv.ParseUint(userIDHeader, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req models.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	var comment models.Comment
	result := repository.GetDB().Preload("User").First(&comment, commentID)
	if result.Error != nil {
		response.NotFound(c, "评论不存在")
		return
	}

	// 验证评论是否属于该任务
	if strconv.FormatUint(uint64(comment.TaskID), 10) != taskID {
		response.BadRequest(c, "评论不属于该任务")
		return
	}

	// 验证是否是评论作者
	if comment.UserID != uint(userID) {
		response.Forbidden(c, "只能编辑自己的评论")
		return
	}

	comment.Content = req.Content
	repository.GetDB().Save(&comment)

	response.Success(c, comment.ToResponse())
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	taskID := c.Param("id")
	commentID := c.Param("cid")
	userIDHeader := c.GetHeader("X-User-ID")

	if userIDHeader == "" {
		response.Unauthorized(c, "缺少用户认证信息")
		return
	}

	userID, err := strconv.ParseUint(userIDHeader, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var comment models.Comment
	result := repository.GetDB().First(&comment, commentID)
	if result.Error != nil {
		response.NotFound(c, "评论不存在")
		return
	}

	// 验证评论是否属于该任务
	if strconv.FormatUint(uint64(comment.TaskID), 10) != taskID {
		response.BadRequest(c, "评论不属于该任务")
		return
	}

	// 验证是否是评论作者
	if comment.UserID != uint(userID) {
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
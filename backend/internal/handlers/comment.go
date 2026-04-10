package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/internal/middleware"
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/services"
	"github.com/taskify/backend/pkg/broadcaster"
	"github.com/taskify/backend/pkg/response"
)

// CommentHandler 评论处理器
type CommentHandler struct {
	commentService *services.CommentService
}

// NewCommentHandler 创建评论处理器
func NewCommentHandler() *CommentHandler {
	return &CommentHandler{
		commentService: services.NewCommentService(),
	}
}

// GetComments 获取任务的评论列表
func (h *CommentHandler) GetComments(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID := uint(parseUint(taskIDStr))

	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 验证任务存在并获取项目ID
	task, err := h.commentService.GetTask(taskID)
	if err != nil {
		response.NotFound(c, "任务不存在")
		return
	}

	// 检查项目访问权限
	hasAccess, err := h.commentService.CheckProjectAccess(userID, task.ProjectID)
	if err != nil || !hasAccess {
		response.Forbidden(c, "您不是该项目成员")
		return
	}

	comments, err := h.commentService.GetComments(taskID)
	if err != nil {
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
func (h *CommentHandler) CreateComment(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID := uint(parseUint(taskIDStr))

	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 验证任务存在并获取项目ID
	task, err := h.commentService.GetTask(taskID)
	if err != nil {
		response.NotFound(c, "任务不存在")
		return
	}

	// 检查项目访问权限
	hasAccess, err := h.commentService.CheckProjectAccess(userID, task.ProjectID)
	if err != nil || !hasAccess {
		response.Forbidden(c, "您不是该项目成员")
		return
	}

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	comment, err := h.commentService.CreateComment(taskID, userID, req.Content)
	if err != nil {
		response.InternalError(c, "创建评论失败")
		return
	}

	// 广播评论添加事件
	broadcaster.Broadcast("comment_added", comment.ToResponse())

	response.Success(c, comment.ToResponse())
}

// UpdateComment 编辑评论
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	taskIDStr := c.Param("id")
	commentIDStr := c.Param("cid")
	commentID := uint(parseUint(commentIDStr))

	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 获取评论信息
	comment, err := h.commentService.GetComment(commentID)
	if err != nil {
		response.NotFound(c, "评论不存在")
		return
	}

	// 验证评论是否属于该任务
	if strconv.FormatUint(uint64(comment.TaskID), 10) != taskIDStr {
		response.BadRequest(c, "评论不属于该任务")
		return
	}

	// 验证任务存在并获取项目ID用于权限检查
	task, err := h.commentService.GetTask(comment.TaskID)
	if err != nil {
		response.NotFound(c, "任务不存在")
		return
	}

	// 检查项目访问权限
	hasAccess, err := h.commentService.CheckProjectAccess(userID, task.ProjectID)
	if err != nil || !hasAccess {
		response.Forbidden(c, "您不是该项目成员")
		return
	}

	// 验证是否是评论作者
	if comment.UserID != userID {
		response.Forbidden(c, "只能编辑自己的评论")
		return
	}

	var req models.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	updatedComment, err := h.commentService.UpdateComment(commentID, req.Content)
	if err != nil {
		response.InternalError(c, "更新评论失败")
		return
	}

	response.Success(c, updatedComment.ToResponse())
}

// DeleteComment 删除评论
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	taskIDStr := c.Param("id")
	commentIDStr := c.Param("cid")
	commentID := uint(parseUint(commentIDStr))

	// 获取当前用户ID
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 获取评论信息
	comment, err := h.commentService.GetComment(commentID)
	if err != nil {
		response.NotFound(c, "评论不存在")
		return
	}

	// 验证评论是否属于该任务
	if strconv.FormatUint(uint64(comment.TaskID), 10) != taskIDStr {
		response.BadRequest(c, "评论不属于该任务")
		return
	}

	// 验证任务存在并获取项目ID用于权限检查
	task, err := h.commentService.GetTask(comment.TaskID)
	if err != nil {
		response.NotFound(c, "任务不存在")
		return
	}

	// 检查项目访问权限
	hasAccess, err := h.commentService.CheckProjectAccess(userID, task.ProjectID)
	if err != nil || !hasAccess {
		response.Forbidden(c, "您不是该项目成员")
		return
	}

	// 验证是否是评论作者
	if comment.UserID != userID {
		response.Forbidden(c, "只能删除自己的评论")
		return
	}

	if err := h.commentService.DeleteComment(commentID); err != nil {
		response.InternalError(c, "删除评论失败")
		return
	}

	response.OK(c, "评论已删除")
}

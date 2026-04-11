package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
	"github.com/taskify/backend/pkg/response"
)

// AuditLogHandler 审计日志处理器
type AuditLogHandler struct {
	auditLogRepo *repository.AuditLogRepository
}

// NewAuditLogHandler 创建审计日志处理器
func NewAuditLogHandler() *AuditLogHandler {
	return &AuditLogHandler{
		auditLogRepo: repository.NewAuditLogRepository(),
	}
}

// AuditLogQueryRequest 审计日志查询请求
type AuditLogQueryRequest struct {
	UserID    string `form:"user_id"`     // 用户ID精确匹配
	EventType string `form:"event_type"`  // 事件类型筛选
	From      string `form:"from"`         // 开始时间
	To        string `form:"to"`           // 结束时间
	Page      int    `form:"page,default=1"`       // 页码
	PageSize  int    `form:"page_size,default=20"` // 每页数量
}

// ListAuditLogs 获取审计日志列表（支持筛选和分页）
// GET /api/admin/audit-logs
func (h *AuditLogHandler) ListAuditLogs(c *gin.Context) {
	var req AuditLogQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	filter := repository.AuditLogFilter{
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	if req.UserID != "" {
		userID, err := strconv.ParseUint(req.UserID, 10, 64)
		if err == nil {
			uid := uint(userID)
			filter.UserID = &uid
		}
	}

	if req.EventType != "" {
		filter.EventType = req.EventType
	}

	result, err := h.auditLogRepo.FindAuditLogs(filter)
	if err != nil {
		response.InternalError(c, "查询审计日志失败")
		return
	}

	// 转换为响应格式
	logs := make([]models.AuditLogResponse, len(result.Logs))
	for i, log := range result.Logs {
		logs[i] = log.ToResponse()
		// 填充用户邮箱
		if log.UserID > 0 {
			if user, err := h.auditLogRepo.GetUserEmail(log.UserID); err == nil {
				logs[i].UserEmail = user.Email
			}
		}
	}

	response.Success(c, gin.H{
		"logs":       logs,
		"total":      result.Total,
		"page":       result.Page,
		"page_size":  result.PageSize,
		"total_pages": result.TotalPages,
	})
}

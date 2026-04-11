package repository

import (
	"github.com/taskify/backend/internal/models"
)

// AuditLogRepository 审计日志仓库
type AuditLogRepository struct{}

// NewAuditLogRepository 创建审计日志仓库
func NewAuditLogRepository() *AuditLogRepository {
	return &AuditLogRepository{}
}

// AuditLogFilter 审计日志查询过滤条件
type AuditLogFilter struct {
	UserID    *uint  // 用户ID精确匹配
	EventType string // 事件类型筛选
	From      string // 开始时间
	To        string // 结束时间
	Page      int    // 页码
	PageSize  int    // 每页数量
}

// AuditLogResult 审计日志查询结果
type AuditLogResult struct {
	Logs       []models.AuditLog
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

// FindAuditLogs 查询审计日志（支持筛选和分页）
func (r *AuditLogRepository) FindAuditLogs(filter AuditLogFilter) (*AuditLogResult, error) {
	var logs []models.AuditLog
	var total int64

	query := DB.Model(&models.AuditLog{})

	// 应用筛选条件
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.EventType != "" {
		query = query.Where("event_type = ?", filter.EventType)
	}
	if filter.From != "" {
		query = query.Where("created_at >= ?", filter.From)
	}
	if filter.To != "" {
		query = query.Where("created_at <= ?", filter.To)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}
	offset := (filter.Page - 1) * filter.PageSize

	if err := query.Order("created_at DESC").Offset(offset).Limit(filter.PageSize).Find(&logs).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &AuditLogResult{
		Logs:       logs,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetUserEmail 根据用户ID获取邮箱
func (r *AuditLogRepository) GetUserEmail(userID uint) (*models.User, error) {
	var user models.User
	if err := DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建审计日志
func (r *AuditLogRepository) Create(log *models.AuditLog) error {
	return DB.Create(log).Error
}

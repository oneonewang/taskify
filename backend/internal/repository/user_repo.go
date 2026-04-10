package repository

import (
	"github.com/taskify/backend/internal/models"
	"gorm.io/gorm"
)

// UserRepository 用户数据访问层
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository() *UserRepository {
	return &UserRepository{db: GetDB()}
}

// FindByEmail 通过邮箱查找用户
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 通过ID查找用户
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// GetAll 获取所有用户
func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

// UserFilter 用户查询过滤条件
type UserFilter struct {
	Email      string // 邮箱精确匹配
	DisplayName string // 显示名称模糊匹配
	IsDisabled *bool  // 禁用状态筛选
	Page       int    // 页码（从1开始）
	PageSize   int    // 每页数量
}

// PaginatedUsers 分页用户结果
type PaginatedUsers struct {
	Users      []models.User `json:"users"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

// FindUsers 分页查询用户
func (r *UserRepository) FindUsers(filter UserFilter) (*PaginatedUsers, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})

	// 应用过滤条件
	if filter.Email != "" {
		query = query.Where("email = ?", filter.Email)
	}
	if filter.DisplayName != "" {
		query = query.Where("display_name LIKE ?", "%"+filter.DisplayName+"%")
	}
	if filter.IsDisabled != nil {
		query = query.Where("is_disabled = ?", *filter.IsDisabled)
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

	// 查询数据
	if err := query.Offset(offset).Limit(filter.PageSize).Order("id DESC").Find(&users).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &PaginatedUsers{
		Users:      users,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}
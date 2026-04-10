package repository

import (
	"github.com/taskify/backend/internal/models"
	"gorm.io/gorm"
)

// CommentRepository 评论数据访问层
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓库
func NewCommentRepository() *CommentRepository {
	return &CommentRepository{db: GetDB()}
}

// Create 创建评论
func (r *CommentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

// Update 更新评论
func (r *CommentRepository) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

// Delete 删除评论
func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Comment{}, id).Error
}

// FindByID 根据ID查找评论
func (r *CommentRepository) FindByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.Preload("User").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetByTask 获取任务的评论列表
func (r *CommentRepository) GetByTask(taskID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Preload("User").Where("task_id = ?", taskID).Order("created_at ASC").Find(&comments).Error
	return comments, err
}

// GetTask 获取任务（用于权限检查）
func (r *CommentRepository) GetTask(taskID uint) (*models.Task, error) {
	var task models.Task
	err := r.db.First(&task, taskID).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

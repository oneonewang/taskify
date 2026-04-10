package services

import (
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
)

// CommentService 评论服务
type CommentService struct {
	commentRepo    *repository.CommentRepository
	membershipRepo *repository.MembershipRepository
}

// NewCommentService 创建评论服务
func NewCommentService() *CommentService {
	return &CommentService{
		commentRepo:    repository.NewCommentRepository(),
		membershipRepo: repository.NewMembershipRepository(),
	}
}

// CheckProjectAccess 检查用户是否有项目访问权限
func (s *CommentService) CheckProjectAccess(userID, projectID uint) (bool, error) {
	// 系统管理员有权限
	isAdmin, err := s.membershipRepo.IsAdmin(userID)
	if err == nil && isAdmin {
		return true, nil
	}
	// 检查是否是项目成员
	return s.membershipRepo.IsMember(userID, projectID)
}

// GetComments 获取任务的评论列表
func (s *CommentService) GetComments(taskID uint) ([]models.Comment, error) {
	return s.commentRepo.GetByTask(taskID)
}

// GetComment 获取评论详情
func (s *CommentService) GetComment(commentID uint) (*models.Comment, error) {
	return s.commentRepo.FindByID(commentID)
}

// GetTask 获取任务（用于权限检查）
func (s *CommentService) GetTask(taskID uint) (*models.Task, error) {
	return s.commentRepo.GetTask(taskID)
}

// CreateComment 创建评论
func (s *CommentService) CreateComment(taskID, userID uint, content string) (*models.Comment, error) {
	comment := &models.Comment{
		Content: content,
		UserID:  userID,
		TaskID:  taskID,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	return s.commentRepo.FindByID(comment.ID)
}

// UpdateComment 更新评论
func (s *CommentService) UpdateComment(commentID uint, content string) (*models.Comment, error) {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return nil, err
	}

	comment.Content = content

	if err := s.commentRepo.Update(comment); err != nil {
		return nil, err
	}

	return s.commentRepo.FindByID(comment.ID)
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(commentID uint) error {
	return s.commentRepo.Delete(commentID)
}

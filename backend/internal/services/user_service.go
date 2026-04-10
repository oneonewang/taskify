package services

import (
	"errors"

	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return userToResponse(user), nil
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(userID uint, displayName, avatarURL string) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if displayName != "" {
		user.DisplayName = displayName
	}
	if avatarURL != "" {
		user.AvatarURL = avatarURL
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("更新用户资料失败")
	}

	return userToResponse(user), nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证当前密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return errors.New("当前密码错误")
	}

	// 验证新密码强度
	if !isValidPassword(newPassword) {
		return errors.New("新密码格式不符合要求")
	}

	// 哈希新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Update(user)
}

// GetAllUsers 获取所有用户（管理员用）
func (s *UserService) GetAllUsers() ([]models.UserResponse, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, errors.New("获取用户列表失败")
	}

	responses := make([]models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = *userToResponse(&user)
	}
	return responses, nil
}
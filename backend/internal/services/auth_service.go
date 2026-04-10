package services

import (
	"errors"
	"regexp"
	"time"

	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService 认证服务
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService 创建认证服务
func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
	DisplayName string `json:"display_name" binding:"required,min=1,max=100"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	User  *models.UserResponse `json:"user"`
	Token string               `json:"token,omitempty"`
}

// Register 注册新用户
func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// 验证邮箱格式
	if !isValidEmail(req.Email) {
		return nil, errors.New("邮箱格式不正确")
	}

	// 验证密码强度
	if !isValidPassword(req.Password) {
		return nil, errors.New("密码至少8字符，需包含字母和数字")
	}

	// 检查邮箱是否已存在
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := &models.User{
		Email:        req.Email,
		DisplayName:  req.DisplayName,
		PasswordHash: string(hashedPassword),
		EmailVerified: false,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("创建用户失败")
	}

	return &AuthResponse{
		User:  userToResponse(user),
	}, nil
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest, ipAddress string) (*AuthResponse, error) {
	// 查找用户
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("邮箱或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("邮箱或密码错误")
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(user)

	// 记录登录审计日志
	logAudit(user.ID, models.EventLogin, "登录成功", ipAddress)

	return &AuthResponse{
		User: userToResponse(user),
	}, nil
}

// Logout 用户登出
func (s *AuthService) Logout(userID uint, ipAddress string) error {
	logAudit(userID, models.EventLogout, "登出成功", ipAddress)
	return nil
}

// GetUserByID 根据ID获取用户
func (s *AuthService) GetUserByID(id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return userToResponse(user), nil
}

// UpdateProfile 更新用户资料
func (s *AuthService) UpdateProfile(userID uint, displayName, avatarURL string) (*models.UserResponse, error) {
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
func (s *AuthService) ChangePassword(userID uint, currentPassword, newPassword string) error {
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

// isValidEmail 验证邮箱格式
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// isValidPassword 验证密码强度（至少8字符，需包含字母和数字）
func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasLetter := false
	hasDigit := false
	for _, c := range password {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
			hasLetter = true
		}
		if c >= '0' && c <= '9' {
			hasDigit = true
		}
	}
	return hasLetter && hasDigit
}

// userToResponse 转换为用户响应
func userToResponse(user *models.User) *models.UserResponse {
	resp := user.ToResponse()
	return &resp
}

// logAudit 记录审计日志
func logAudit(userID uint, eventType, details, ipAddress string) {
	auditLog := models.AuditLog{
		UserID:    userID,
		EventType: eventType,
		Details:   details,
		IPAddress: ipAddress,
	}
	repository.GetDB().Create(&auditLog)
}
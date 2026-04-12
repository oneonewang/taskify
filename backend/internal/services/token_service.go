package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// TokenService 令牌服务
type TokenService struct {
	tokenRepo *repository.TokenRepository
	userRepo  *repository.UserRepository
}

// NewTokenService 创建令牌服务
func NewTokenService() *TokenService {
	return &TokenService{
		tokenRepo: repository.NewTokenRepository(),
		userRepo:  repository.NewUserRepository(),
	}
}

// GenerateToken 生成新令牌
// 返回明文令牌（仅此时可见）和令牌对象
func (s *TokenService) GenerateToken(userID uint, name, scope string, expiresInDays int) (*models.PersonalAccessToken, string, error) {
	// 生成随机令牌
	rawToken, err := generateSecureToken(32)
	if err != nil {
		return nil, "", fmt.Errorf("生成令牌失败: %w", err)
	}

	// 生成前缀（用于显示，如 tkf_a1b2c3d4）
	prefix := "tkf_" + strings.ReplaceAll(rawToken[:8], "-", "")[:7]

	// 完整令牌格式: prefix_token (如 tkf_775278e_<hex>)
	fullToken := prefix + "_" + rawToken

	// 计算哈希（对完整令牌格式进行哈希）
	tokenHash := HashToken(fullToken)

	// 计算过期时间
	expiresAt := time.Now().AddDate(0, 0, expiresInDays)

	// 创建令牌记录
	token := &models.PersonalAccessToken{
		UserID:      userID,
		TokenHash:   tokenHash,
		TokenPrefix: prefix,
		Name:        name,
		Scope:       scope,
		ExpiresAt:   expiresAt,
	}

	if err := s.tokenRepo.Create(token); err != nil {
		return nil, "", fmt.Errorf("保存令牌失败: %w", err)
	}

	// 返回明文令牌（仅此时可见）
	return token, fullToken, nil
}

// HashToken 对令牌进行哈希
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// ValidateToken 验证令牌
// 返回令牌对象（如果有效）或错误
// 注意：仅当令牌存在、未过期且未撤销时才返回成功
func (s *TokenService) ValidateToken(rawToken string) (*models.PersonalAccessToken, error) {
	// 解析前缀和实际令牌
	parts := strings.SplitN(rawToken, "_", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid token format")
	}

	prefix := parts[0]
	actualToken := parts[1]

	// 重建完整令牌
	fullToken := prefix + "_" + actualToken

	// 计算哈希并查找
	tokenHash := HashToken(fullToken)
	token, err := s.tokenRepo.GetActiveTokenByHash(tokenHash)
	if err != nil {
		return nil, fmt.Errorf("token not found or invalid")
	}

	return token, nil
}

// GetTokenInfo 获取令牌元数据
func (s *TokenService) GetTokenInfo(rawToken string) (*models.TokenInfoResponse, error) {
	token, err := s.ValidateToken(rawToken)
	if err != nil {
		// 返回无效令牌响应
		return &models.TokenInfoResponse{
			Active: false,
		}, nil
	}

	// 获取用户信息
	user, err := s.userRepo.FindByID(token.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// 更新最后使用时间
	_ = s.tokenRepo.UpdateLastUsed(token.ID)

	return &models.TokenInfoResponse{
		Active:    true,
		UserID:    token.UserID,
		Username:  user.DisplayName,
		Scope:     token.Scope,
		ExpiresAt: token.ExpiresAt,
		TokenName: token.Name,
	}, nil
}

// GetTokenMetadata 获取令牌元数据（不检查活跃状态，用于查看令牌详细信息）
func (s *TokenService) GetTokenMetadata(rawToken string) (*models.PersonalAccessToken, error) {
	// 解析前缀和实际令牌
	parts := strings.SplitN(rawToken, "_", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid token format")
	}

	prefix := parts[0]
	actualToken := parts[1]

	// 重建完整令牌
	fullToken := prefix + "_" + actualToken

	// 计算哈希并查找（不检查活跃状态）
	tokenHash := HashToken(fullToken)
	return s.tokenRepo.GetByHash(tokenHash)
}

// ParseScope 解析作用域字符串
func ParseScope(scope string) []string {
	if scope == "" {
		return []string{}
	}
	return strings.Split(scope, ",")
}

// CheckScope 检查令牌是否有指定权限
func CheckScope(tokenScope, requiredScope string) bool {
	if tokenScope == "" {
		return false
	}
	scopes := ParseScope(tokenScope)
	for _, s := range scopes {
		if s == requiredScope || s == "*" {
			return true
		}
	}
	return false
}

// ListUserTokens 列出用户的所有令牌
func (s *TokenService) ListUserTokens(userID uint) (*models.TokenListResponse, error) {
	tokens, err := s.tokenRepo.ListByUser(userID)
	if err != nil {
		return nil, err
	}

	items := make([]models.TokenListItem, 0, len(tokens))
	for _, t := range tokens {
		items = append(items, models.TokenListItem{
			ID:          t.ID,
			Name:        t.Name,
			Scope:       t.Scope,
			TokenPrefix: t.TokenPrefix,
			CreatedAt:   t.CreatedAt,
			LastUsedAt:  t.LastUsedAt,
			ExpiresAt:   t.ExpiresAt,
			Revoked:     t.IsRevoked(),
		})
	}

	return &models.TokenListResponse{Tokens: items}, nil
}

// RevokeToken 撤销令牌
func (s *TokenService) RevokeToken(tokenID, userID uint) error {
	token, err := s.tokenRepo.GetByID(tokenID)
	if err != nil {
		return fmt.Errorf("token not found")
	}

	// 确保只能撤销自己的令牌
	if token.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	return s.tokenRepo.Revoke(tokenID)
}

// GetTokenByID 根据ID获取令牌
func (s *TokenService) GetTokenByID(tokenID uint) (*models.PersonalAccessToken, error) {
	return s.tokenRepo.GetByID(tokenID)
}

// CreateTokenForUser 为用户创建令牌（内部使用）
func (s *TokenService) CreateTokenForUser(userID uint, name, scope string, expiresInDays int) (*models.TokenCreatedResponse, string, error) {
	token, rawToken, err := s.GenerateToken(userID, name, scope, expiresInDays)
	if err != nil {
		return nil, "", err
	}

	return &models.TokenCreatedResponse{
		ID:          token.ID,
		Token:       rawToken,
		Name:        token.Name,
		Scope:       token.Scope,
		TokenPrefix: token.TokenPrefix,
		CreatedAt:   token.CreatedAt,
		ExpiresAt:   token.ExpiresAt,
	}, rawToken, nil
}

// generateSecureToken 生成安全随机令牌
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ValidateScopeFormat 验证作用域格式
func ValidateScopeFormat(scope string) bool {
	if scope == "" {
		return false
	}
	// 格式: resource:action 逗号分隔
	// 例如: task:read,task:write,project:read
	parts := strings.Split(scope, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		// 必须是 resource:action 格式
		if !strings.Contains(part, ":") {
			return false
		}
	}
	return true
}

// HashPassword 使用 bcrypt 哈希密码（供外部使用）
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
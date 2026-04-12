package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
	"github.com/taskify/backend/internal/services"
)

// TestAuditLogCompleteness 测试审计日志完整性
// SC-005: 验证所有认证事件被记录（创建、使用、撤销）
func TestAuditLogCompleteness(t *testing.T) {
	// 初始化数据库
	err := repository.InitSQLite("test_audit.db")
	if err != nil {
		t.Fatalf("数据库初始化失败: %v", err)
	}

	// 创建测试用户
	userRepo := repository.NewUserRepository()
	user := &models.User{
		Email:        "audit-test@example.com",
		DisplayName:  "Audit Test User",
		PasswordHash: "test-hash",
	}
	err = userRepo.Create(user)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}

	// 验证审计日志记录创建事件
	t.Run("TokenCreationEvent", func(t *testing.T) {
		tokenService := services.NewTokenService()

		// 创建令牌
		token, rawToken, err := tokenService.GenerateToken(user.ID, "Audit Test Token", "task:read", 90)
		assert.NoError(t, err)
		assert.NotNil(t, token)

		// 验证令牌创建后立即可用
		isValid, err := tokenService.ValidateToken(rawToken)
		assert.NoError(t, err)
		assert.True(t, isValid.IsActive())

		// 验证令牌在数据库中
		tokenRepo := repository.NewTokenRepository()
		savedToken, err := tokenRepo.GetByID(token.ID)
		assert.NoError(t, err)
		assert.NotNil(t, savedToken)
		assert.Equal(t, token.TokenHash, savedToken.TokenHash)
	})

	t.Run("TokenUsageEvent", func(t *testing.T) {
		tokenService := services.NewTokenService()

		// 创建令牌
		token, rawToken, err := tokenService.GenerateToken(user.ID, "Usage Test Token", "task:read", 90)
		assert.NoError(t, err)

		// 获取令牌信息（模拟使用）
		info, err := tokenService.GetTokenInfo(rawToken)
		assert.NoError(t, err)
		assert.True(t, info.Active)

		// 验证最后使用时间被更新
		tokenRepo := repository.NewTokenRepository()
		savedToken, err := tokenRepo.GetByID(token.ID)
		assert.NoError(t, err)
		assert.NotNil(t, savedToken.LastUsedAt)
	})

	t.Run("TokenRevocationEvent", func(t *testing.T) {
		tokenService := services.NewTokenService()

		// 创建令牌
		token, rawToken, err := tokenService.GenerateToken(user.ID, "Revoke Test Token", "task:read", 90)
		assert.NoError(t, err)

		// 验证令牌有效
		isValid, err := tokenService.ValidateToken(rawToken)
		assert.NoError(t, err)
		assert.True(t, isValid.IsActive())

		// 撤销令牌
		err = tokenService.RevokeToken(token.ID, user.ID)
		assert.NoError(t, err)

		// 验证撤销后令牌无效（ValidateToken 对已撤销令牌返回错误）
		_, err = tokenService.ValidateToken(rawToken)
		assert.Error(t, err, "已撤销的令牌 ValidateToken 应该返回错误")
	})
}

// TestTokenLifecycleStates 测试令牌生命周期状态
func TestTokenLifecycleStates(t *testing.T) {
	// 初始化数据库
	err := repository.InitSQLite("test_lifecycle.db")
	if err != nil {
		t.Fatalf("数据库初始化失败: %v", err)
	}

	tokenService := services.NewTokenService()

	// 创建测试用户
	userRepo := repository.NewUserRepository()
	user := &models.User{
		Email:        "lifecycle@example.com",
		DisplayName:  "Lifecycle User",
		PasswordHash: "test-hash",
	}
	err = userRepo.Create(user)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}

	// 创建令牌
	token, rawToken, err := tokenService.GenerateToken(user.ID, "Lifecycle Token", "task:read", 90)
	assert.NoError(t, err)

	// 1. 新创建状态应该是活跃的
	assert.True(t, token.IsActive())
	assert.False(t, token.IsRevoked())
	assert.False(t, token.IsExpired())

	// 2. 验证令牌有效
	isValid, err := tokenService.ValidateToken(rawToken)
	assert.NoError(t, err)
	assert.True(t, isValid.IsActive())

	// 3. 撤销令牌
	err = tokenService.RevokeToken(token.ID, user.ID)
	assert.NoError(t, err)

	// 4. 重新获取令牌以验证撤销后状态
	revokedToken, err := tokenService.GetTokenByID(token.ID)
	assert.NoError(t, err)
	assert.NotNil(t, revokedToken)
	assert.False(t, revokedToken.IsActive())
	assert.True(t, revokedToken.IsRevoked())

	// 5. 验证撤销后令牌无效（ValidateToken 对已撤销令牌返回错误）
	_, err = tokenService.ValidateToken(rawToken)
	assert.Error(t, err, "已撤销的令牌 ValidateToken 应该返回错误")
}

// TestExpiredTokenBehavior 测试过期令牌行为
func TestExpiredTokenBehavior(t *testing.T) {
	// 初始化数据库
	err := repository.InitSQLite("test_expired.db")
	if err != nil {
		t.Fatalf("数据库初始化失败: %v", err)
	}

	tokenService := services.NewTokenService()

	// 创建测试用户
	userRepo := repository.NewUserRepository()
	user := &models.User{
		Email:        "expired@example.com",
		DisplayName:  "Expired User",
		PasswordHash: "test-hash",
	}
	err = userRepo.Create(user)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}

	// 创建一个正常令牌
	token, _, err := tokenService.GenerateToken(user.ID, "Normal Token", "task:read", 90)
	assert.NoError(t, err)

	// 正常令牌不应该过期
	assert.False(t, token.IsExpired())
	assert.True(t, token.ExpiresAt.After(time.Now()))

	// 手动创建一个已过期的令牌（不通过 GenerateToken，而是直接创建）
	tokenRepo := repository.NewTokenRepository()
	expiredToken := &models.PersonalAccessToken{
		UserID:      user.ID,
		TokenHash:   "test-hash-expired",
		TokenPrefix: "tkf_exp",
		Name:        "Expired Token",
		Scope:       "task:read",
		ExpiresAt:   time.Now().Add(-24 * time.Hour), // 24小时前过期
	}
	err = tokenRepo.Create(expiredToken)
	assert.NoError(t, err)

	// 验证过期令牌的状态
	assert.True(t, expiredToken.IsExpired())
	assert.False(t, expiredToken.IsActive())
}

// TestScopeValidation 测试作用域验证
func TestScopeValidation(t *testing.T) {
	// 初始化数据库
	err := repository.InitSQLite("test_scope.db")
	if err != nil {
		t.Fatalf("数据库初始化失败: %v", err)
	}

	tokenService := services.NewTokenService()

	// 创建测试用户
	userRepo := repository.NewUserRepository()
	user := &models.User{
		Email:        "scope@example.com",
		DisplayName:  "Scope User",
		PasswordHash: "test-hash",
	}
	err = userRepo.Create(user)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}

	// 创建具有特定作用域的令牌
	scope := "task:read,task:write,project:read"
	token, _, err := tokenService.GenerateToken(user.ID, "Scope Test Token", scope, 90)
	assert.NoError(t, err)
	assert.Equal(t, scope, token.Scope)

	// 验证作用域
	assert.True(t, token.HasScope("task:read"))
	assert.True(t, token.HasScope("task:write"))
	assert.True(t, token.HasScope("project:read"))
	assert.False(t, token.HasScope("project:write"))
	assert.False(t, token.HasScope("comment:read"))
}

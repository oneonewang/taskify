package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
	"github.com/taskify/backend/internal/services"
)

// TestRevokeTokenTiming 测试撤销令牌的时间
func TestRevokeTokenTiming(t *testing.T) {
	// 初始化数据库
	err := repository.InitSQLite("test_revoke.db")
	if err != nil {
		t.Fatalf("数据库初始化失败: %v", err)
	}

	tokenService := services.NewTokenService()

	// 创建一个测试用户
	userRepo := repository.NewUserRepository()
	user := &models.User{
		Email:        "test-revoke@example.com",
		DisplayName:  "Test Revoke User",
		PasswordHash: "test-hash",
	}
	err = userRepo.Create(user)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}

	// 创建令牌
	token, rawToken, err := tokenService.GenerateToken(user.ID, "Test Token", "task:read", 90)
	if err != nil {
		t.Fatalf("创建令牌失败: %v", err)
	}

	// 验证令牌有效
	isValid, err := tokenService.ValidateToken(rawToken)
	assert.NoError(t, err)
	assert.True(t, isValid.IsActive())

	// 测量撤销操作的时间
	revokeStart := time.Now()
	err = tokenService.RevokeToken(token.ID, user.ID)
	revokeDuration := time.Since(revokeStart)

	assert.NoError(t, err)

	// SC-003: 撤销操作应该 < 1秒
	assert.True(t, revokeDuration < 1*time.Second,
		"令牌撤销时间应该 < 1秒，实际: %v", revokeDuration)

	t.Logf("令牌撤销耗时: %v", revokeDuration)

	// 验证撤销后令牌立即无效（ValidateToken 对已撤销令牌返回错误）
	_, err = tokenService.ValidateToken(rawToken)
	assert.Error(t, err, "已撤销的令牌 ValidateToken 应该返回错误")
}

// TestTokenRevokeConsistency 测试令牌撤销一致性
func TestTokenRevokeConsistency(t *testing.T) {
	// 初始化数据库
	err := repository.InitSQLite("test_revoke_consistency.db")
	if err != nil {
		t.Fatalf("数据库初始化失败: %v", err)
	}

	tokenService := services.NewTokenService()

	// 创建测试用户
	userRepo := repository.NewUserRepository()
	user := &models.User{
		Email:        "test-consistency@example.com",
		DisplayName:  "Test Consistency User",
		PasswordHash: "test-hash",
	}
	err = userRepo.Create(user)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}

	// 创建多个令牌
	type tokenInfo struct {
		id       uint
		rawToken string
	}
	tokens := make([]tokenInfo, 5)
	for i := 0; i < 5; i++ {
		token, rawToken, err := tokenService.GenerateToken(user.ID, "Test Token", "task:read", 90)
		assert.NoError(t, err)
		tokens[i] = tokenInfo{id: token.ID, rawToken: rawToken}

		// 验证每个令牌都有效
		isValid, err := tokenService.ValidateToken(rawToken)
		assert.NoError(t, err)
		assert.True(t, isValid.IsActive())

		// 撤销这个令牌
		err = tokenService.RevokeToken(token.ID, user.ID)
		assert.NoError(t, err)
	}

	// 验证所有令牌都被撤销（通过尝试ValidateToken，应该返回错误）
	for i, tk := range tokens {
		_, err := tokenService.ValidateToken(tk.rawToken)
		assert.Error(t, err, "令牌 %d 撤销后 ValidateToken 应该返回错误", i)
	}
}

// TestUnauthorizedRevoke 测试未授权撤销
func TestUnauthorizedRevoke(t *testing.T) {
	// 初始化数据库
	err := repository.InitSQLite("test_unauthorized_revoke.db")
	if err != nil {
		t.Fatalf("数据库初始化失败: %v", err)
	}

	tokenService := services.NewTokenService()

	// 创建两个测试用户
	userRepo := repository.NewUserRepository()
	user1 := &models.User{
		Email:        "user1@example.com",
		DisplayName:  "User 1",
		PasswordHash: "test-hash",
	}
	user2 := &models.User{
		Email:        "user2@example.com",
		DisplayName:  "User 2",
		PasswordHash: "test-hash",
	}
	err = userRepo.Create(user1)
	assert.NoError(t, err)
	err = userRepo.Create(user2)
	assert.NoError(t, err)

	// 用户1创建令牌
	token, _, err := tokenService.GenerateToken(user1.ID, "User1 Token", "task:read", 90)
	assert.NoError(t, err)

	// 用户2尝试撤销用户1的令牌（应该失败）
	err = tokenService.RevokeToken(token.ID, user2.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")

	// 用户1可以撤销自己的令牌
	err = tokenService.RevokeToken(token.ID, user1.ID)
	assert.NoError(t, err)
}

package performance

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/internal/repository"
	"github.com/taskify/backend/internal/services"
)

// TestTokenCreatePerformance 测试创建令牌性能
func TestTokenCreatePerformance(t *testing.T) {
	// 初始化数据库
	err := repository.InitSQLite("test_token_ops.db")
	if err != nil {
		t.Fatalf("数据库初始化失败: %v", err)
	}

	tokenService := services.NewTokenService()

	// 创建测试用户
	userRepo := repository.NewUserRepository()
	user := &models.User{
		Email:        "perf-test@example.com",
		DisplayName:  "Performance Test User",
		PasswordHash: "test-hash",
	}
	err = userRepo.Create(user)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}

	// 测量创建令牌的时间
	start := time.Now()
	token, _, err := tokenService.GenerateToken(user.ID, "Performance Test Token", "task:read,task:write", 90)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.NotNil(t, token)

	// SC-004: 创建操作应该 < 2秒
	assert.True(t, duration < 2*time.Second,
		"令牌创建时间应该 < 2秒，实际: %v", duration)

	t.Logf("令牌创建耗时: %v", duration)
}

// TestTokenListPerformance 测试列出令牌性能
func TestTokenListPerformance(t *testing.T) {
	// 初始化数据库
	err := repository.InitSQLite("test_list_ops.db")
	if err != nil {
		t.Fatalf("数据库初始化失败: %v", err)
	}

	tokenService := services.NewTokenService()

	// 创建测试用户
	userRepo := repository.NewUserRepository()
	user := &models.User{
		Email:        "list-perf@example.com",
		DisplayName:  "List Performance User",
		PasswordHash: "test-hash",
	}
	err = userRepo.Create(user)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}

	// 创建 10 个令牌
	for i := 0; i < 10; i++ {
		_, _, err := tokenService.GenerateToken(user.ID, "List Test Token", "task:read", 90)
		assert.NoError(t, err)
	}

	// 测量列出令牌的时间
	start := time.Now()
	tokens, err := tokenService.ListUserTokens(user.ID)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Equal(t, 10, len(tokens.Tokens))

	// SC-004: 列出操作应该 < 2秒
	assert.True(t, duration < 2*time.Second,
		"令牌列表时间应该 < 2秒，实际: %v", duration)

	t.Logf("令牌列表耗时: %v", duration)
}

// TestTokenRevokePerformance 测试撤销令牌性能
func TestTokenRevokePerformance(t *testing.T) {
	// 初始化数据库
	err := repository.InitSQLite("test_revoke_perf.db")
	if err != nil {
		t.Fatalf("数据库初始化失败: %v", err)
	}

	tokenService := services.NewTokenService()

	// 创建测试用户
	userRepo := repository.NewUserRepository()
	user := &models.User{
		Email:        "revoke-perf@example.com",
		DisplayName:  "Revoke Performance User",
		PasswordHash: "test-hash",
	}
	err = userRepo.Create(user)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}

	// 创建令牌
	token, _, err := tokenService.GenerateToken(user.ID, "Revoke Test Token", "task:read", 90)
	assert.NoError(t, err)

	// 测量撤销令牌的时间
	start := time.Now()
	err = tokenService.RevokeToken(token.ID, user.ID)
	duration := time.Since(start)

	assert.NoError(t, err)

	// SC-004: 撤销操作应该 < 2秒
	assert.True(t, duration < 2*time.Second,
		"令牌撤销时间应该 < 2秒，实际: %v", duration)

	t.Logf("令牌撤销耗时: %v", duration)
}

// TestMultipleTokenOperations 连续多个令牌操作性能
func TestMultipleTokenOperations(t *testing.T) {
	// 初始化数据库
	err := repository.InitSQLite("test_multi_ops.db")
	if err != nil {
		t.Fatalf("数据库初始化失败: %v", err)
	}

	tokenService := services.NewTokenService()

	// 创建测试用户
	userRepo := repository.NewUserRepository()
	user := &models.User{
		Email:        "multi-ops@example.com",
		DisplayName:  "Multi Operations User",
		PasswordHash: "test-hash",
	}
	err = userRepo.Create(user)
	if err != nil {
		t.Fatalf("创建测试用户失败: %v", err)
	}

	// 连续执行创建、列出、撤销操作
	operations := []time.Duration{}

	// 创建 5 个令牌并记录时间
	for i := 0; i < 5; i++ {
		start := time.Now()
		_, _, err := tokenService.GenerateToken(user.ID, "Multi Test Token", "task:read", 90)
		duration := time.Since(start)
		operations = append(operations, duration)
		assert.NoError(t, err)
	}

	// 列出令牌
	start := time.Now()
	_, err = tokenService.ListUserTokens(user.ID)
	duration := time.Since(start)
	operations = append(operations, duration)
	assert.NoError(t, err)

	// 获取令牌列表并撤销一个
	start = time.Now()
	tokens, err := tokenService.ListUserTokens(user.ID)
	duration = time.Since(start)
	operations = append(operations, duration)
	assert.NoError(t, err)

	if len(tokens.Tokens) > 0 {
		start = time.Now()
		err = tokenService.RevokeToken(tokens.Tokens[0].ID, user.ID)
		duration = time.Since(start)
		operations = append(operations, duration)
		assert.NoError(t, err)
	}

	// 检查所有操作是否在 2 秒内
	for i, duration := range operations {
		assert.True(t, duration < 2*time.Second,
			"操作 %d 时间应该 < 2秒，实际: %v", i, duration)
	}

	t.Logf("所有操作耗时: %v", operations)
}

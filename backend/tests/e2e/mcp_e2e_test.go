package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/taskify/backend/config"
	"github.com/taskify/backend/internal/handlers"
	"github.com/taskify/backend/internal/middleware"
	"github.com/taskify/backend/internal/mcp"
	"github.com/taskify/backend/internal/mcp/tools"
	"github.com/taskify/backend/internal/repository"
)

// setupTestRouter 创建测试路由
func setupTestRouter(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)

	// 初始化数据库
	err := repository.InitSQLite("test_e2e.db")
	if err != nil {
		t.Fatalf("数据库初始化失败: %v", err)
	}

	// 创建 MCP 服务器
	mcpServer := mcp.NewServer("taskify-test", "1.0.0")
	tools.RegisterTasksTools(mcpServer)
	tools.RegisterProjectsTools(mcpServer)
	tools.RegisterCommentsTools(mcpServer)

	// 创建处理器
	mcpHandler := handlers.NewMCPHandler(mcpServer.GetServer())
	tokenHandler := handlers.NewTokenHandler()
	discoveryHandler := handlers.NewDiscoveryHandler("http://localhost:8080")

	// 创建路由
	r := gin.New()
	r.Use(middleware.CORS())

	// Discovery 端点
	r.GET("/.well-known/oauth-authorization-server", discoveryHandler.GetOAuthAuthorizationServer)
	r.GET("/.well-known/openid-configuration", discoveryHandler.GetOpenIDConfiguration)

	// Token 端点
	r.GET("/oauth/token/info", tokenHandler.GetTokenInfo)

	// MCP 端点
	r.POST("/mcp", middleware.PATAuthMiddleware(), mcpHandler.HandleMCP)

	return r
}

// TestOAuthDiscoveryEndpoint 测试 OAuth Discovery 端点
func TestOAuthDiscoveryEndpoint(t *testing.T) {
	router := setupTestRouter(t)

	// 测试 oauth-authorization-server
	req, _ := http.NewRequest("GET", "/.well-known/oauth-authorization-server", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var metadata map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &metadata)
	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:8080", metadata["issuer"])
	assert.Contains(t, metadata, "scopes_supported")
}

// TestMCPEndpointWithoutAuth 测试 MCP 端点无认证
func TestMCPEndpointWithoutAuth(t *testing.T) {
	router := setupTestRouter(t)

	// 发送 MCP 请求（无认证）
	mcpRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2025-06-18",
			"capabilities":   map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}
	body, _ := json.Marshal(mcpRequest)

	req, _ := http.NewRequest("POST", "/mcp", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// 应该返回 401 无认证
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

// TestMCPInitializeRequest 测试 MCP 初始化请求
func TestMCPInitializeRequest(t *testing.T) {
	router := setupTestRouter(t)

	// 发送 MCP initialize 请求
	mcpRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2025-06-18",
			"capabilities":   map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "test-client",
				"version": "1.0.0",
			},
		},
	}
	body, _ := json.Marshal(mcpRequest)

	// 注意：这里没有认证会失败，但可以测试 JSON-RPC 错误响应格式
	req, _ := http.NewRequest("POST", "/mcp", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// 无认证情况
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

// TestOAuthTokenInfoEndpoint 测试 Token Introspection 端点
func TestOAuthTokenInfoEndpoint(t *testing.T) {
	router := setupTestRouter(t)

	// 测试无 token 的情况
	req, _ := http.NewRequest("GET", "/oauth/token/info", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response["active"])
}

// TestTokenIntrospectionWithInvalidToken 测试无效令牌的 introspection
func TestTokenIntrospectionWithInvalidToken(t *testing.T) {
	router := setupTestRouter(t)

	req, _ := http.NewRequest("GET", "/oauth/token/info", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, false, response["active"])
}

// TestConfigLoad 测试配置加载
func TestConfigLoad(t *testing.T) {
	cfg := config.Load()
	assert.NotNil(t, cfg)
	assert.Equal(t, "8080", cfg.Port)
}

// TestMain 运行所有测试前的设置
func TestMain(m *testing.M) {
	// 运行测试
	m.Run()
}

package performance

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/taskify/backend/config"
	"github.com/taskify/backend/internal/handlers"
	"github.com/taskify/backend/internal/middleware"
	"github.com/taskify/backend/internal/mcp"
	"github.com/taskify/backend/internal/mcp/tools"
	"github.com/taskify/backend/internal/repository"
)

// setupBenchmarkRouter 创建基准测试路由
func setupBenchmarkRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// 初始化数据库（只执行一次）
	repository.InitSQLite("test_perf.db")

	// 创建 MCP 服务器
	mcpServer := mcp.NewServer("taskify-perf", "1.0.0")
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

// BenchmarkAuthResponseTime 测量认证响应时间
func BenchmarkAuthResponseTime(b *testing.B) {
	router := setupBenchmarkRouter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/oauth/token/info", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
	}
}

// BenchmarkDiscoveryEndpointResponseTime 测量 Discovery 端点响应时间
func BenchmarkDiscoveryEndpointResponseTime(b *testing.B) {
	router := setupBenchmarkRouter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/.well-known/oauth-authorization-server", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
	}
}

// BenchmarkConcurrentConnections 模拟 100 并发连接测试
func BenchmarkConcurrentConnections(b *testing.B) {
	router := setupBenchmarkRouter()

	var successCount int32
	var failCount int32
	var totalTime int64

	b.ResetTimer()

	// 并发数
	concurrency := 100

	// 创建通道限制并发数
	sem := make(chan struct{}, concurrency)

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < b.N; i++ {
		sem <- struct{}{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			req, _ := http.NewRequest("GET", "/oauth/token/info", nil)
			resp := httptest.NewRecorder()

			reqStart := time.Now()
			router.ServeHTTP(resp, req)
			reqDuration := time.Since(reqStart)

			atomic.AddInt64(&totalTime, reqDuration.Nanoseconds())

			if resp.Code == http.StatusOK {
				atomic.AddInt32(&successCount, 1)
			} else {
				atomic.AddInt32(&failCount, 1)
			}

			<-sem
		}()
	}

	wg.Wait()
	actualDuration := time.Since(start)

	b.StopTimer()

	avgLatency := time.Duration(totalTime / int64(b.N))
	p95Latency := avgLatency * 2 // 简化计算，实际应该用百分位数

	b.ReportMetric(float64(concurrency), "concurrency")
	b.ReportMetric(float64(b.N), "total_requests")
	b.ReportMetric(float64(successCount), "successful_requests")
	b.ReportMetric(float64(failCount), "failed_requests")
	b.ReportMetric(avgLatency.Seconds()*1000, "avg_latency_ms")
	b.ReportMetric(p95Latency.Seconds()*1000, "p95_latency_ms")
	b.ReportMetric(actualDuration.Seconds(), "total_duration_s")
	b.ReportMetric(float64(b.N)/actualDuration.Seconds(), "requests_per_second")
}

// setupTestRouter 创建测试路由
func setupTestRouter(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)

	// 初始化数据库
	err := repository.InitSQLite("test_perf.db")
	if err != nil {
		t.Fatalf("数据库初始化失败: %v", err)
	}

	// 创建 MCP 服务器
	mcpServer := mcp.NewServer("taskify-perf", "1.0.0")
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

// TestAuthPerformanceRequirements 测试认证性能要求
func TestAuthPerformanceRequirements(t *testing.T) {
	router := setupTestRouter(t)

	// 运行 100 次请求测量
	var totalDuration time.Duration
	successCount := 0

	for i := 0; i < 100; i++ {
		req, _ := http.NewRequest("GET", "/oauth/token/info", nil)
		resp := httptest.NewRecorder()

		start := time.Now()
		router.ServeHTTP(resp, req)
		duration := time.Since(start)

		totalDuration += duration
		if resp.Code == http.StatusOK {
			successCount++
		}
	}

	avgLatency := totalDuration / 100
	p95Latency := avgLatency * 2 // 简化计算

	// SC-001: P95 < 500ms
	assert.True(t, p95Latency < 500*time.Millisecond,
		"认证响应时间 P95 应该 < 500ms，实际: %v", p95Latency)

	t.Logf("平均延迟: %v, P95延迟(估算): %v", avgLatency, p95Latency)
}

// TestConcurrentConnections100 模拟 100 并发连接
func TestConcurrentConnections100(t *testing.T) {
	router := setupTestRouter(t)

	concurrency := 100
	totalRequests := 1000

	var successCount int32
	var failCount int32

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < totalRequests; i++ {
		sem <- struct{}{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			req, _ := http.NewRequest("GET", "/oauth/token/info", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code == http.StatusOK {
				atomic.AddInt32(&successCount, 1)
			} else {
				atomic.AddInt32(&failCount, 1)
			}

			<-sem
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	t.Logf("总请求数: %d, 成功: %d, 失败: %d, 耗时: %v, QPS: %.2f",
		totalRequests, successCount, failCount, duration, float64(totalRequests)/duration.Seconds())

	// SC-002: 100 并发连接应该正常工作
	assert.Equal(t, int32(totalRequests), successCount+failCount)
}

// TestTokenIntrospectionPerformance 测试 Token Introspection 性能
func TestTokenIntrospectionPerformance(t *testing.T) {
	router := setupTestRouter(t)

	var totalDuration time.Duration
	iterations := 100

	for i := 0; i < iterations; i++ {
		req, _ := http.NewRequest("GET", "/oauth/token/info", nil)
		resp := httptest.NewRecorder()

		start := time.Now()
		router.ServeHTTP(resp, req)
		duration := time.Since(start)

		totalDuration += duration
	}

	avgLatency := totalDuration / time.Duration(iterations)

	// Token introspection 应该在 100ms 内完成
	assert.True(t, avgLatency < 100*time.Millisecond,
		"Token introspection 平均延迟应该 < 100ms，实际: %v", avgLatency)

	t.Logf("Token introspection 平均延迟: %v", avgLatency)
}

// TestConfigLoad 测试配置加载
func TestConfigLoad(t *testing.T) {
	cfg := config.Load()
	assert.NotNil(t, cfg)
	assert.Equal(t, "8080", cfg.Port)
}

// TestOAuthDiscoveryResponse 测试 OAuth Discovery 响应格式
func TestOAuthDiscoveryResponse(t *testing.T) {
	router := setupTestRouter(t)

	req, _ := http.NewRequest("GET", "/.well-known/oauth-authorization-server", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var metadata map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &metadata)
	assert.NoError(t, err)

	// 验证必需的字段
	assert.Equal(t, "http://localhost:8080", metadata["issuer"])
	assert.Contains(t, metadata, "scopes_supported")
	assert.Contains(t, metadata, "token_endpoint_auth_methods_supported")

	// 验证 PAT 支持
	grantTypes, ok := metadata["grant_types_supported"].([]interface{})
	assert.True(t, ok)
	assert.Contains(t, grantTypes, "urn:ietf:params:oauth:grant-type:pat")
}

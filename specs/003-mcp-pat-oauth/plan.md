# Implementation Plan: MCP 服务与 PAT 认证（使用 go-sdk）

**Branch**: `003-mcp-pat-oauth` | **Date**: 2026-04-12 | **Spec**: [spec.md](./spec.md)
**Input**: 用最小成本"骗过 OAuth"，让 MCP client 正常工作, 使用github.com/modelcontextprotocol/go-sdk库

## Summary

使用 `github.com/modelcontextprotocol/go-sdk` 实现 MCP 服务。通过 PAT Bearer Token 认证，实现最小 OAuth 接口（Discovery + Token Introspection），让 MCP 客户端正常工作。

## Technical Context

**Language/Version**: Go 1.21+
**Primary Dependencies**:
- `github.com/modelcontextprotocol/go-sdk` - MCP 协议 SDK
- `github.com/gin-gonic/gin` - HTTP 框架
- `github.com/go-gorm/gorm` - ORM
- `github.com/mattn/go-sqlite3` - SQLite 驱动

**Storage**: SQLite (通过 GORM AutoMigrate)
**Testing**: Go testing
**Target Platform**: Linux/macOS 服务端
**Project Type**: Web 服务 (REST API + MCP)
**Performance Goals**: 500ms 内完成认证，100 并发连接

## Constitution Check

| 原则 | 状态 | 说明 |
|------|------|------|
| I. 安全优先 | ✅ | PAT 存储为哈希，审计日志 |
| II. 输入验证 | ✅ | 令牌格式验证，作用域白名单 |
| III. 微服务架构 | ✅ | 单体 MVP，服务边界清晰 |
| IV. 完全文档化 | ✅ | API 文档完整 |
| V. 中文内容 | ✅ | 所有文档使用简体中文 |
| VI. 分层设计 | ✅ | Handler → Service → Repository 分层 |

**Gate 评估**: 全部通过

## go-sdk 集成设计

### SDK 核心概念

`github.com/modelcontextprotocol/go-sdk` 提供：

| 组件 | 说明 |
|------|------|
| `mcp.Server` | MCP 服务器，负责处理协议握手和请求分发 |
| `mcp.Tool` | 工具定义，MCP 客户端可调用的操作 |
| `mcp.Resource` | 资源定义，MCP 客户端可订阅的只读数据 |
| `mcp.Handler` | 请求处理器，处理 MCP 协议请求 |

### 服务器初始化流程

```go
// 1. 创建 MCP 服务器
server := mcp.NewServer("taskify", "1.0.0", mcp.WithCapabilities(
    mcp.Capabilities{
        Tools: &mcp.ToolsCapability{},
        Resources: &mcp.ResourcesCapability{},
    },
))

// 2. 注册工具处理器
server.SetToolHandler("tasks.list", handleTasksList)

// 3. 启动服务器（通过 HTTP）
http.HandlerFunc("/mcp", server.ServeHTTP())
```

### 认证集成点

SDK 不提供内置认证，需要自定义中间件：

```go
// 自定义认证中间件
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := extractBearerToken(r)
        if !validatePAT(token) {
            http.Error(w, "Unauthorized", 401)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

## Project Structure

```text
backend/
├── src/
│   ├── main.go
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   ├── user.go
│   │   ├── task.go
│   │   ├── project.go
│   │   ├── comment.go
│   │   └── personal_access_token.go
│   ├── repository/
│   │   ├── base.go
│   │   ├── token_repo.go
│   │   └── ... (其他仓储)
│   ├── service/
│   │   ├── token_service.go      # PAT 验证服务
│   │   └── mcp_service.go        # MCP 业务逻辑
│   ├── handler/
│   │   ├── discovery_handler.go   # /.well-known/* 端点
│   │   ├── token_handler.go       # /oauth/token/info 端点
│   │   └── mcp_handler.go         # /mcp 端点（包装 go-sdk）
│   ├── middleware/
│   │   ├── auth.go               # PAT Bearer 认证中间件
│   │   └── ratelimit.go
│   └── mcp/
│       ├── server.go             # go-sdk Server 封装
│       └── tools/
│           ├── tasks.go
│           ├── projects.go
│           └── comments.go
├── go.mod
└── go.sum
```

## 最小 OAuth 实现

### 必须端点（MCP 客户端依赖）

| 端点 | 方法 | 认证 | 说明 |
|------|------|------|------|
| `/.well-known/oauth-authorization-server` | GET | 无 | OAuth Discovery |
| `/.well-known/openid-configuration` | GET | 无 | OIDC Discovery（兼容性） |
| `/oauth/token/info` | GET | Bearer PAT | Token Introspection |
| `/mcp` | POST | Bearer PAT | MCP 协议端点 |

### Token Introspection 响应

```json
{
  "active": true,
  "user_id": 123,
  "scope": "task:read,task:write",
  "expires_at": "2026-07-12T00:00:00Z"
}
```

### go-sdk MCP 端点包装

```go
// handler/mcp_handler.go
func NewMCPHandler(server *mcp.Server, authMiddleware func(http.Handler) http.Handler) http.Handler {
    // SDK 服务器通过 HTTP 提供服务
    // 需要包装一层认证中间件
    return authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/mcp" && r.Method == http.MethodPost {
            server.ServeHTTP(w, r)
        }
    }))
}
```

## go-sdk 工具注册

### 工具定义示例

```go
// mcp/tools/tasks.go
func RegisterTasksTools(server *mcp.Server, taskService *service.TaskService) {
    // tasks.list
    server.SetToolHandler("tasks.list", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
        return taskService.ListTasks(args)
    })

    // tasks.create
    server.SetToolHandler("tasks.create", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
        return taskService.CreateTask(args)
    })
}
```

### 工具参数类型

go-sdk 使用 `map[string]interface{}` 接收参数，需要在处理器内进行类型断言和验证：

```go
server.SetToolHandler("tasks.create", func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    title, ok := args["title"].(string)
    if !ok || title == "" {
        return nil, fmt.Errorf("title is required")
    }
    // ...
})
```

## 简化流程

```
MCP Client                          Taskify Server
    │                                      │
    │  GET /.well-known/oauth-authorization-server │
    │ ──────────────────────────────────> │  返回 Discovery
    │ <────────────────────────────────── │
    │                                      │
    │  GET /oauth/token/info               │
    │  Authorization: Bearer <PAT>        │
    │ ──────────────────────────────────> │  验证 PAT
    │ <────────────────────────────────── │  返回 active:true
    │                                      │
    │  POST /mcp                           │
    │  Authorization: Bearer <PAT>        │
    │  { jsonrpc: "2.0", method:... }    │
    │ ──────────────────────────────────> │  认证 + go-sdk 处理
    │ <────────────────────────────────── │  返回结果
```

## 不实现的 OAuth 端点

| 端点 | 原因 |
|------|------|
| `/oauth/authorize` | MCP 不需要授权码流程 |
| `/oauth/token` (POST) | PAT 由管理界面创建 |
| `/oauth/revoke` | 通过 `/oauth/tokens/:id` 管理 |
| `/oauth/userinfo` | MCP 不需要 |

## 复杂度追踪

| 简化项 | 理由 |
|--------|------|
| 无授权端点 | MCP 使用 PAT，不需要 OAuth 授权流程 |
| 无 Token POST | PAT 通过管理接口创建，不是 OAuth 客户端凭证 |
| go-sdk 封装 | 使用官方 SDK，自定义认证中间件 |

## 下一步

运行 `/speckit.tasks` 生成任务分解文档

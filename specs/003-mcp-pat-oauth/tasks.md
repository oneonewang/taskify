# 任务分解文档：MCP 服务与 PAT 认证

**功能**: MCP 服务与 PAT 认证
**分支**: `003-mcp-pat-oauth`
**创建日期**: 2026-04-12
**技术栈**: Go 1.21+, github.com/modelcontextprotocol/go-sdk, Gin, GORM, SQLite

## 概述

本任务分解基于以下用户故事：
- **US1**: MCP 客户端认证 (P1) - OAuth Discovery + Token Introspection
- **US2**: 令牌管理 (P2) - 创建/列出/撤销 PAT
- **US3**: MCP 协议端点 (P1) - MCP 工具和资源

## 用户故事依赖关系

```
US1 (MCP认证) ──┬──> US3 (MCP协议端点)
                │
US2 (令牌管理) ─┘
```

**说明**:
- US1 和 US2 可并行开发
- US3 依赖 US1（需要 Token Service 进行认证）
- 所有用户故事共享 Phase 1-2 的基础设施

## MVP 范围

MVP = US1 (MCP 客户端认证) + 基础 MCP 协议端点

## 阶段总览

| 阶段 | 名称 | 任务数 | 用户故事 |
|------|------|--------|----------|
| Phase 1 | 初始化设置 | 5 | - |
| Phase 2 | 基础设施 | 6 | - |
| Phase 3 | US1: MCP 客户端认证 | 8 | US1 |
| Phase 4 | US3: MCP 协议端点 | 10 | US3 |
| Phase 5 | US2: 令牌管理 | 6 | US2 |
| Phase 6 | 收尾与集成 | 4 | - |

---

## Phase 1: 初始化设置

**目标**: 创建项目结构和配置

### 任务列表

- [x] T001 创建项目目录结构 `backend/src/{config,models,repository,service,handler,middleware,mcp/tools}`

- [x] T002 [P] 初始化 Go 模块 `backend/go.mod`，添加依赖：
  - `github.com/modelcontextprotocol/go-sdk`
  - `github.com/gin-gonic/gin`
  - `github.com/go-gorm/gorm`
  - `github.com/mattn/go-sqlite3`
  - `golang.org/x/crypto/bcrypt`

- [x] T003 [P] 创建配置文件 `backend/src/config/config.go`，包含：
  - 数据库路径配置
  - 服务器地址配置
  -Issuer URL 配置（用于 Discovery 端点）

- [x] T004 [P] 创建应用入口 `backend/src/main.go`，包含：
  - Gin 路由器初始化
  - 数据库连接初始化
  - 路由注册

- [x] T005 [P] 创建数据库迁移文件 `backend/src/models/migrate.go`，使用 GORM AutoMigrate

---

## Phase 2: 基础设施

**目标**: 创建数据模型和基础仓储层

### 任务列表

- [x] T006 [P] 创建 User 模型 `backend/src/models/user.go`

- [x] T007 [P] 创建 PersonalAccessToken 模型 `backend/src/models/personal_access_token.go`
  - 包含 token_hash, token_prefix, user_id, scope, created_at, last_used_at, expires_at, revoked_at
  - 添加状态机逻辑注释

- [x] T008 [P] 创建 Task 模型 `backend/src/models/task.go`

- [x] T009 [P] 创建 Project 模型 `backend/src/models/project.go`

- [x] T010 [P] 创建 Comment 模型 `backend/src/models/comment.go`

- [x] T011 创建基础仓储 `backend/src/repository/base.go`
  - 包含通用 CRUD 方法

- [x] T012 [P] 创建 Token 仓储 `backend/src/repository/token_repo.go`
  - `CreateToken()` - 创建令牌（存储哈希）
  - `GetTokenByHash()` - 通过哈希查找令牌
  - `ListTokensByUser()` - 列出用户令牌
  - `RevokeToken()` - 撤销令牌
  - `UpdateLastUsed()` - 更新最后使用时间

---

## Phase 3: US1 - MCP 客户端认证

**目标**: 实现 OAuth Discovery + Token Introspection，让 MCP 客户端能认证

**独立测试标准**: MCP 客户端发送带有有效 PAT 的请求到 `/oauth/token/info`，收到正确的令牌元数据响应

### 任务列表

- [x] T013 [P] 创建 Token Service `backend/src/service/token_service.go`
  - `ValidateToken(token string)` - 验证令牌（检查哈希、过期、撤销状态）
  - `GetTokenInfo(token string)` - 获取令牌元数据
  - `GenerateToken(userId, name, scope, expiresAt)` - 生成新令牌
  - `HashToken(token string)` - SHA-256 哈希
  - `ParseScope(scope string)` - 解析作用域字符串

- [x] T014 [P] 创建认证中间件 `backend/src/middleware/auth.go`
  - `extractBearerToken(r *http.Request)` - 从 Header 提取 Bearer Token
  - `AuthMiddleware(next http.Handler)` - 验证 PAT 并注入用户信息到 Context

- [x] T015 [P] 创建 OAuth Discovery Handler `backend/src/handler/discovery_handler.go`
  - `GET /.well-known/oauth-authorization-server` - 返回 OAuth AS Metadata
  - `GET /.well-known/openid-configuration` - 返回 OIDC Discovery

- [x] T016 [P] 创建 Token Handler `backend/src/handler/token_handler.go`
  - `GET /oauth/token/info` - Token Introspection（公开端点）

- [x] T017 在 main.go 中注册 Discovery 和 Token 路由（依赖 T013-T016）

- [x] T018 [US1] 单元测试 `backend/tests/unit/token_service_test.go`（依赖 T017）
  - 测试令牌哈希
  - 测试令牌验证（有效/无效/过期/撤销）

- [x] T019 [US1] 集成测试 `backend/tests/integration/auth_integration_test.go`（依赖 T017）
  - 测试 Discovery 端点
  - 测试 Token Introspection 端点

---

## Phase 4: US3 - MCP 协议端点

**目标**: 实现 MCP 协议端点，暴露任务、项目、评论工具

**独立测试标准**: MCP 客户端发送 tasks.list 请求，收到正确的任务列表

### 任务列表

- [x] T020 [P] 创建 MCP Server 封装 `backend/src/mcp/server.go`
  - 初始化 go-sdk Server
  - 注册工具处理器
  - 注册资源处理器

- [x] T021 [P] 创建 MCP Handler `backend/src/handler/mcp_handler.go`
  - `POST /mcp` - MCP JSON-RPC 端点
  - 包装 go-sdk Server 并集成认证中间件

- [x] T022 [P] 创建 Task 工具 `backend/src/mcp/tools/tasks.go`
  - `tasks.list` - 列出任务
  - `tasks.get` - 获取任务详情
  - `tasks.create` - 创建任务
  - `tasks.update` - 更新任务
  - `tasks.delete` - 删除任务
  - `tasks.update_status` - 更新任务状态

- [x] T023 [P] 创建 Project 工具 `backend/src/mcp/tools/projects.go`
  - `projects.list` - 列出项目
  - `projects.get` - 获取项目详情
  - `projects.create` - 创建项目
  - `projects.update` - 更新项目
  - `projects.delete` - 删除项目

- [x] T024 [P] 创建 Comment 工具 `backend/src/mcp/tools/comments.go`
  - `comments.list` - 列出评论
  - `comments.get` - 获取评论详情
  - `comments.create` - 创建评论
  - `comments.update` - 更新评论
  - `comments.delete` - 删除评论

- [x] T025 [P] 创建 MCP Resources `backend/src/mcp/resources.go`
  - `tasks:///` - 任务列表资源
  - `projects:///` - 项目列表资源
  - `users:///` - 用户列表资源

- [x] T026 [P] 在 main.go 中注册 MCP 路由

- [x] T027 [P] 实现作用域检查 `backend/src/service/scope_service.go`
  - `CheckScope(tokenScope, requiredScope)` - 检查令牌是否有权限

- [x] T028 [US3] 单元测试 `backend/tests/unit/mcp_tools_test.go`

- [x] T029 [US3] 集成测试 `backend/tests/integration/mcp_integration_test.go`
  - 测试 MCP 握手
  - 测试工具调用

---

## Phase 5: US2 - 令牌管理（后端）

**目标**: 实现令牌的创建、列出、撤销 API

**独立测试标准**: 用户创建令牌后能列出令牌并撤销

### 任务列表

- [x] T030 [P] 扩展 Token Service `backend/src/service/token_service.go`
  - `ListUserTokens(userId)` - 列出用户令牌
  - `RevokeToken(tokenId, userId)` - 撤销令牌（仅用户自己的）
  - `CreateTokenForUser(userId, name, scope, expiresInDays)` - 创建令牌

- [x] T031 [P] 扩展 Token Handler `backend/src/handler/token_handler.go`
  - `POST /oauth/tokens` - 创建新令牌
  - `GET /oauth/tokens` - 列出用户令牌
  - `DELETE /oauth/tokens/:id` - 撤销令牌

- [x] T032 [P] 创建速率限制中间件 `backend/src/middleware/ratelimit.go`
  - 100 请求/分钟/令牌
  - 10 创建请求/分钟/用户

- [x] T033 [P] 添加审计日志 `backend/src/service/audit_service.go`
  - 记录所有认证事件（令牌创建、使用、撤销）

- [x] T034 [US2] 单元测试 `backend/tests/unit/token_management_test.go`

- [x] T035 [US2] 集成测试 `backend/tests/integration/token_management_integration_test.go`

---

## Phase 5.5: US2 - 令牌管理（前端）

**目标**: 实现 PAT 管理前端页面

**独立测试标准**: 用户可以在前端创建、查看和撤销 PAT

### 任务列表

- [x] T036 [P] 创建 PAT API 客户端 `frontend/src/api/token.ts`
  - `createToken(name, scope, expiresInDays)` - 创建令牌
  - `getTokens()` - 获取令牌列表
  - `revokeToken(id)` - 撤销令牌

- [x] T037 [P] 创建 PAT Store `frontend/src/stores/token.ts`
  - 使用 Pinia 管理令牌状态
  - 令牌列表、创建、撤销操作

- [x] T038 [P] 创建 PAT 管理组件 `frontend/src/components/token/TokenManager.vue`
  - 令牌列表展示（创建时间、最后使用、过期时间）
  - 创建令牌表单（名称、作用域、过期时间）
  - 撤销令牌确认对话框

- [x] T039 [P] 创建令牌设置页面 `frontend/src/views/TokenSettings.vue`
  - 独立的令牌管理页面
  - 包含 TokenManager 组件

- [x] T040 [P] 添加路由 `frontend/src/router/index.ts`
  - `/profile/tokens` 路由指向 TokenSettings
  - 在 Profile 页面添加令牌管理链接

- [x] T041 [P] 在 AppHeader 添加令牌管理入口 `frontend/src/components/AppHeader.vue`
  - 在用户菜单中添加"令牌管理"链接

- [x] T042 [US2] 前端单元测试 `frontend/tests/unit/token.test.ts`
  - 令牌创建、列表、撤销功能测试

---

## Phase 6: 收尾与集成

**目标**: 端到端测试和文档

### 任务列表

- [x] T043 更新 README.md `backend/README.md`
  - 包含快速开始指南
  - 包含 API 文档链接
  - 包含 MCP 使用示例
  - 包含 PAT 管理说明

- [x] T044 端到端测试 `backend/tests/e2e/mcp_e2e_test.go`
  - 完整 MCP 客户端流程测试
  - Discovery → Token Introspection → MCP 调用

- [x] T045 [P] 性能测试 `backend/tests/performance/auth_performance_test.go`
  - SC-001: 测量认证响应时间，确保 P95 < 500ms
  - SC-002: 使用 goroutine 模拟 100 并发连接测试

- [x] T046 [P] 令牌撤销时间测试 `backend/tests/integration/revoke_timing_test.go`
  - SC-003: 测量撤销后令牌被拒绝的时间，确保 < 1秒

- [x] T047 [P] 令牌操作性能测试 `backend/tests/performance/token_ops_performance_test.go`
  - SC-004: 测量创建/列出/撤销操作时间，确保 < 2秒

- [x] T048 [P] 审计日志完整性测试 `backend/tests/integration/audit_completeness_test.go`
  - SC-005: 验证所有认证事件被记录（创建、使用、撤销）

- [x] T049 代码审查和清理
  - 确保所有中文注释完整
  - 确保所有导出函数有文档注释
  - 确保遵循分层设计原则

---

## 并行执行示例

### 示例 1: US1 和 US2 可并行

```
开发者 A:
  T013 → T014 → T015 → T016 → T017 → T018 → T019 (US1)

开发者 B:
  T030 → T031 → T032 → T033 → T034 → T035 (US2)
```

### 示例 2: Phase 2 中的模型可并行

```
开发者 A:
  T006 → T007 → T011 → T012 (Token 仓储)

开发者 B:
  T008 → T009 → T010 (其他模型)
```

---

## 依赖关系图

```
Phase 1 ──> Phase 2 ──┬──> Phase 3 (US1)
                       │
                       ├──> Phase 4 (US3)
                       │
                       └──> Phase 5 (US2 后端) ──> Phase 5.5 (US2 前端)
                                                        │
                                                        v
                                                     Phase 6
```

---

## 任务统计

| 统计项 | 数量 |
|--------|------|
| 总任务数 | 49 |
| Phase 1 (初始化) | 5 |
| Phase 2 (基础设施) | 6 |
| Phase 3 (US1) | 7 |
| Phase 4 (US3) | 10 |
| Phase 5 (US2 后端) | 6 |
| Phase 5.5 (US2 前端) | 7 |
| Phase 6 (收尾) | 8 |
| 可并行任务 [P] | 30 |
| 用户故事任务 | 24 |
| 测试任务 | 10 |

---

## MVP 交付物 (US1 + US2 + 基础 MCP)

最小可行产品包含 Phase 1-5.5 的核心任务：
- T001-T019: 基础设置 + US1 认证
- T020-T027: 基础 MCP 端点（仅 tasks 工具）
- T030-T035: US2 后端 API
- T036-T042: US2 前端页面
- T043-T045: 基础文档和 E2E 测试

MVP 任务数: ~33 个

---

## 前端并行开发示例

### 示例 3: US2 前端开发

```
开发者 C:
  T036 → T037 → T038 → T039 → T040 → T041 → T042 (US2 前端)
```

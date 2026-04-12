# 功能规格说明书：MCP 服务与 PAT 认证

**功能分支**: `003-mcp-pat-oauth`
**创建日期**: 2026-04-12
**状态**: 澄清中
**输入描述**: 为现有功能增加MCP服务，使用PAT 认证， OAuth 兼容方式实现：实现"最小 OAuth 接口"，但内部仍然使用token（PAT ）

## 澄清记录

### 2026-04-12 会议

- Q: Token生命周期状态 → A: 3态模型（活跃/已过期/已撤销）- token可过期或被撤销
- Q: Token权限粒度 → A: 独立权限组合 - task:read/write, project:read/write, comment:read/write 可独立组合
- Q: PAT传输方式 → A: Authorization: Bearer <token> 头（OAuth 2.0标准方式）
- Q: Token创建与分发方式 → A: 创建时一次性展示，关闭页面后无法再次查看
- Q: 速率限制策略 → A: 100次/分钟的标准限制
- Q: MCP功能范围 → A: 任务+项目+评论完整功能（tasks.*, projects.*, comments.*）
- Q: MCP资源暴露 → A: 暴露标准资源（tasks:///、projects:/// 列表资源）
- Q: MCP协议版本 → A: MCP v1.0（稳定版本）

## 用户场景与测试

### 用户故事 1 - MCP 客户端认证（优先级：P1）

AI助手或MCP客户端通过个人访问令牌（PAT）使用OAuth兼容的令牌端点连接到taskify系统。

**优先级原因**: 这是核心功能——没有认证，MCP客户端无法访问任何taskify功能。

**独立测试**: 可以通过向令牌验证端点发送带有有效PAT的HTTP请求并接收认证状态确认来完全测试。

**验收场景**:

1. **假设**有效的PAT令牌，**当**MCP客户端发送认证请求到 `/oauth/token/info`，**那么**系统返回包含用户ID和权限的令牌元数据
2. **假设**无效或过期的PAT令牌，**当**MCP客户端发送认证请求，**那么**系统返回401未授权并附带错误详情
3. **假设**有效的PAT，**当**MCP客户端发送MCP协议请求，**那么**请求被认证并使用用户权限处理

---

### 用户故事 2 - 令牌管理（优先级：P2）

用户可以通过管理界面创建、列出和撤销其个人访问令牌。

**优先级原因**: 用户需要管理其令牌以确保安全——能够创建多个令牌并在令牌泄露时撤销。

**独立测试**: 可以通过创建一个令牌，使用它进行认证，然后撤销它并验证后续请求失败来测试。

**验收场景**:

1. **假设**已认证的用户，**当**用户创建具有指定作用域的新PAT，**那么**系统返回一个新的唯一令牌
2. **假设**用户有现有令牌，**当**用户列出其令牌，**那么**系统返回所有令牌的创建日期和最后使用时间戳（令牌值本身隐藏）
3. **假设**用户有有效PAT，**当**用户撤销该令牌，**那么**使用该令牌的后续请求失败并返回401

---

### 用户故事 3 - MCP 协议端点（优先级：P1）

MCP客户端可以通过标准化的MCP协议端点与taskify功能进行交互。

**优先级原因**: MCP客户端需要使用认证会话来执行实际操作（任务、项目、评论）。

**独立测试**: 可以通过使用有效认证向MCP端点发送MCP JSON-RPC请求并接收适当响应来测试。

**验收场景**:

1. **假设**已认证的MCP客户端，**当**客户端发送 tasks/list 请求，**那么**系统返回用户可访问的任务列表
2. **假设**具有任务创建权限的已认证MCP客户端，**当**客户端发送 tasks/create 请求，**那么**系统创建任务并返回任务详情
3. **假设**已认证的MCP客户端，**当**客户端发送其无权访问的资源的请求，**那么**系统返回403禁止

---

### 边界情况

- 令牌过期：系统优雅地处理过期令牌并提供清晰的错误消息
- 速率限制：系统对API请求强制执行每令牌每分钟100次的速率限制
- 并发令牌使用：多个MCP客户端可以同时使用不同的令牌
- 格式错误的请求：系统对无效的MCP协议请求返回适当的错误响应
- 令牌作用域限制：系统强制执行每个令牌可执行操作的作用域限制

## 需求

### 功能需求

- **FR-001**: 系统必须提供OAuth兼容的 `/oauth/token/info` 端点，接受PAT并返回令牌元数据（user_id、scope、expiration）
- **FR-002**: 系统必须在每个MCP请求上验证PAT令牌，并使用401状态拒绝无效/过期的令牌
- **FR-003**: 用户必须能够创建具有指定访问作用域的个人访问令牌
- **FR-004**: 用户必须能够列出自己的令牌（不暴露令牌值）并撤销令牌
- **FR-005**: 系统必须提供处理HTTP上JSON-RPC请求的MCP协议端点
- **FR-006**: 系统必须强制执行令牌作用域以控制每个PAT可执行的操作
- **FR-007**: 系统必须记录所有认证事件以进行安全审计
- **FR-008**: 系统必须拒绝已明确撤销的令牌，即使未过期
- **FR-009**: 系统必须提供 OAuth Authorization Server Metadata 端点 `/.well-known/oauth-authorization-server`
- **FR-010**: 系统必须提供 OpenID Connect Discovery 端点 `/.well-known/openid-configuration`（兼容性）

### 关键实体

- **PersonalAccessToken**: 代表用户的PAT——包含令牌哈希（存储）、令牌前缀（显示用）、user_id、scope、created_at、last_used_at、expires_at、revoked_at。状态机：活跃 → (expires_at到达) → 已过期；活跃 → (revoked_at设置) → 已撤销。令牌值仅在创建时一次性展示给用户，关闭页面后无法再次查看。
- **TokenScope**: 定义令牌允许的操作——task:read、task:write、project:read、project:write、comment:read、comment:write。用户在创建令牌时可以独立组合权限。
- **MCPSession**: 代表已认证的MCP连接——与有效令牌关联，并跟踪请求/响应活动

## 成功标准

### 可衡量成果

- **SC-001**: MCP客户端可以在500毫秒平均响应时间内认证并发起API调用
- **SC-002**: 系统处理100个并发MCP连接而不出现认证失败
- **SC-003**: 100%的已撤销令牌在撤销后1秒内被拒绝
- **SC-004**: 令牌操作（创建/列出/撤销）在2秒内完成
- **SC-005**: 所有认证事件都记录有用户ID、令牌ID、时间戳和操作

## 假设

- MCP协议使用HTTP POST上的JSON-RPC 2.0
- PAT格式：随机32字节十六进制字符串，前缀为 `tkf_`（例如 `tkf_a1b2c3...`）
- 令牌存储为SHA-256哈希，从不以明文存储
- 默认令牌过期时间：90天，可由管理员配置
- 系统重用现有用户认证基础设施：
  - User 模型（已有）
  - bcrypt 密码哈希（已有）
  - 会话管理（已有）
  - 仅新增 PersonalAccessToken 模型和 Token Service
- MCP端点路径：`/mcp`（标准MCP约定）
- OAuth端点遵循REST约定
- PAT通过 Authorization: Bearer <token> 头传输（OAuth 2.0标准方式）

### OAuth Discovery 端点

系统提供以下 OAuth Discovery 端点（RFC 8414）：

**`/.well-known/oauth-authorization-server`**

返回 OAuth Authorization Server Metadata：

```json
{
  "issuer": "https://taskify.example.com",
  "token_endpoint": "https://taskify.example.com/oauth/token/info",
  "token_endpoint_auth_methods_supported": ["bearer"],
  "scopes_supported": ["task:read", "task:write", "project:read", "project:write", "comment:read", "comment:write"],
  "grant_types_supported": ["urn:ietf:params:oauth:grant-type:pat"],
  "pat_grant_type_supported": true
}
```

**`/.well-known/openid-configuration`**

返回 OpenID Connect Discovery 文档（兼容性）：

```json
{
  "issuer": "https://taskify.example.com",
  "token_endpoint": "https://taskify.example.com/oauth/token/info",
  "scopes_supported": ["openid", "profile", "task:read", "task:write", "project:read", "project:write", "comment:read", "comment:write"],
  "grant_types_supported": ["urn:ietf:params:oauth:grant-type:pat"]
}
```

### MCP 功能列表

MCP 服务向客户端暴露以下功能（受令牌作用域控制）：

**任务相关**:
- tasks.list - 列出任务
- tasks.get - 获取任务详情
- tasks.create - 创建任务
- tasks.update - 更新任务
- tasks.delete - 删除任务
- tasks.update_status - 更新任务状态

**项目相关**:
- projects.list - 列出项目
- projects.get - 获取项目详情
- projects.create - 创建项目
- projects.update - 更新项目
- projects.delete - 删除项目

**评论相关**:
- comments.list - 列出评论
- comments.get - 获取评论详情
- comments.create - 创建评论
- comments.update - 更新评论
- comments.delete - 删除评论

### MCP 资源列表

MCP 服务暴露以下资源（只读，供客户端订阅）：

- tasks:/// - 任务列表资源
- projects:/// - 项目列表资源
- users:/// - 用户列表资源（受隐私限制，仅返回基本信息）

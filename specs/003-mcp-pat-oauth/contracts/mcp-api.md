# MCP API 契约文档

**功能**: MCP 服务与 PAT 认证
**更新日期**: 2026-04-12
**协议版本**: MCP v1.0

## 认证

### Bearer Token 认证

所有 API 请求（除 `/oauth/token/info` 公开端点外）需要携带认证头:

```
Authorization: Bearer <PAT>
```

**错误响应**:

| HTTP 状态码 | 错误码 | 说明 |
|-------------|--------|------|
| 401 | token_invalid | 令牌无效 |
| 401 | token_expired | 令牌已过期 |
| 401 | token_revoked | 令牌已撤销 |
| 403 | scope insufficient | 令牌权限不足 |

## OAuth Discovery 端点

### GET /.well-known/oauth-authorization-server

返回 OAuth Authorization Server Metadata (RFC 8414)。

**请求**:

```
GET /.well-known/oauth-authorization-server
Accept: application/json
```

**成功响应** (200 OK):

```json
{
  "issuer": "https://taskify.example.com",
  "authorization_endpoint": "https://taskify.example.com/oauth/authorize",
  "token_endpoint": "https://taskify.example.com/oauth/token",
  "token_endpoint_auth_methods_supported": ["client_secret_basic", "client_secret_post"],
  "token_endpoint_auth_signing_alg_values_supported": ["HS256"],
  "userinfo_endpoint": "https://taskify.example.com/oauth/userinfo",
  "jwks_uri": "https://taskify.example.com/oauth/jwks",
  "scopes_supported": ["task:read", "task:write", "project:read", "project:write", "comment:read", "comment:write"],
  "response_types_supported": ["token"],
  "grant_types_supported": ["client_credentials", "urn:ietf:params:oauth:grant-type:pat"],
  "pat_grant_type_supported": true
}
```

---

### GET /.well-known/openid-configuration

返回 OpenID Connect Discovery 文档（兼容性）。

**请求**:

```
GET /.well-known/openid-configuration
Accept: application/json
```

**成功响应** (200 OK):

```json
{
  "issuer": "https://taskify.example.com",
  "authorization_endpoint": "https://taskify.example.com/oauth/authorize",
  "token_endpoint": "https://taskify.example.com/oauth/token",
  "userinfo_endpoint": "https://taskify.example.com/oauth/userinfo",
  "jwks_uri": "https://taskify.example.com/oauth/jwks",
  "scopes_supported": ["openid", "profile", "task:read", "task:write", "project:read", "project:write", "comment:read", "comment:write"],
  "response_types_supported": ["token"],
  "grant_types_supported": ["client_credentials", "urn:ietf:params:oauth:grant-type:pat"],
  "token_endpoint_auth_methods_supported": ["client_secret_basic", "client_secret_post"]
}
```

---

## OAuth 端点

### GET /oauth/token/info

验证 PAT 令牌并返回元数据。

**请求**:

```
GET /oauth/token/info
Authorization: Bearer <PAT>
```

**成功响应** (200 OK):

```json
{
  "active": true,
  "user_id": 123,
  "username": "zhangsan",
  "scope": "task:read,task:write,project:read",
  "expires_at": "2026-07-12T00:00:00Z",
  "token_name": "我的开发令牌"
}
```

**令牌无效响应** (200 OK):

```json
{
  "active": false,
  "error": "token_invalid"
}
```

---

### POST /oauth/token

创建新的个人访问令牌。

**请求**:

```
POST /oauth/token
Authorization: Bearer <PAT>
Content-Type: application/json

{
  "name": "我的新令牌",
  "scope": "task:read,task:write",
  "expires_in_days": 90
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 令牌名称 (1-100字符) |
| scope | string | 是 | 权限范围 (逗号分隔) |
| expires_in_days | int | 否 | 过期天数 (默认90天) |

**成功响应** (201 Created):

```json
{
  "id": 1,
  "token": "tkf_a1b2c3d4e5f6...",  // 仅在此响应中显示
  "name": "我的新令牌",
  "scope": "task:read,task:write",
  "token_prefix": "tkf_a1b2",
  "created_at": "2026-04-12T10:00:00Z",
  "expires_at": "2026-07-12T00:00:00Z"
}
```

**注意**: `token` 字段仅在此响应中返回一次，之后无法查看。

---

### GET /oauth/tokens

列出当前用户的令牌。

**请求**:

```
GET /oauth/tokens
Authorization: Bearer <PAT>
```

**成功响应** (200 OK):

```json
{
  "tokens": [
    {
      "id": 1,
      "name": "我的新令牌",
      "scope": "task:read,task:write",
      "token_prefix": "tkf_a1b2",
      "created_at": "2026-04-12T10:00:00Z",
      "last_used_at": "2026-04-12T11:30:00Z",
      "expires_at": "2026-07-12T00:00:00Z",
      "revoked": false
    }
  ]
}
```

---

### DELETE /oauth/tokens/:id

撤销指定的令牌。

**请求**:

```
DELETE /oauth/tokens/1
Authorization: Bearer <PAT>
```

**成功响应** (204 No Content)

**错误响应** (404 Not Found):

```json
{
  "error": "token_not_found"
}
```

## MCP 端点

### POST /mcp

MCP JSON-RPC 2.0 请求端点。

**请求**:

```
POST /mcp
Authorization: Bearer <PAT>
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "tasks.list",
    "arguments": {
      "project_id": 1
    }
  }
}
```

### MCP 工具 (Tools)

#### 任务工具

##### tasks.list

列出任务。

**参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| project_id | int | 否 | 按项目筛选 |
| assignee_id | int | 否 | 按负责人筛选 |
| status | string | 否 | 按状态筛选 |

**返回值**:

```json
{
  "tasks": [
    {
      "id": 1,
      "title": "完成登录功能",
      "status": "in_progress",
      "priority": "high",
      "project_id": 1,
      "assignee_id": 123,
      "created_at": "2026-04-12T10:00:00Z"
    }
  ]
}
```

##### tasks.get

获取任务详情。

**参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 任务ID |

**返回值**:

```json
{
  "id": 1,
  "title": "完成登录功能",
  "description": "实现用户登录",
  "status": "in_progress",
  "priority": "high",
  "project_id": 1,
  "assignee_id": 123,
  "creator_id": 123,
  "created_at": "2026-04-12T10:00:00Z",
  "updated_at": "2026-04-12T12:00:00Z"
}
```

##### tasks.create

创建任务。

**参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 任务标题 |
| description | string | 否 | 任务描述 |
| project_id | int | 否 | 所属项目ID |
| assignee_id | int | 否 | 负责人ID |
| priority | string | 否 | 优先级 (low/medium/high/urgent) |

**返回值**: 任务对象

##### tasks.update

更新任务。

**参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 任务ID |
| title | string | 否 | 任务标题 |
| description | string | 否 | 任务描述 |
| assignee_id | int | 否 | 负责人ID |
| priority | string | 否 | 优先级 |

**返回值**: 任务对象

##### tasks.delete

删除任务。

**参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 任务ID |

**返回值**:

```json
{
  "success": true
}
```

##### tasks.update_status

更新任务状态。

**参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 任务ID |
| status | string | 是 | 新状态 (pending/in_progress/done/cancelled) |

**返回值**: 任务对象

---

#### 项目工具

##### projects.list

列出项目。

**参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| owner_id | int | 否 | 按所有者筛选 |

**返回值**:

```json
{
  "projects": [
    {
      "id": 1,
      "name": "Taskify 项目",
      "description": "主要产品开发",
      "owner_id": 123,
      "created_at": "2026-04-12T10:00:00Z"
    }
  ]
}
```

##### projects.get

获取项目详情。

**参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 项目ID |

**返回值**: 项目对象

##### projects.create

创建项目。

**参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 项目名称 |
| description | string | 否 | 项目描述 |

**返回值**: 项目对象

##### projects.update

更新项目。

**参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 项目ID |
| name | string | 否 | 项目名称 |
| description | string | 否 | 项目描述 |

**返回值**: 项目对象

##### projects.delete

删除项目。

**参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 项目ID |

**返回值**:

```json
{
  "success": true
}
```

---

#### 评论工具

##### comments.list

列出评论。

**参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| task_id | int | 是 | 任务ID |

**返回值**:

```json
{
  "comments": [
    {
      "id": 1,
      "content": "这个任务进展如何？",
      "task_id": 1,
      "user_id": 123,
      "created_at": "2026-04-12T10:00:00Z"
    }
  ]
}
```

##### comments.get

获取评论详情。

**参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 评论ID |

**返回值**: 评论对象

##### comments.create

创建评论。

**参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| task_id | int | 是 | 任务ID |
| content | string | 是 | 评论内容 |

**返回值**: 评论对象

##### comments.update

更新评论。

**参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 评论ID |
| content | string | 是 | 新评论内容 |

**返回值**: 评论对象

##### comments.delete

删除评论。

**参数**:

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 评论ID |

**返回值**:

```json
{
  "success": true
}
```

---

### MCP 资源 (Resources)

MCP 资源是只读数据，供客户端订阅。

#### tasks:///

任务列表资源 URI。

**响应格式**:

```json
{
  "uri": "tasks:///",
  "name": "任务列表",
  "mimeType": "application/json",
  "content": {
    "tasks": [...]
  }
}
```

#### projects:///

项目列表资源 URI。

**响应格式**:

```json
{
  "uri": "projects:///",
  "name": "项目列表",
  "mimeType": "application/json",
  "content": {
    "projects": [...]
  }
}
```

#### users:///

用户列表资源 URI（受隐私限制）。

**响应格式**:

```json
{
  "uri": "users:///",
  "name": "用户列表",
  "mimeType": "application/json",
  "content": {
    "users": [
      {
        "id": 123,
        "username": "zhangsan"
      }
    ]
  }
}
```

## 速率限制

| 端点 | 限制 |
|------|------|
| 所有 API 端点 | 100 请求/分钟/令牌 |
| POST /oauth/token (创建令牌) | 10 请求/分钟/用户 |

**速率限制响应** (429 Too Many Requests):

```json
{
  "error": "rate_limit_exceeded",
  "retry_after": 60
}
```

## 错误响应格式

所有错误响应遵循统一格式:

```json
{
  "error": "error_code",
  "message": "人类可读的错误消息",
  "details": {}  // 可选的额外信息
}
```

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
| token_invalid | 401 | 令牌无效 |
| token_expired | 401 | 令牌已过期 |
| token_revoked | 401 | 令牌已撤销 |
| scope_insufficient | 403 | 权限不足 |
| rate_limit_exceeded | 429 | 超过速率限制 |
| invalid_request | 400 | 请求参数无效 |
| not_found | 404 | 资源不存在 |
| internal_error | 500 | 服务器内部错误 |

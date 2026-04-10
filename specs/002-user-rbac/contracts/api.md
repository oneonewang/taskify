# Phase 1 API Contracts: 用户管理与 RBAC

## 认证接口

### POST /api/auth/register - 用户注册

**Request**:
```json
{
  "email": "user@example.com",
  "password": "SecurePass123",
  "display_name": "张三"
}
```

**Response 201**:
```json
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "display_name": "张三",
      "avatar_url": "",
      "email_verified": true
    },
    "token": "session_id_cookie"
  }
}
```

**Error Responses**:
| Code | Message | Condition |
|------|---------|-----------|
| 400 | 邮箱格式不正确 | Invalid email format |
| 400 | 密码至少8字符，需包含字母和数字 | Weak password |
| 409 | 邮箱已被注册 | Email already exists |

---

### POST /api/auth/login - 用户登录

**Request**:
```json
{
  "email": "user@example.com",
  "password": "SecurePass123"
}
```

**Response 200**:
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "display_name": "张三",
      "avatar_url": "",
      "email_verified": true
    }
  }
}
```
**Set-Cookie**: `session_id=xxx; HttpOnly; Path=/; MaxAge=86400`

**Error Responses**:
| Code | Message | Condition |
|------|---------|-----------|
| 401 | 邮箱或密码错误 | Invalid credentials |
| 403 | 账户已被禁用 | User disabled |

---

### POST /api/auth/logout - 用户登出

**Request**: (Cookie 中已包含 session)

**Response 200**:
```json
{
  "code": 0,
  "message": "登出成功"
}
```
**Set-Cookie**: `session_id=; HttpOnly; Path=/; MaxAge=0`

---

## 用户管理接口

### GET /api/users/me - 获取当前用户

**Auth**: Required

**Response 200**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "display_name": "张三",
    "avatar_url": "https://...",
    "email_verified": true,
    "last_login_at": "2026-04-10T10:00:00Z",
    "created_at": "2026-04-10T09:00:00Z"
  }
}
```

---

### PUT /api/users/me - 更新当前用户资料

**Auth**: Required

**Request**:
```json
{
  "display_name": "李四",
  "avatar_url": "https://example.com/avatar.png"
}
```

**Response 200**:
```json
{
  "code": 0,
  "message": "更新成功",
  "data": {
    "id": 1,
    "display_name": "李四",
    "avatar_url": "https://example.com/avatar.png"
  }
}
```

---

### PUT /api/users/me/password - 修改密码

**Auth**: Required

**Request**:
```json
{
  "current_password": "OldPass123",
  "new_password": "NewPass456"
}
```

**Response 200**:
```json
{
  "code": 0,
  "message": "密码修改成功"
}
```

**Error Responses**:
| Code | Message | Condition |
|------|---------|-----------|
| 400 | 新密码格式不符合要求 | Weak password |
| 401 | 当前密码错误 | Wrong current password |

---

## 角色管理接口

### GET /api/roles - 获取角色列表

**Auth**: Required (admin only)

**Response 200**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "roles": [
      {
        "id": 1,
        "name": "admin",
        "description": "系统管理员",
        "is_system": true,
        "scope": "system",
        "permissions": ["users.manage", "projects.manage", ...]
      },
      {
        "id": 2,
        "name": "owner",
        "description": "项目所有者",
        "is_system": true,
        "scope": "project",
        "permissions": ["projects.edit", "tasks.create", ...]
      }
    ]
  }
}
```

---

### POST /api/roles - 创建自定义角色

**Auth**: Required (admin only)

**Request**:
```json
{
  "name": "项目经理",
  "description": "可管理项目进度",
  "scope": "project",
  "permission_ids": [3, 4, 5, 6]
}
```

**Response 201**:
```json
{
  "code": 0,
  "message": "角色创建成功",
  "data": {
    "id": 5,
    "name": "项目经理",
    "description": "可管理项目进度",
    "is_system": false,
    "scope": "project"
  }
}
```

---

### PUT /api/roles/:id - 更新角色

**Auth**: Required (admin only)

**Request**:
```json
{
  "name": "高级项目经理",
  "description": "更新后的描述",
  "permission_ids": [3, 4, 5, 6, 7]
}
```

**Response 200**:
```json
{
  "code": 0,
  "message": "角色更新成功"
}
```

---

### DELETE /api/roles/:id - 删除角色

**Auth**: Required (admin only)

**Response 200**:
```json
{
  "code": 0,
  "message": "角色删除成功"
}
```

**Error Responses**:
| Code | Message | Condition |
|------|---------|-----------|
| 400 | 无法删除系统预定义角色 | Cannot delete system role |
| 400 | 角色已被分配给用户 | Role in use |

---

## 权限管理接口

### GET /api/permissions - 获取所有权限

**Auth**: Required (admin only)

**Response 200**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "permissions": [
      { "id": 1, "resource": "users", "action": "view", "description": "查看用户" },
      { "id": 2, "resource": "users", "action": "manage", "description": "管理用户" },
      { "id": 3, "resource": "projects", "action": "view", "description": "查看项目" },
      ...
    ]
  }
}
```

---

## 用户角色分配接口

### GET /api/users/:id/roles - 获取用户角色

**Auth**: Required (admin only)

**Response 200**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "system_roles": ["admin"],
    "project_memberships": [
      { "project_id": 1, "project_name": "项目A", "role": "owner" },
      { "project_id": 2, "project_name": "项目B", "role": "member" }
    ]
  }
}
```

---

### POST /api/users/:id/roles - 分配系统角色

**Auth**: Required (admin only)

**Request**:
```json
{
  "role_id": 1
}
```

**Response 200**:
```json
{
  "code": 0,
  "message": "角色分配成功"
}
```

---

### DELETE /api/users/:id/roles/:role_id - 移除用户角色

**Auth**: Required (admin only)

**Response 200**:
```json
{
  "code": 0,
  "message": "角色移除成功"
}
```

**Error Responses**:
| Code | Message | Condition |
|------|---------|-----------|
| 400 | 无法移除最后一个系统管理员 | Cannot remove last admin |

---

## 项目成员管理接口

### GET /api/projects/:id/members - 获取项目成员

**Auth**: Required (project member)

**Response 200**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "members": [
      {
        "user_id": 1,
        "display_name": "张三",
        "avatar_url": "...",
        "role": "owner",
        "joined_at": "2026-04-01T00:00:00Z"
      }
    ]
  }
}
```

---

### POST /api/projects/:id/members - 添加项目成员

**Auth**: Required (owner only)

**Request**:
```json
{
  "user_email": "newmember@example.com",
  "role_id": 3
}
```

**Response 201**:
```json
{
  "code": 0,
  "message": "成员添加成功",
  "data": {
    "user_id": 5,
    "role": "member"
  }
}
```

---

### PUT /api/projects/:id/members/:user_id - 更新成员角色

**Auth**: Required (owner only)

**Request**:
```json
{
  "role_id": 4
}
```

**Response 200**:
```json
{
  "code": 0,
  "message": "角色更新成功"
}
```

---

### DELETE /api/projects/:id/members/:user_id - 移除项目成员

**Auth**: Required (owner only)

**Response 200**:
```json
{
  "code": 0,
  "message": "成员移除成功"
}
```

---

## 审计日志接口

### GET /api/audit-logs - 获取审计日志

**Auth**: Required (admin only)

**Query Params**:
- `user_id` (optional): 筛选特定用户
- `event_type` (optional): 筛选事件类型
- `from` (optional): 开始时间 (ISO8601)
- `to` (optional): 结束时间 (ISO8601)
- `page` (default: 1): 页码
- `page_size` (default: 20): 每页数量

**Response 200**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "logs": [
      {
        "id": 1,
        "user_id": 1,
        "user_email": "admin@example.com",
        "event_type": "login",
        "details": "登录成功",
        "ip_address": "192.168.1.1",
        "created_at": "2026-04-10T10:00:00Z"
      }
    ],
    "pagination": {
      "total": 100,
      "page": 1,
      "page_size": 20
    }
  }
}
```

---

## 错误响应格式

所有 API 错误响应遵循统一格式：

```json
{
  "code": <错误码>,
  "message": "<错误消息>",
  "data": null
}
```

**错误码定义**：
| Code | HTTP Status | Description |
|------|-------------|-------------|
| 0 | 200/201 | 成功 |
| 400 | 400 | 请求参数错误 |
| 401 | 401 | 未认证 |
| 403 | 403 | 权限不足 |
| 404 | 404 | 资源不存在 |
| 409 | 409 | 资源冲突 |
| 500 | 500 | 服务器内部错误 |
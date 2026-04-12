# 数据模型文档

**功能**: MCP 服务与 PAT 认证
**更新日期**: 2026-04-12

## 实体关系图

```
┌─────────────┐       ┌─────────────────────┐
│    User     │ 1 ─── N│ PersonalAccessToken │
└─────────────┘       └─────────────────────┘
       │                        │
       │ 1                      │
       N                        │
       │                        │
       ▼                        ▼
┌─────────────┐        ┌─────────────────────┐
│   Task     │ N ─── 1│      Project         │
└─────────────┘        └─────────────────────┘
       │
       │ 1 ─── N
       ▼
┌─────────────┐
│   Comment   │
└─────────────┘
```

## 详细数据模型

### User (用户)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | uint | PK, AUTO | 用户ID |
| username | string | UNIQUE, NOT NULL | 用户名 |
| email | string | UNIQUE, NOT NULL | 邮箱 |
| password_hash | string | NOT NULL | 密码哈希 (bcrypt) |
| created_at | datetime | NOT NULL | 创建时间 |
| updated_at | datetime | NOT NULL | 更新时间 |

### PersonalAccessToken (个人访问令牌)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | uint | PK, AUTO | 令牌ID |
| user_id | uint | FK → User, NOT NULL | 所属用户 |
| token_hash | string | UNIQUE, NOT NULL, INDEX | 令牌 SHA-256 哈希 |
| token_prefix | string | NOT NULL | 令牌前缀 (显示用, 如 `tkf_a1b2`) |
| name | string | NOT NULL | 令牌名称 (用户指定) |
| scope | string | NOT NULL | 权限范围 (逗号分隔: `task:read,task:write,...`) |
| created_at | datetime | NOT NULL | 创建时间 |
| last_used_at | datetime | NULL | 最后使用时间 |
| expires_at | datetime | NOT NULL | 过期时间 |
| revoked_at | datetime | NULL | 撤销时间 (NULL = 未撤销) |

**状态机**:

```
活跃 (revoked_at=NULL AND expires_at > NOW)
    │
    ├──[revoked_at 设置]──→ 已撤销 (revoked_at != NULL)
    │
    └──[expires_at 到期]──→ 已过期 (expires_at <= NOW)
```

**索引**:
- `idx_token_hash` ON token_hash (UNIQUE)
- `idx_user_id` ON user_id
- `idx_expires_at` ON expires_at

### Project (项目)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | uint | PK, AUTO | 项目ID |
| name | string | NOT NULL | 项目名称 |
| description | string | NULL | 项目描述 |
| owner_id | uint | FK → User, NOT NULL | 所有者用户ID |
| created_at | datetime | NOT NULL | 创建时间 |
| updated_at | datetime | NOT NULL | 更新时间 |

**索引**:
- `idx_owner_id` ON owner_id

### Task (任务)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | uint | PK, AUTO | 任务ID |
| title | string | NOT NULL | 任务标题 |
| description | string | NULL | 任务描述 |
| status | string | NOT NULL, DEFAULT 'pending' | 任务状态 (pending/in_progress/done/cancelled) |
| priority | string | NOT NULL, DEFAULT 'medium' | 优先级 (low/medium/high/urgent) |
| project_id | uint | FK → Project | 所属项目ID |
| assignee_id | uint | FK → User | 负责人ID |
| creator_id | uint | FK → User, NOT NULL | 创建人ID |
| created_at | datetime | NOT NULL | 创建时间 |
| updated_at | datetime | NOT NULL | 更新时间 |

**索引**:
- `idx_project_id` ON project_id
- `idx_assignee_id` ON assignee_id
- `idx_creator_id` ON creator_id
- `idx_status` ON status

### Comment (评论)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | uint | PK, AUTO | 评论ID |
| content | string | NOT NULL | 评论内容 |
| task_id | uint | FK → Task, NOT NULL | 所属任务ID |
| user_id | uint | FK → User, NOT NULL | 评论用户ID |
| created_at | datetime | NOT NULL | 创建时间 |
| updated_at | datetime | NOT NULL | 更新时间 |

**索引**:
- `idx_task_id` ON task_id
- `idx_user_id` ON user_id

## 令牌作用域 (TokenScope)

作用域格式: `resource:action` 逗号分隔

| 资源 | 读操作 | 写操作 |
|------|--------|--------|
| task | task:read | task:write |
| project | project:read | project:write |
| comment | comment:read | comment:write |

**有效作用域组合示例**:
- `task:read` - 仅读取任务
- `task:read,task:write` - 读写任务
- `task:read,project:read,comment:read` - 只读所有资源
- `task:read,task:write,project:read,project:write,comment:read,comment:write` - 完全访问

## 验证规则

### PersonalAccessToken 创建时

| 字段 | 验证规则 |
|------|----------|
| name | 长度 1-100 字符，非空 |
| scope | 必须匹配 `^[a-z_:]+(,[a-z_:]+)*$` 格式 |
| expires_at | 必须晚于当前时间 |

### 令牌状态检查

| 状态 | 拒绝条件 |
|------|----------|
| revoked_at != NULL | 令牌已撤销 |
| expires_at <= NOW | 令牌已过期 |
| token_hash 不存在 | 无效令牌 |

## 假设

- 所有时间戳使用 UTC
- 令牌哈希使用 SHA-256
- 默认过期时间: 90 天 (可配置)

# 数据模型: Taskify

**日期**: 2026-04-08
**Branch**: 001-kanban-platform

---

## 实体关系图

```
┌─────────────┐       ┌─────────────┐       ┌─────────────┐
│    User     │       │   Project   │       │    Task     │
├─────────────┤       ├─────────────┤       ├─────────────┤
│ ID          │◄──────│ ID          │◄──────│ ID          │
│ Name        │       │ Name        │       │ Title       │
│ Role        │       │ Description │       │ Description │
│ Avatar      │       │ CreatedAt   │       │ Status      │
└─────────────┘       └─────────────┘       │ AssigneeID  │
      │                                       │ ProjectID   │
      │                                       │ Position    │
      │                                       │ CreatedAt   │
      │                                       │ UpdatedAt   │
      │                                       └─────────────┘
      │                                             │
      │               ┌─────────────┐               │
      │               │   Comment   │               │
      │               ├─────────────┤               │
      └──────────────►│ ID          │◄──────────────┘
                      │ Content     │
                      │ UserID      │
                      │ TaskID      │
                      │ CreatedAt   │
                      │ UpdatedAt   │
                      └─────────────┘
```

---

## 实体定义

### User (用户)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | uint | PK, AUTO | 用户ID |
| name | string | REQUIRED, MAX=50 | 用户姓名 |
| role | string | REQUIRED, ENUM | 角色: "product_manager", "engineer" |
| avatar | string | MAX=100 | 头像标识符(颜色代码) |
| created_at | datetime | AUTO | 创建时间 |

### Project (项目)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | uint | PK, AUTO | 项目ID |
| name | string | REQUIRED, MAX=100 | 项目名称 |
| description | string | MAX=500 | 项目描述 |
| created_at | datetime | AUTO | 创建时间 |

### Task (任务)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | uint | PK, AUTO | 任务ID |
| title | string | REQUIRED, MAX=200 | 任务标题 |
| description | string | MAX=2000 | 任务描述 |
| status | string | REQUIRED, ENUM | 状态: "todo", "in_progress", "review", "done" |
| position | int | DEFAULT=0 | 在列中的排序位置 |
| assignee_id | uint | FK(User.ID), REQUIRED | 负责人ID |
| project_id | uint | FK(Project.ID), REQUIRED | 所属项目ID |
| created_at | datetime | AUTO | 创建时间 |
| updated_at | datetime | AUTO | 更新时间 |

**状态枚举映射**:
- `todo` → 待办
- `in_progress` → 进行中
- `review` → 审核中
- `done` → 已完成

### Comment (评论)

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| id | uint | PK, AUTO | 评论ID |
| content | string | REQUIRED, MAX=2000 | 评论内容 |
| user_id | uint | FK(User.ID), REQUIRED | 评论者ID |
| task_id | uint | FK(Task.ID), REQUIRED | 所属任务ID |
| created_at | datetime | AUTO | 创建时间 |
| updated_at | datetime | AUTO | 更新时间 |

---

## 索引设计

| 表 | 索引类型 | 字段 | 用途 |
|----|----------|------|------|
| tasks | INDEX | project_id, status | 按项目和状态查询任务 |
| tasks | INDEX | assignee_id | 按负责人查询任务 |
| comments | INDEX | task_id | 按任务查询评论 |
| comments | INDEX | user_id | 按用户查询评论 |

---

## 数据验证规则

根据Taskify章程(安全优先):

1. **标题验证**: 非空，最大200字符，移除HTML标签
2. **描述验证**: 最大2000字符，允许纯文本
3. **评论验证**: 非空，最大2000字符
4. **状态验证**: 必须是预定义枚举值之一
5. **外键验证**: 负责人必须是存在的用户ID

---

## 预定义数据 (Seed Data)

### 用户 (5个)

| ID | Name | Role | Avatar |
|----|------|------|--------|
| 1 | 张明 | product_manager | #FF6B6B |
| 2 | 李伟 | engineer | #4ECDC4 |
| 3 | 王芳 | engineer | #45B7D1 |
| 4 | 刘强 | engineer | #96CEB4 |
| 5 | 陈静 | engineer | #FFEAA7 |

### 项目 (3个)

| ID | Name | Description |
|----|------|-------------|
| 1 | 任务管理重构 | 重构现有任务管理系统，提升性能和可维护性 |
| 2 | 移动端开发 | 开发iOS和Android移动应用 |
| 3 | API网关升级 | 升级API网关，支持更多协议和认证方式 |

### 看板列 (所有项目共用)

| Status | 中文名称 | 默认Position |
|--------|----------|--------------|
| todo | 待办 | 0 |
| in_progress | 进行中 | 1 |
| review | 审核中 | 2 |
| done | 已完成 | 3 |

---

## 迁移策略

开发环境(SQLite)到生产环境(PostgreSQL):

1. 使用Gorm AutoMigrate进行开发期 schema管理
2. 生产环境使用正式迁移工具(如golang-migrate)
3. 所有迁移必须向后兼容
4. 大字段添加使用默认值避免锁表

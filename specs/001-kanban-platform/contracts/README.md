# API 契约: Taskify REST API

**日期**: 2026-04-08
**Branch**: 001-kanban-platform

---

## 通用规范

### Base URL
```
开发环境: http://localhost:8080/api
```

### 统一响应格式

```typescript
// 成功响应
interface SuccessResponse<T> {
  success: true
  data: T
  message?: string
}

// 错误响应
interface ErrorResponse {
  success: false
  error: {
    code: string      // 错误码，如 "VALIDATION_ERROR"
    message: string   // 人类可读错误信息
    details?: any     // 详细错误信息
  }
}
```

### HTTP 状态码

| 状态码 | 含义 |
|--------|------|
| 200 | 成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

### 请求头

```
Content-Type: application/json
Accept: application/json
```

---

## 用户 API

### GET /api/users

获取所有用户列表。

**响应** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "张明",
      "role": "product_manager",
      "avatar": "#FF6B6B"
    },
    {
      "id": 2,
      "name": "李伟",
      "role": "engineer",
      "avatar": "#4ECDC4"
    }
  ]
}
```

---

## 项目 API

### GET /api/projects

获取所有项目列表。

**响应** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "任务管理重构",
      "description": "重构现有任务管理系统，提升性能和可维护性",
      "created_at": "2026-04-08T10:00:00Z"
    }
  ]
}
```

### GET /api/projects/:id

获取项目详情。

**路径参数**
| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 项目ID |

**响应** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "任务管理重构",
    "description": "重构现有任务管理系统，提升性能和可维护性",
    "created_at": "2026-04-08T10:00:00Z",
    "task_counts": {
      "todo": 5,
      "in_progress": 3,
      "review": 2,
      "done": 10
    }
  }
}
```

---

## 任务 API

### GET /api/projects/:id/tasks

获取项目的所有任务。

**路径参数**
| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 项目ID |

**查询参数**
| 参数 | 类型 | 说明 |
|------|------|------|
| status | string | 可选，按状态过滤 |

**响应** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "设计数据库Schema",
      "description": "设计任务和评论表结构",
      "status": "todo",
      "position": 0,
      "assignee": {
        "id": 2,
        "name": "李伟",
        "avatar": "#4ECDC4"
      },
      "project_id": 1,
      "created_at": "2026-04-08T10:00:00Z",
      "updated_at": "2026-04-08T10:00:00Z"
    }
  ]
}
```

### POST /api/projects/:id/tasks

创建新任务。

**路径参数**
| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 项目ID |

**请求体**
```json
{
  "title": "设计数据库Schema",      // 必填，最大200字符
  "description": "设计任务和评论表结构",  // 可选，最大2000字符
  "assignee_id": 2                  // 必填，用户ID
}
```

**响应** `201 Created`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "设计数据库Schema",
    "description": "设计任务和评论表结构",
    "status": "todo",
    "position": 0,
    "assignee": {
      "id": 2,
      "name": "李伟",
      "avatar": "#4ECDC4"
    },
    "project_id": 1,
    "created_at": "2026-04-08T10:00:00Z",
    "updated_at": "2026-04-08T10:00:00Z"
  }
}
```

**错误响应** `400 Bad Request`
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "标题不能为空"
  }
}
```

### PUT /api/tasks/:id

更新任务信息。

**路径参数**
| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 任务ID |

**请求体**
```json
{
  "title": "更新后的标题",          // 可选
  "description": "更新后的描述",    // 可选
  "assignee_id": 3                 // 可选，新的负责人ID
}
```

**响应** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "更新后的标题",
    "description": "更新后的描述",
    "status": "todo",
    "position": 0,
    "assignee": {
      "id": 3,
      "name": "王芳",
      "avatar": "#45B7D1"
    },
    "project_id": 1,
    "updated_at": "2026-04-08T11:00:00Z"
  }
}
```

### PUT /api/tasks/:id/status

更新任务状态(看板拖放)。

**路径参数**
| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 任务ID |

**请求体**
```json
{
  "status": "in_progress",         // 必填，状态枚举
  "position": 2                     // 必填，新位置(0-based)
}
```

**响应** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "status": "in_progress",
    "position": 2,
    "updated_at": "2026-04-08T11:00:00Z"
  }
}
```

**状态枚举值**
| 值 | 中文 |
|----|------|
| todo | 待办 |
| in_progress | 进行中 |
| review | 审核中 |
| done | 已完成 |

### DELETE /api/tasks/:id

删除任务。

**路径参数**
| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 任务ID |

**响应** `200 OK`
```json
{
  "success": true,
  "message": "任务已删除"
}
```

---

## 评论 API

### GET /api/tasks/:id/comments

获取任务的评论列表。

**路径参数**
| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 任务ID |

**响应** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "content": "这个设计看起来不错",
      "user": {
        "id": 3,
        "name": "王芳",
        "avatar": "#45B7D1"
      },
      "task_id": 1,
      "created_at": "2026-04-08T10:30:00Z",
      "updated_at": "2026-04-08T10:30:00Z"
    }
  ]
}
```

### POST /api/tasks/:id/comments

添加评论。

**路径参数**
| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 任务ID |

**请求头**
```
X-User-ID: 2    // 当前登录用户ID (简单认证)
```

**请求体**
```json
{
  "content": "这个设计看起来不错"    // 必填，最大2000字符
}
```

**响应** `201 Created`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "content": "这个设计看起来不错",
    "user": {
      "id": 2,
      "name": "李伟",
      "avatar": "#4ECDC4"
    },
    "task_id": 1,
    "created_at": "2026-04-08T10:30:00Z",
    "updated_at": "2026-04-08T10:30:00Z"
  }
}
```

### PUT /api/tasks/:id/comments/:cid

编辑评论。

**路径参数**
| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 任务ID |
| cid | integer | 评论ID |

**请求头**
```
X-User-ID: 2    // 当前登录用户ID
```

**请求体**
```json
{
  "content": "更新后的评论内容"    // 必填，最大2000字符
}
```

**响应** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": 1,
    "content": "更新后的评论内容",
    "user": {
      "id": 2,
      "name": "李伟",
      "avatar": "#4ECDC4"
    },
    "task_id": 1,
    "created_at": "2026-04-08T10:30:00Z",
    "updated_at": "2026-04-08T11:00:00Z"
  }
}
```

**错误响应** `403 Forbidden` (不是评论作者)
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "只能编辑自己的评论"
  }
}
```

### DELETE /api/tasks/:id/comments/:cid

删除评论。

**路径参数**
| 参数 | 类型 | 说明 |
|------|------|------|
| id | integer | 任务ID |
| cid | integer | 评论ID |

**请求头**
```
X-User-ID: 2    // 当前登录用户ID
```

**响应** `200 OK`
```json
{
  "success": true,
  "message": "评论已删除"
}
```

**错误响应** `403 Forbidden` (不是评论作者)
```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "只能删除自己的评论"
  }
}
```

---

## 实时事件 API (SSE)

### GET /api/events

建立SSE连接，接收实时事件。

**请求头**
```
Accept: text/event-stream
```

**响应** `200 OK` (流式)

```
event: task_created
data: {"id":1,"title":"新任务","status":"todo","project_id":1}

event: task_updated
data: {"id":1,"title":"更新后的任务","status":"in_progress"}

event: task_deleted
data: {"id":1}

event: task_moved
data: {"id":1,"from":"todo","to":"in_progress","position":2}

event: comment_added
data: {"id":1,"task_id":1,"content":"新评论","user_id":2}
```

**连接管理**
- 心跳间隔: 30秒
- 自动重连: EventSource API 自动处理
- 超时: 5分钟无活动自动断开

---

## 错误码

| 错误码 | HTTP状态 | 说明 |
|--------|----------|------|
| VALIDATION_ERROR | 400 | 请求参数验证失败 |
| UNAUTHORIZED | 401 | 缺少认证信息 |
| FORBIDDEN | 403 | 无权限操作 |
| NOT_FOUND | 404 | 资源不存在 |
| INTERNAL_ERROR | 500 | 服务器内部错误 |

---

## 版本历史

| 版本 | 日期 | 变更 |
|------|------|------|
| v1.0 | 2026-04-08 | 初始版本 |

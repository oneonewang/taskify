# Taskify 看板平台

团队生产力平台，支持项目看板、任务管理、团队协作和实时更新。

## 技术栈

**后端**: Go 1.21+ | Gin | Gorm | SQLite (开发) / PostgreSQL (生产)

**前端**: Vue.js 3 | Element Plus | Pinia | vue-draggable-plus | TypeScript

## 快速启动

### 前置要求

- Go 1.21+
- Node.js 18+
- npm / pnpm

### 后端启动

```bash
cd backend

# 安装依赖
go mod tidy

# 启动开发服务器 (SQLite)
go run cmd/server/main.go
```

服务器运行在 http://localhost:8080

### 前端启动

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端运行在 http://localhost:5173

## 项目结构

```
taskify/
├── backend/                    # Go 后端
│   ├── cmd/server/main.go      # 入口点
│   ├── internal/
│   │   ├── models/             # 数据模型
│   │   ├── handlers/           # HTTP 处理器
│   │   ├── services/           # 业务逻辑
│   │   ├── repository/         # 数据访问
│   │   └── middleware/         # 中间件
│   ├── config/                 # 配置管理
│   └── pkg/response/           # 统一响应
│
├── frontend/                   # Vue.js 前端
│   ├── src/
│   │   ├── api/               # API 客户端
│   │   ├── components/        # Vue 组件
│   │   ├── composables/       # 组合式函数
│   │   ├── stores/            # Pinia 状态管理
│   │   └── views/             # 页面
│   └── ...
│
└── specs/                      # 需求和设计文档
```

## 功能特性

- [x] 用户身份选择 (US0)
- [x] 查看项目和看板 (US1)
- [x] 查看团队成员 (US2)
- [x] 创建和分配任务 (US3)
- [x] 移动任务状态 (US4)
- [x] 高亮我的任务 (US5)
- [x] 添加任务评论 (US6)
- [x] 查看任务详情 (US7)
- [x] 编辑我的评论 (US8)
- [x] 删除我的评论 (US9)
- [x] 实时更新 (SSE) - T068-T071

## 用户管理与 RBAC

- [x] 用户注册与登录（会话认证）
- [x] 管理员后台（用户管理、角色管理）
- [x] 系统角色（超级管理员、管理员、普通用户）
- [x] 项目角色（所有者、成员、访客）
- [x] 项目成员管理（添加、移除、角色分配）
- [x] 基于角色的权限控制

## API 端点

### 认证

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/auth/register | 用户注册 |
| POST | /api/auth/login | 用户登录 |
| POST | /api/auth/logout | 用户登出 |

### 用户

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/users/me | 获取当前用户 |
| PUT | /api/users/me | 更新个人信息 |
| PUT | /api/users/me/password | 修改密码 |

### 项目

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/projects | 获取项目列表 |
| GET | /api/projects/:id | 获取项目详情 |
| POST | /api/projects | 创建项目 |
| PUT | /api/projects/:id | 更新项目（所有者） |
| DELETE | /api/projects/:id | 删除项目（所有者） |
| POST | /api/projects/:id/archive | 归档项目（所有者） |
| GET | /api/projects/:id/members | 获取项目成员 |
| POST | /api/projects/:id/members | 添加项目成员（所有者） |
| PUT | /api/projects/:id/members/:uid | 更新成员角色（所有者） |
| DELETE | /api/projects/:id/members/:uid | 移除项目成员（所有者） |

### 任务

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/projects/:id/tasks | 获取项目任务 |
| POST | /api/projects/:id/tasks | 创建任务（成员） |
| PUT | /api/tasks/:id | 更新任务 |
| PUT | /api/tasks/:id/status | 更新任务状态 (拖放) |
| DELETE | /api/tasks/:id | 删除任务 |

### 评论

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/tasks/:id/comments | 获取评论列表 |
| POST | /api/tasks/:id/comments | 添加评论 |
| PUT | /api/tasks/:id/comments/:cid | 编辑评论 |
| DELETE | /api/tasks/:id/comments/:cid | 删除评论 |

### 管理员

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/admin/users | 获取用户列表 |
| GET | /api/admin/roles | 获取角色列表 |
| POST | /api/admin/roles | 创建角色 |
| PUT | /api/admin/roles/:id | 更新角色 |
| DELETE | /api/admin/roles/:id | 删除角色 |
| PUT | /api/admin/roles/:id/permissions | 设置角色权限 |
| GET | /api/admin/permissions | 获取权限列表 |
| GET | /api/admin/users/:id/roles | 获取用户角色 |
| POST | /api/admin/users/:id/roles | 分配系统角色 |
| DELETE | /api/admin/users/:id/roles/:rid | 移除系统角色 |

### 实时

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/events | SSE 实时事件流 |

## 开发指南

### 添加新模型

1. 在 `backend/internal/models/` 创建模型文件
2. 在 `backend/internal/repository/sqlite.go` 的 AutoMigrate 中注册
3. 在 handler 中添加 CRUD 端点

### 添加新前端组件

1. 在 `frontend/src/components/` 创建 Vue 组件
2. 在 Pinia store 中添加状态管理
3. 在 API 客户端中添加后端调用

## 配置

### 后端环境变量

参考 `backend/.env.example`:

```bash
PORT=8080
DATABASE_URL=              # 空值 = SQLite
GIN_MODE=debug             # debug 或 release
```

### 前端环境变量

参考 `frontend/.env.example`:

```bash
# Vite 开发环境使用代理，无需配置
```

## 调试

### 查看 SSE 事件

```bash
curl -N http://localhost:8080/api/events
```

### 运行测试

```bash
# 后端
cd backend
go test -v ./...

# 前端
cd frontend
npm run dev
```

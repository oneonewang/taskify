# Taskify 快速入门

**日期**: 2026-04-08
**Branch**: 001-kanban-platform

---

## 环境要求

| 组件 | 版本 | 说明 |
|------|------|------|
| Go | 1.21+ | 后端运行环境 |
| Node.js | 18+ | 前端构建工具 |
| npm/pnpm | 最新 | 包管理器 |
| SQLite | - | 开发数据库(内置) |

---

## 后端启动 (开发环境)

```bash
# 进入后端目录
cd backend

# 安装依赖
go mod tidy

# 启动开发服务器 (SQLite)
go run cmd/server/main.go

# 服务器运行在 http://localhost:8080
```

### 环境变量配置

```bash
# .env 文件 (backend/)
GIN_MODE=debug
DATABASE_URL=taskify.db        # SQLite (开发)
# DATABASE_URL=postgres://...  # PostgreSQL (生产)
PORT=8080
```

---

## 前端启动 (开发环境)

```bash
# 新开终端，进入前端目录
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 前端运行在 http://localhost:5173
```

### 主要脚本

```bash
npm run dev      # 开发服务器 with HMR
npm run build    # 生产构建
npm run preview  # 预览生产构建
npm run test     # 运行测试
npm run lint     # 代码检查
```

---

## 目录结构

```
taskify/
├── backend/                    # Go后端
│   ├── cmd/server/main.go      # 入口点
│   ├── internal/
│   │   ├── models/             # 数据模型
│   │   ├── handlers/           # HTTP处理器
│   │   ├── services/           # 业务逻辑
│   │   ├── repository/         # 数据访问
│   │   └── middleware/         # 中间件
│   └── pkg/response/           # 统一响应
│
├── frontend/                   # Vue.js前端
│   ├── src/
│   │   ├── api/               # API客户端
│   │   ├── components/        # 组件
│   │   ├── stores/            # Pinia状态
│   │   └── views/             # 页面
│   └── ...
│
└── specs/                     # 文档
```

---

## API 端点概览

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/projects | 获取项目列表 |
| GET | /api/projects/:id | 获取项目详情 |
| GET | /api/projects/:id/tasks | 获取项目任务 |
| POST | /api/projects/:id/tasks | 创建任务 |
| PUT | /api/tasks/:id | 更新任务 |
| PUT | /api/tasks/:id/status | 更新任务状态(拖放) |
| DELETE | /api/tasks/:id | 删除任务 |
| GET | /api/tasks/:id/comments | 获取评论列表 |
| POST | /api/tasks/:id/comments | 添加评论 |
| PUT | /api/tasks/:id/comments/:cid | 编辑评论 |
| DELETE | /api/tasks/:id/comments/:cid | 删除评论 |
| GET | /api/users | 获取用户列表 |
| GET | /api/events | SSE实时事件流 |

详细API契约见 `contracts/` 目录。

---

## 常用开发任务

### 添加新模型

1. 在 `backend/internal/models/` 创建模型文件
2. 运行 `go generate` 更新 GORM 关联
3. 在 handler 中添加 CRUD 端点
4. 添加单元测试

### 添加新前端组件

1. 在 `frontend/src/components/` 创建 Vue 组件
2. 在 Pinia store 中添加状态管理
3. 在 API client 中添加后端调用
4. 使用 Element Plus 组件库

### 调试实时更新

```bash
# 查看SSE事件
curl -N http://localhost:8080/api/events
```

---

## 测试

### 后端测试

```bash
cd backend
go test -v ./...
go test -cover ./...
```

### 前端测试

```bash
cd frontend
npm run test        # Vitest单元测试
npm run test:e2e    # Playwright E2E测试
```

---

## 下一步

1. 阅读 [数据模型](data-model.md) 了解数据库设计
2. 阅读 `contracts/` 目录了解API契约
3. 阅读 [研究文档](research.md) 了解技术选型
4. 查看 `tasks.md` 获取任务列表

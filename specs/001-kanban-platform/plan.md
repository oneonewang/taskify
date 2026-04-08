# Implementation Plan: Taskify 看板平台

**Branch**: `001-kanban-platform` | **Date**: 2026-04-08 | **Spec**: [spec.md](./spec.md)

## Summary

构建Taskify团队生产力平台：Go+Gin+Gorm后端REST API（项目/任务/通知），Vue.js+Element Plus前端看板（拖放+实时更新），SQLite开发/PostgreSQL生产。

## Technical Context

**Language/Version**: Go 1.21+ (最新稳定版)
**Primary Dependencies**:
- Backend: Gin (web framework), Gorm (ORM), go-playwright (E2E testing), SSE
- Frontend: Vue.js 3, Element Plus, vue-draggable-plus (drag-and-drop), Pinia
**Storage**: SQLite (开发环境), PostgreSQL (生产环境)
**Testing**: go test + httptest (后端), Vitest + Playwright (前端)
**Target Platform**: Linux服务器 (后端), 桌面浏览器 (前端)
**Project Type**: Web应用 (REST API后端 + 单页应用前端)
**Performance Goals**: <100ms API响应 (p95), 60fps 看板拖放
**Constraints**: 本地存储会话持久化 (刷新保持, 关闭重置), 无密码认证(预定义用户选择)
**Scale/Scope**: 5用户, 3项目, 每项目4看板列, 预计50+任务

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Gate | Status | Notes |
|------|--------|-------|
| 安全优先 (输入验证) | ✅ PASS | 所有用户输入须验证 |
| 微服务架构 | ⚠️ VIOLATION | 单体应用 - 章程要求微服务，但spec未要求拆分 |
| 完全文档化 | ✅ PASS | 所有代码需中文注释 |
| 中文内容 | ✅ PASS | 所有文档使用简体中文 |

**Violation Justification (微服务)**:
- Spec明确要求单体看板应用，功能边界清晰（项目/任务/通知）
- 小团队(5人)不需要服务拆分带来的运维复杂度
- 后续可通过模块化设计为微服务迁移做准备
- 替代方案(立即拆分)不适合v1.0交付目标

### Post-Design Constitution Check

| Gate | Status | Notes |
|------|--------|-------|
| 安全优先 | ✅ PASS | 研究确认使用参数化查询(Gorm)，防止SQL注入 |
| 微服务架构 | ⚠️ VIOLATION | 维持 - 单体架构，violation已论证 |
| 完全文档化 | ✅ PASS | 所有代码文件需包含中文注释 |
| 中文内容 | ✅ PASS | 文档和注释均使用简体中文 |
| 输入验证 | ✅ PASS | API契约定义了所有验证规则 |

## Project Structure

### Documentation (this feature)

```text
specs/001-kanban-platform/
├── plan.md              # 本文件
├── research.md          # Phase 0 研究成果
├── data-model.md        # Phase 1 数据模型
├── quickstart.md        # Phase 1 快速入门
├── contracts/           # Phase 1 接口契约
│   └── README.md
└── tasks.md             # Phase 2 任务列表 (/speckit.tasks 命令生成)
```

### Source Code (仓库根目录)

```text
d:/work/taskify/
├── backend/                    # Go后端服务
│   ├── cmd/
│   │   └── server/              # 主入口
│   │       └── main.go
│   ├── internal/
│   │   ├── models/              # 数据模型
│   │   │   ├── user.go
│   │   │   ├── project.go
│   │   │   ├── task.go
│   │   │   └── comment.go
│   │   ├── handlers/             # HTTP处理器
│   │   │   ├── project.go
│   │   │   ├── task.go
│   │   │   └── notification.go
│   │   ├── services/             # 业务逻辑
│   │   │   └── kanban.go
│   │   ├── repository/           # 数据访问
│   │   │   └── sqlite.go
│   │   └── middleware/          # 中间件
│   │       └── cors.go
│   ├── pkg/
│   │   └── response/            # 统一响应
│   │       └── response.go
│   ├── config/
│   │   └── config.go            # 配置管理
│   ├── go.mod
│   └── go.sum
├── frontend/                    # Vue.js前端
│   ├── src/
│   │   ├── api/                 # API客户端
│   │   │   ├── project.ts
│   │   │   ├── task.ts
│   │   │   └── notification.ts
│   │   ├── components/
│   │   │   ├── kanban/          # 看板组件
│   │   │   │   ├── KanbanBoard.vue
│   │   │   │   ├── KanbanColumn.vue
│   │   │   │   └── KanbanCard.vue
│   │   │   └── task/            # 任务组件
│   │   │       └── TaskDetail.vue
│   │   ├── stores/              # 状态管理 (Pinia)
│   │   │   ├── user.ts
│   │   │   ├── project.ts
│   │   │   └── kanban.ts
│   │   ├── views/               # 页面
│   │   │   ├── UserSelect.vue
│   │   │   └── ProjectKanban.vue
│   │   ├── router/
│   │   │   └── index.ts
│   │   ├── App.vue
│   │   └── main.ts
│   ├── index.html
│   ├── vite.config.ts
│   ├── package.json
│   └── tsconfig.json
├── specs/                       # 文档
├── CLAUDE.md                    # 项目章程
└── README.md
```

**Structure Decision**: 采用选项2(前后端分离Web应用)，backend与frontend完全解耦，通过REST API通信。

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| 微服务架构违规(单体应用) | Spec明确要求单体应用，v1.0交付目标明确 | 拆分为微服务需额外API网关、服务发现、数据库隔离，适合v2.0 |

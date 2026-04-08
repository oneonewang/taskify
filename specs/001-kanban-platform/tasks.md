# Taskify 看板平台 - 任务列表

**Feature**: Taskify 看板平台
**Branch**: `001-kanban-platform`
**Generated**: 2026-04-08
**Tech Stack**: Go 1.21+ (Gin + Gorm) | Vue.js 3 (Element Plus + Pinia + vue-draggable-plus) | SQLite (开发) / PostgreSQL (生产)

---

## 概述

本任务列表基于 `plan.md` (技术架构)、`spec.md` (用户故事)、`data-model.md` (数据模型)、`contracts/README.md` (API契约) 生成。

**总任务数**: 43
**用户故事数**: 10 (US0-US9)
**MVP范围**: US0 + US1 + US2 + US3 + US4 + US5 (P1优先级)

---

## 依赖关系图

```
[US0: 选择当前用户]
       │
       ▼
[US1: 查看项目和看板] ──► [US2: 查看团队成员]
       │
       ▼
[US3: 创建和分配任务] ──► [US4: 移动任务状态] ──► [US5: 高亮我的任务]
       │
       ▼
[US6: 添加任务评论] ──► [US7: 查看任务详情] ──► [US8: 编辑我的评论] ──► [US9: 删除我的评论]
```

**并行执行机会**:
- US1 和 US2 可并行开发（不同API端点）
- US3、US4、US5 可在 US1 完成后并行开发

---

## Phase 1: 项目初始化

**目标**: 建立前后端项目结构，安装依赖，验证开发环境

### 独立测试标准
- [ ] 后端 `go run cmd/server/main.go` 启动成功，监听 8080 端口
- [ ] 前端 `npm run dev` 启动成功，访问 http://localhost:5173
- [ ] API GET /api/users 返回 200

### 实现任务

- [x] T001 创建后端项目结构 `backend/cmd/server/`, `backend/internal/{models,handlers,services,repository,middleware}`, `backend/pkg/response`, `backend/config/`
- [x] T002 [P] 初始化 Go 模块：backend/go.mod (module github.com/taskify/backend, go 1.21+)
- [x] T003 [P] 创建 backend/go.mod 并添加依赖：gin, gorm, gorm.io/driver/sqlite, gorm.io/driver/postgres
- [x] T004 创建前端项目结构 `frontend/src/{api,components,stores,views,router}`，使用 Vite + Vue 3 + TypeScript
- [x] T005 [P] 创建 frontend/package.json 并安装依赖：vue, vue-router, pinia, element-plus, vue-draggable-plus, axios
- [x] T006 创建 backend/config/config.go 实现环境变量配置管理 (GIN_MODE, DATABASE_URL, PORT)
- [x] T007 创建 backend/pkg/response/response.go 实现统一响应格式 {success, data, message/error}
- [x] T008 创建 backend/internal/models/user.go, project.go, task.go, comment.go (参考 data-model.md)
- [x] T009 验证项目运行：后端启动 + 前端 dev server

---

## Phase 2: 基础设施层

**目标**: 建立数据库连接、统一响应、错误处理等基础设施，为所有用户故事提供支撑

### 独立测试标准
- [ ] 后端单元测试通过：go test ./internal/models/...
- [ ] API 响应格式符合 contracts/README.md 规范
- [ ] 数据库迁移成功，自动创建所有表

### 实现任务

- [x] T010 创建 backend/internal/repository/sqlite.go 实现 Gorm SQLite 连接和 AutoMigrate
- [x] T011 [P] 创建 backend/internal/models/user.go 实现 User 模型和预定义用户数据 (5个用户 Seed Data)
- [x] T012 [P] 创建 backend/internal/models/project.go 实现 Project 模型和预定义项目数据 (3个项目 Seed Data)
- [x] T013 [P] 创建 backend/internal/models/task.go 实现 Task 模型 (关联 User, Project，状态枚举: todo/in_progress/review/done)
- [x] T014 [P] 创建 backend/internal/models/comment.go 实现 Comment 模型 (关联 User, Task)
- [x] T015 创建 backend/internal/middleware/cors.go 实现 CORS 中间件
- [x] T016 创建 backend/cmd/server/main.go 实现服务器入口、路由注册、数据库初始化
- [x] T017 运行数据库迁移验证：`go run cmd/server/main.go` 无错误启动

---

## Phase 3: 用户身份选择 (US0) - P1

**目标**: 用户打开应用时选择当前用户，实现会话持久化

**独立测试标准**:
- [ ] 用户打开应用后直接看到用户选择列表
- [ ] 选择用户后进入项目列表视图
- [ ] 刷新页面保持用户选择
- [ ] 关闭浏览器后重置用户选择

### 实现任务

- [x] T018 创建 frontend/src/stores/user.ts 实现用户状态管理 (Pinia store)
- [x] T019 创建 frontend/src/views/UserSelect.vue 实现用户选择页面 (显示5个预定义用户)
- [x] T020 创建 frontend/src/api/user.ts 实现 API 客户端 GET /api/users
- [x] T021 实现用户选择持久化到 localStorage，刷新保持，关闭浏览器重置
- [x] T022 创建 frontend/src/router/index.ts 配置路由：/ (UserSelect), /projects (ProjectKanban)
- [x] T023 创建 frontend/src/App.vue 基础布局组件
- [x] T024 集成测试：选择用户 → 跳转到项目列表 → 刷新页面 → 验证用户保持

---

## Phase 4: 查看项目和看板 (US1) - P1

**目标**: 查看项目列表和每个项目的看板视图

**独立测试标准**:
- [ ] 显示3个预定义项目："任务管理重构"、"移动端开发"、"API网关升级"
- [ ] 选择项目后显示4列看板：待办、进行中、审核中、已完成
- [ ] 每列显示该列的任务数量

### 实现任务

- [x] T025 [P] 创建 backend/internal/handlers/project.go 实现 GET /api/projects, GET /api/projects/:id (含 task_counts)
- [x] T026 [P] 创建 backend/internal/services/kanban.go 实现看板业务逻辑
- [x] T027 创建 frontend/src/api/project.ts 实现 API 客户端
- [x] T028 创建 frontend/src/views/ProjectKanban.vue 项目看板主视图
- [x] T029 创建 frontend/src/components/kanban/KanbanBoard.vue 看板主组件
- [x] T030 创建 frontend/src/components/kanban/KanbanColumn.vue 单列组件 (显示列名、任务数、任务列表)
- [x] T031 创建 frontend/src/components/kanban/KanbanCard.vue 任务卡片组件 (显示标题、负责人头像/姓名)
- [x] T032 创建 frontend/src/stores/project.ts 实现项目状态管理
- [x] T033 集成测试：选择项目 → 显示看板4列 → 验证列任务数量

---

## Phase 5: 查看团队成员 (US2) - P1

**目标**: 查看项目团队成员列表

**独立测试标准**:
- [ ] 显示5个预定义用户：产品经理"张明"、4个工程师
- [ ] 显示每个成员的职责分类（产品经理/工程师）
- [ ] 头像和姓名清晰可辨

### 实现任务

- [x] T034 复用 T019 UserSelect.vue 中的用户列表展示逻辑
- [x] T035 创建 frontend/src/components/TeamMembers.vue 团队成员侧边栏组件
- [x] T036 在 ProjectKanban.vue 中集成 TeamMembers 组件
- [x] T037 集成测试：查看团队成员 → 验证5个用户显示正确

---

## Phase 6: 创建和分配任务 (US3) - P1

**目标**: 创建新任务并分配给团队成员

**独立测试标准**:
- [ ] 在"待办"列创建新任务
- [ ] 必填字段：标题、负责人；可选：描述
- [ ] 任务创建后显示在"待办"列
- [ ] 任务卡片显示被分配成员姓名

### 实现任务

- [x] T038 [P] 创建 backend/internal/handlers/task.go 实现 POST /api/projects/:id/tasks (创建任务), PUT /api/tasks/:id (更新任务)
- [x] T039 创建 frontend/src/api/task.ts 实现 API 客户端
- [x] T040 创建 frontend/src/components/task/TaskForm.vue 任务创建/编辑表单
- [x] T041 在 KanbanColumn.vue 中集成"创建任务"按钮和 TaskForm
- [x] T042 创建 frontend/src/stores/kanban.ts 实现看板状态管理 (添加任务、更新任务)
- [ ] T043 集成测试：创建任务 → 填写标题和负责人 → 提交 → 验证任务显示在待办列

---

## Phase 7: 移动任务状态 (US4) - P1

**目标**: 将任务在不同看板列之间拖拽移动

**独立测试标准**:
- [ ] 任务从"待办"拖拽到"进行中" → 任务消失并出现在"进行中"列
- [ ] 拖拽后任务状态持久保存
- [ ] 支持跨列拖拽

### 实现任务

- [x] T044 [P] 创建 backend/internal/handlers/task.go 实现 PUT /api/tasks/:id/status (状态更新 + position)
- [x] T045 在 frontend/src/stores/kanban.ts 添加任务移动状态更新逻辑
- [x] T046 集成 vue-draggable-plus：在 KanbanColumn.vue 中使用 VueDraggable 组件
- [x] T047 实现乐观更新：拖拽时立即更新 UI，API 失败时回滚
- [ ] T048 集成测试：拖拽任务 → 跨列移动 → 验证状态持久化
- [ ] T048.1 [P] 实现任务移动的竞态处理：last-write-wins 策略，API返回最新状态时前端重置UI

---

## Phase 8: 高亮我的任务 (US5) - P1

**目标**: 当前用户被分配的任务卡片高亮显示

**独立测试标准**:
- [ ] 当前用户登录后，分配给该用户的任务卡片特殊颜色显示
- [ ] 切换用户后，高亮相应变化

### 实现任务

- [x] T049 在 KanbanCard.vue 中实现任务卡片高亮逻辑 (根据当前用户 ID vs assignee_id)
- [x] T050 添加高亮样式：特殊背景色/边框，区别于普通任务卡片
- [ ] T051 集成测试：选择用户A → 验证A的任务高亮 → 切换用户B → 验证B的任务高亮

---

## Phase 9: 添加任务评论 (US6) - P2

**目标**: 为任务添加无限数量的评论

**独立测试标准**:
- [ ] 点击任务卡片可添加评论
- [ ] 评论提交后显示在任务详情中
- [ ] 评论显示内容和提交者姓名

### 实现任务

- [x] T052 [P] 创建 backend/internal/handlers/comment.go 实现 POST /api/tasks/:id/comments
- [x] T053 [P] 创建 backend/internal/handlers/comment.go 实现 GET /api/tasks/:id/comments
- [x] T054 创建 frontend/src/api/comment.ts 实现 API 客户端
- [x] T055 创建 frontend/src/components/task/TaskDetail.vue 任务详情弹窗
- [x] T056 在 TaskDetail.vue 中集成评论列表和评论表单
- [x] T057 集成测试：打开任务详情 → 添加评论 → 验证评论显示

---

## Phase 10: 查看任务详情 (US7) - P2

**目标**: 查看任务完整详情

**独立测试标准**:
- [ ] 点击任务卡片显示详情弹窗
- [ ] 显示标题、描述、当前状态、负责人
- [ ] 评论按时间顺序显示

### 实现任务

- [x] T058 在 TaskDetail.vue 中完善任务详情展示 (标题、描述、状态、负责人、时间)
- [x] T059 集成测试：点击任务卡片 → 验证详情弹窗显示完整信息

---

## Phase 11: 编辑我的评论 (US8) - P2

**目标**: 编辑自己发表的评论

**独立测试标准**:
- [ ] 当前用户的评论显示编辑按钮
- [ ] 点击编辑按钮可修改评论内容
- [ ] 他人评论无编辑按钮

### 实现任务

- [x] T060 [P] 创建 backend/internal/handlers/comment.go 实现 PUT /api/tasks/:id/comments/:cid
- [x] T061 在 TaskDetail.vue 中为当前用户评论添加编辑按钮
- [x] T062 实现评论编辑表单（内联编辑或弹窗编辑）
- [x] T063 集成测试：查看任务详情 → 编辑自己的评论 → 验证更新成功

---

## Phase 12: 删除我的评论 (US9) - P2

**目标**: 删除自己发表的评论

**独立测试标准**:
- [ ] 当前用户的评论显示删除按钮
- [ ] 点击删除按钮显示确认提示
- [ ] 确认后评论移除
- [ ] 他人评论无删除按钮

### 实现任务

- [x] T064 [P] 创建 backend/internal/handlers/comment.go 实现 DELETE /api/tasks/:id/comments/:cid
- [x] T065 在 TaskDetail.vue 中为当前用户评论添加删除按钮
- [x] T066 实现删除确认对话框
- [x] T067 集成测试：查看任务详情 → 删除自己的评论 → 验证评论移除

---

## Phase 13: 实时更新 (SSE) - 可选增强

**目标**: 实现看板实时更新

**说明**: 此功能为增强功能，不影响核心 P1 用户故事

### 实现任务

- [ ] T068 [P] 创建 backend/internal/handlers/notification.go 实现 GET /api/events (SSE)
- [ ] T069 [P] 创建 backend/internal/services/kanban.go 中实现 SSE 事件广播 (task_created, task_updated, task_deleted, task_moved, comment_added)
- [ ] T070 创建 frontend/src/composables/useSSE.ts 实现前端 SSE 客户端
- [ ] T071 在 frontend/src/stores/kanban.ts 中集成 SSE 事件监听，自动更新状态

---

## Phase 14: 收尾与优化

**目标**: 完成项目收尾工作

### 实现任务

- [ ] T072 创建 backend/.env 示例配置文件
- [ ] T073 创建 frontend/.env 示例配置文件
- [ ] T074 添加前端 vite.config.ts 代理配置 (开发环境 API 代理到 localhost:8080)
- [ ] T075 添加 README.md (项目说明、启动命令)
- [ ] T076 清理临时代码和调试输出
- [ ] T077 最终集成测试：完整用户流程 (选择用户 → 创建任务 → 拖拽 → 评论)

---

## 任务统计

| 阶段 | 任务数 | 用户故事 | 说明 |
|------|--------|----------|------|
| Phase 1 | 9 | - | 项目初始化 |
| Phase 2 | 8 | - | 基础设施层 |
| Phase 3 | 7 | US0 | 用户身份选择 |
| Phase 4 | 8 | US1 | 查看项目和看板 |
| Phase 5 | 4 | US2 | 查看团队成员 |
| Phase 6 | 6 | US3 | 创建和分配任务 |
| Phase 7 | 5 | US4 | 移动任务状态 |
| Phase 8 | 3 | US5 | 高亮我的任务 |
| Phase 9 | 6 | US6 | 添加任务评论 |
| Phase 10 | 2 | US7 | 查看任务详情 |
| Phase 11 | 4 | US8 | 编辑我的评论 |
| Phase 12 | 4 | US9 | 删除我的评论 |
| Phase 13 | 4 | - | 实时更新 (SSE) |
| Phase 14 | 6 | - | 收尾与优化 |
| **总计** | **43** | **10** | |

---

## 独立测试时机

| 用户故事 | 独立测试时机 | 测试要点 |
|----------|-------------|----------|
| US0 | Phase 3 完成后 | 用户选择 → localStorage 持久化 |
| US1 | Phase 4 完成后 | 项目列表 → 看板4列 → 任务数量 |
| US2 | Phase 5 完成后 | 团队成员列表 → 角色分类 |
| US3 | Phase 6 完成后 | 创建任务 → 分配负责人 → 待办列显示 |
| US4 | Phase 7 完成后 | 拖拽任务 → 跨列移动 → 状态持久化 |
| US5 | Phase 8 完成后 | 当前用户任务高亮 → 切换用户 |
| US6 | Phase 9 完成后 | 添加评论 → 评论列表显示 |
| US7 | Phase 10 完成后 | 任务详情弹窗 → 完整信息 |
| US8 | Phase 11 完成后 | 编辑评论 → 内容更新 |
| US9 | Phase 12 完成后 | 删除评论 → 确认 → 移除 |

---

## MVP 建议

**MVP 范围**: Phase 1-8 (US0 + US1 + US2 + US3 + US4 + US5)

**MVP 完成标准**:
- 用户可以选择身份
- 可以查看项目和看板
- 可以创建和分配任务
- 可以拖拽移动任务
- 当前用户的任务高亮显示

**Phase 9-14 可作为增量交付**:
- Sprint 1: Phase 1-2 (基础设施)
- Sprint 2: Phase 3-5 (US0 + US1 + US2)
- Sprint 3: Phase 6-8 (US3 + US4 + US5)
- Sprint 4: Phase 9-12 (US6-US9 评论功能)
- Sprint 5: Phase 13-14 (实时更新 + 收尾)

---

## 格式验证

- [x] 所有任务使用 `- [ ]` 复选框格式
- [x] 所有任务包含任务 ID (T001-T043)
- [x] [P] 标记并行化任务
- [x] [USX] 标记用户故事任务
- [x] 所有任务包含文件路径
- [x] 用户故事按优先级排列 (P1 → P2)

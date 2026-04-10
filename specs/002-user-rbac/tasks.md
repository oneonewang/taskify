# 任务清单：用户管理与基于角色的访问控制

**输入**: 设计文档 from `/specs/002-user-rbac/`
**前置条件**: plan.md, spec.md, data-model.md, contracts/api.md, research.md

**测试**: 本功能说明书未明确要求测试，采用手动验证方式。

**组织**: 任务按用户故事分组，每个故事可独立实现、测试和交付。

## 格式: `[ID] [P?] [Story] 描述`

- **[P]**: 可并行执行（不同文件，无依赖）
- **[Story]**: 所属用户故事（如 US1, US2, US3）
- 描述中包含精确文件路径

---

## Phase 1: 初始化设置

**目的**: 项目初始化和基础结构

- [x] T001 创建后端项目结构 (`backend/internal/{models,services,handlers,middleware,repository}`)
- [x] T002 创建前端项目结构 (`frontend/src/{views/admin,views/projects,stores,api,router}`)
- [x] T003 [P] 安装后端依赖: `github.com/gorilla/sessions`, `golang.org/x/crypto/bcrypt`, `github.com/gin-contrib/sessions`
- [x] T004 [P] 安装前端依赖: Pinia, Vue Router (已在 package.json)

---

## Phase 2: 基础层（阻塞前置条件）

**目的**: 所有用户故事实现前必须完成的核心基础设施

**⚠️ 关键**: 此阶段未完成前，任何用户故事工作都无法开始

- [x] T005 [P] 创建 User 模型 (`backend/internal/models/user.go`) - Email, PasswordHash, DisplayName, AvatarURL, EmailVerified, LastLoginAt
- [x] T006 [P] 创建 Role 模型 (`backend/internal/models/role.go`) - Name, Description, IsSystem, Scope
- [x] T007 [P] 创建 Permission 模型 (`backend/internal/models/permission.go`) - Resource, Action, Description
- [x] T008 [P] 创建 RolePermission 模型 (`backend/internal/models/role_permission.go`) - RoleID, PermissionID
- [x] T009 [P] 创建 ProjectMembership 模型 (`backend/internal/models/project_membership.go`) - UserID, ProjectID, RoleID, JoinedAt, 唯一约束
- [x] T010 [P] 创建 AuditLog 模型 (`backend/internal/models/audit_log.go`) - UserID, EventType, Details, IPAddress
- [x] T011 [P] 更新 Project 模型 (`backend/internal/models/project.go`) - 添加 IsArchived 字段
- [x] T012 [P] 更新 repository/sqlite.go - AutoMigrate 添加新模型
- [x] T013 创建数据库预定义数据填充 (`backend/internal/repository/seed.go`) - 预定义角色和权限
- [x] T014 创建 Session 配置 (`backend/config/session.go`) - gorilla/sessions 内存存储，Cookie 安全设置
- [x] T015 创建认证中间件 (`backend/internal/middleware/auth.go`) - 验证 session，获取 user_id
- [x] T016 创建 RBAC 中间件 (`backend/internal/middleware/rbac.go`) - 权限检查，项目成员验证

**检查点**: 基础层完成 - 用户故事实现可以开始了

---

## Phase 3: US1 - 用户注册与登录 (优先级: P1) 🎯 MVP

**目标**: 用户可以使用邮箱密码注册并登录系统

**独立测试**: 使用 curl 测试注册、登录、登出接口，验证 session cookie 设置

- [x] T017 [US1] 创建 AuthService (`backend/internal/services/auth_service.go`) - 注册、登录验证、登出
- [x] T018 [US1] 创建 AuthHandler (`backend/internal/handlers/auth.go`) - POST /api/auth/register, POST /api/auth/login, POST /api/auth/logout
- [x] T019 [US1] 创建 UserRepository (`backend/internal/repository/user_repo.go`) - 按邮箱查询、创建用户
- [x] T020 [US1] 添加密码哈希 - bcrypt cost=12
- [x] T021 [US1] 添加 AuditLog 记录 - 登录/登出事件
- [x] T022 [US1] 注册 auth 路由，添加中间件

**检查点**: US1 完成 - 用户可以注册、登录、登出

---

## Phase 4: US2 - 用户资料管理 (优先级: P2)

**目标**: 已登录用户可以查看和更新个人资料

**独立测试**: 登录后获取 /api/users/me，更新资料验证返回数据变化

- [x] T023 [US2] 创建 UserService (`backend/internal/services/user_service.go`) - 获取/更新当前用户、修改密码
- [x] T024 [US2] 创建 UserHandler (`backend/internal/handlers/user.go`) - GET /api/users/me, PUT /api/users/me, PUT /api/users/me/password
- [x] T025 [US2] 添加权限控制 - 只有本人可以查看/更新自己资料

**检查点**: US2 完成 - 用户可以管理自己的资料

---

## Phase 5: US3 - 角色管理 (优先级: P2)

**目标**: 管理员可以创建、编辑、删除自定义角色

**独立测试**: 管理员调用角色 CRUD 接口，验证系统预定义角色不可删除

- [x] T026 [US3] 创建 RoleService (`backend/internal/services/role_service.go`) - 角色 CRUD
- [x] T027 [US3] 创建 RoleHandler (`backend/internal/handlers/role.go`) - GET /api/roles, POST /api/roles, PUT /api/roles/:id, DELETE /api/roles/:id
- [x] T028 [US3] 创建 RoleRepository (`backend/internal/repository/role_repo.go`)
- [x] T029 [US3] 添加保护 - 禁止删除 IsSystem=true 的角色
- [x] T030 [US3] 添加保护 - 禁止删除已被分配的角色

**检查点**: US3 完成 - 管理员可以管理系统角色

---

## Phase 6: US4 - 权限配置 (优先级: P2)

**目标**: 管理员可以为角色分配和移除权限

**独立测试**: 编辑角色权限，验证权限变更立即生效

- [x] T031 [US4] 创建 PermissionService (`backend/internal/services/permission_service.go`)
- [x] T032 [US4] 创建 PermissionHandler (`backend/internal/handlers/permission.go`) - GET /api/permissions
- [x] T033 [US4] 扩展 RoleHandler - PUT /api/roles/:id/permissions
- [x] T034 [US4] 创建 PermissionRepository (`backend/internal/repository/permission_repo.go`)
- [x] T035 [US4] 添加 AuditLog 记录 - 权限变更

**检查点**: US4 完成 - 管理员可以配置角色权限

---

## Phase 7: US5 - 用户角色分配 (优先级: P1)

**目标**: 管理员可以为用户分配系统角色和项目成员角色

**独立测试**: 为用户分配角色，验证用户获得相应权限

- [x] T036 [US5] 创建 MembershipService (`backend/internal/services/membership_service.go`)
- [x] T037 [US5] 创建 MembershipHandler (`backend/internal/handlers/membership.go`) - /api/users/:id/roles, /api/projects/:id/members
- [x] T038 [US5] 创建 MembershipRepository (`backend/internal/repository/membership_repo.go`)
- [x] T039 [US5] 添加保护 - 禁止移除最后一个系统管理员
- [x] T040 [US5] 添加 /api/users/:id/project-memberships 端点

**检查点**: US5 完成 - 管理员可以为用户分配角色

---

## Phase 8: US6+US7 - 项目管理与项目成员邀请 (优先级: P1)

**目标**: 用户可以创建项目并管理项目成员

**独立测试**: 创建项目成为 Owner，添加/移除/修改成员角色

- [x] T041 [US6] 改造 ProjectService (`backend/internal/services/project_service.go`) - CRUD、自动分配 Owner
- [x] T042 [US6] 改造 ProjectHandler (`backend/internal/handlers/project.go`) - POST /api/projects, PUT /api/projects/:id, DELETE /api/projects/:id, POST /api/projects/:id/archive
  - **注意**: DELETE /api/projects/:id 必须 cascade 删除所有关联任务和评论
- [x] T043 [US6] 添加成员管理端点 - GET /api/projects/:id/members, POST /api/projects/:id/members, PUT /api/projects/:id/members/:user_id, DELETE /api/projects/:id/members/:user_id
- [x] T044 [US6] 添加权限检查 - 只有 Owner 可以管理成员、删除/归档项目
- [x] T045 [US6] 添加 AuditLog 记录 - 项目创建/删除/归档、成员变更

**检查点**: US6+US7 完成 - 项目创建和成员管理功能可用

---

## Phase 9: US8 - 前端登录与注册 (优先级: P1)

**目标**: 前端提供登录和注册页面

**独立测试**: 访问登录页，输入凭据，验证成功登录后跳转

- [x] T046 [US8] 创建登录页 (`frontend/src/views/Login.vue`)
- [x] T047 [US8] 创建注册页 (`frontend/src/views/Register.vue`)
- [x] T048 [US8] 创建 auth API 调用 (`frontend/src/api/auth.ts`)
- [x] T049 [US8] 创建 authStore (`frontend/src/stores/auth.ts`)
- [x] T050 [US8] 配置 Vue Router - `/login`, `/register`
- [x] T051 [US8] 添加路由守卫 - 未登录重定向到 /login
- [x] T052 [US8] 实现登出功能

**检查点**: US8 完成 - 前端用户可以注册和登录

---

## Phase 10: US9 - 管理员后台用户管理 (优先级: P1)

**目标**: 管理员可以在前端管理所有用户

**独立测试**: 访问 /admin/users，查看用户列表，修改用户角色

- [x] T053 [US9] 创建管理员布局 (`frontend/src/views/admin/Layout.vue`)
- [x] T054 [US9] 创建用户管理页面 (`frontend/src/views/admin/Users.vue`)
- [x] T055 [US9] 扩展 user API (`frontend/src/api/user.ts`)
- [x] T056 [US9] 添加路由 - `/admin`, `/admin/users`
- [x] T057 [US9] 添加路由守卫 - 非管理员访问 /admin 时重定向

**检查点**: US9 完成 - 管理员可以前端管理用户

---

## Phase 11: US10 - 管理员角色与权限管理 (优先级: P2)

**目标**: 管理员可以在前端配置角色和权限

**独立测试**: 访问 /admin/roles，编辑角色权限，验证变更

- [x] T058 [US10] 创建角色管理页面 (`frontend/src/views/admin/Roles.vue`)
- [x] T059 [US10] 创建权限表格组件 (`frontend/src/components/admin/PermissionTable.vue`)
- [x] T060 [US10] 扩展 role API (`frontend/src/api/role.ts`)
- [x] T061 [US10] 添加路由 - `/admin/roles`
- [x] T062 [US10] 创建审计日志查看 (`frontend/src/views/admin/AuditLogs.vue`)

**检查点**: US10 完成 - 管理员可以前端配置角色权限

---

## Phase 12: US11+US12 - 项目管理前端与前端权限控制 (优先级: P1)

**目标**: 前端实现项目管理页面和按钮级别权限控制

**独立测试**: 以不同角色登录，验证按钮显示/隐藏正确

- [x] T063 [US11] 创建项目管理页面 (`frontend/src/views/projects/ProjectList.vue`)
- [x] T064 [US11] 创建新建项目对话框 (`frontend/src/components/project/CreateProjectDialog.vue`)
- [x] T065 [US11] 创建项目设置页面 (`frontend/src/views/projects/ProjectSettings.vue`)
- [x] T066 [US11] 创建项目成员管理组件 (`frontend/src/components/project/ProjectMembers.vue`)
- [x] T067 [US11] 扩展 project API (`frontend/src/api/project.ts`)
- [x] T068 [US11] 添加路由 - `/projects`, `/projects/new`, `/projects/:id/settings`
- [x] T069 [US12] 创建权限函数 (`frontend/src/composables/usePermission.ts`)
- [x] T070 [US12] 改造看板页面 - Guest 隐藏创建/编辑任务按钮 (`frontend/src/views/ProjectKanban.vue`)
- [x] T071 [US12] 改造看板页面 - Member 隐藏删除项目按钮
- [x] T072 [US12] 全局权限控制 - 动态显示/隐藏操作按钮

**检查点**: US11+US12 完成 - 项目管理和前端权限控制完成

---

## Phase 13: 改造现有 API (FR-018)

**目的**: 为现有 API 添加认证和权限中间件

- [x] T073 改造 task.go - 所有端点添加 auth 中间件
- [x] T074 改造 task.go - 添加项目成员权限检查
- [x] T075 改造 project.go - 添加 auth 中间件
- [x] T076 改造 comment.go - 添加 auth 中间件
- [x] T077 [P] 更新前端 API 调用 - 所有请求携带 session cookie

**检查点**: 现有 API 全部添加权限控制

---

## Phase 14: 收尾与横切关注点

**目的**: 完善文档、安全加固、清理

- [x] T078 [P] 更新 README.md
- [x] T079 [P] 添加中文注释 - 检查所有新增文件的函数、类、模块注释
- [x] T080 代码清理 - 移除旧的 SeedUsers/SeedProjects 相关代码
- [x] T081 安全检查 - 确认密码 bcrypt，session HttpOnly Cookie
- [x] T082 验证 quickstart.md (N/A - 文件不存在)
- [x] T083 提交所有更改

---

## 依赖关系与执行顺序

### Phase 依赖

- Phase 1 → 无依赖
- Phase 2 → 依赖 Phase 1（阻塞所有用户故事）
- Phase 3-12 → 依赖 Phase 2
- Phase 13 → 依赖 Phase 3 (US1)
- Phase 14 → 依赖所有用户故事完成

### 用户故事依赖

| 故事 | 依赖 | 说明 |
|------|------|------|
| US1 (P1) | Phase 2 | 无需其他故事 |
| US2 (P2) | US1 | 需要认证基础 |
| US3 (P2) | US1 | 需要认证基础 |
| US4 (P2) | US3 | 需要角色基础 |
| US5 (P1) | US1 | 需要认证基础 |
| US6+US7 (P1) | US1 | 需要认证基础 |
| US8 (P1) | US1 | 需要后端 API 就绪 |
| US9 (P1) | US5 | 需要用户角色 API |
| US10 (P2) | US4 | 需要权限 API |
| US11+US12 (P1) | US6+US7 | 需要项目管理 API |

### 并行机会

- Phase 1: T001-T004 可并行
- Phase 2: T005-T012 (模型) 可并行
- Phase 3-12: 不同开发者可并行负责不同故事
- Phase 13: US1 完成后即可开始

---

## 任务统计

| Phase | 任务数 | 描述 |
|-------|--------|------|
| Phase 1 | 4 | 初始化设置 |
| Phase 2 | 12 | 基础层 |
| Phase 3 (US1) | 6 | 用户注册登录 |
| Phase 4 (US2) | 3 | 用户资料管理 |
| Phase 5 (US3) | 5 | 角色管理 |
| Phase 6 (US4) | 5 | 权限配置 |
| Phase 7 (US5) | 5 | 用户角色分配 |
| Phase 8 (US6+US7) | 5 | 项目管理与成员 |
| Phase 9 (US8) | 7 | 前端登录注册 |
| Phase 10 (US9) | 5 | 管理员用户管理 |
| Phase 11 (US10) | 5 | 管理员角色权限 |
| Phase 12 (US11+US12) | 10 | 项目管理前端+权限 |
| Phase 13 | 5 | 改造现有 API |
| Phase 14 | 6 | 收尾 |
| **总计** | **83** | |

---

## MVP 策略 (US1 仅)

1. 完成 Phase 1: 初始化
2. 完成 Phase 2: 基础层（关键）
3. 完成 Phase 3: US1 后端注册登录
4. **停止并验证**: 测试注册、登录、登出
5. 准备部署/demo

---

## 检查点清单

- [ ] Phase 1 完成
- [ ] Phase 2 完成 - 6 个模型、中间件就绪
- [ ] US1 完成 - 后端注册/登录/登出
- [ ] US2 完成 - 用户资料
- [ ] US3 完成 - 角色 CRUD
- [ ] US4 完成 - 权限配置
- [ ] US5 完成 - 用户角色分配
- [ ] US6+US7 完成 - 项目管理+成员
- [ ] US8 完成 - 前端登录注册
- [ ] US9 完成 - 管理员用户管理
- [ ] US10 完成 - 管理员角色权限
- [ ] US11+US12 完成 - 项目管理前端+按钮权限
- [ ] Phase 13 完成 - 现有 API 添加权限
- [ ] Phase 14 完成 - 文档完善

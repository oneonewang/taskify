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

- [x] T005 [P] 创建 User 模型 (`backend/internal/models/user.go`) - Email, PasswordHash, DisplayName, AvatarURL, EmailVerified, IsDisabled, LastLoginAt
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
  - **补充**: 修改密码后使其他会话失效（仅当前会话有效）

**检查点**: US2 完成 - 用户可以管理自己的资料

---

## Phase 4b: US2 补充 - 个人资料前端 (优先级: P2)

**目标**: 用户可以通过前端页面管理个人资料

**独立测试**: 登录后访问 /profile，修改资料验证变更

- [x] T026 [US2] 创建个人资料页面 (`frontend/src/views/Profile.vue`)
- [x] T027 [US2] 添加个人资料路由 - `/profile`
- [x] T028 [US2] 添加路由守卫 - 已登录用户才能访问

**检查点**: US2 前端完成 - 用户可以通过前端管理个人资料

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

## Phase 9b: US9b - 批量导入用户 (优先级: P2)

**目标**: 管理员可以通过 Excel 文件批量导入用户

**独立测试**: 上传 Excel 文件，验证导入结果（成功数量、失败报告）

- [x] T093 [US9b] 添加批量导入接口 - POST /api/admin/users/import，解析 .xlsx 文件
- [x] T094 [US9b] 添加导入模板下载接口 - GET /api/admin/users/import/template
- [x] T095 [US9b] 实现 Excel 解析 - 读取邮箱、显示名称、角色字段
- [x] T096 [US9b] 实现批量创建逻辑 - 统一默认密码（admin123），跳过已存在邮箱
- [x] T097 [US9b] 添加 AuditLog 记录 - 批量导入事件
- [x] T098 [US9b] 创建批量导入组件 (`frontend/src/components/admin/ImportUsersDialog.vue`)
- [x] T099 [US9b] 在用户管理页面添加"批量导入"按钮 (`frontend/src/views/admin/Users.vue`)
- [x] T100 [US9b] 添加导入结果展示组件

**检查点**: US9b 完成 - 管理员可以批量导入用户

---

## Phase 10b: US9c - 用户查询、重置密码与禁用 (优先级: P2)

**目标**: 管理员可以查询用户、重置密码和禁用账号

**独立测试**: 在用户管理页面完成搜索、重置密码、禁用/启用操作

- [x] T098 [US9c] 添加用户查询接口 - GET /api/admin/users，支持邮箱、显示名称、角色、状态筛选，分页
- [x] T099 [US9c] 添加重置密码接口 - POST /api/admin/users/:id/reset-password
- [x] T100 [US9c] 添加禁用用户接口 - POST /api/admin/users/:id/disable
- [x] T101 [US9c] 添加启用用户接口 - POST /api/admin/users/:id/enable
- [x] T102 [US9c] [P] 更新 UserRepository - 添加按条件查询方法（含 is_disabled 筛选）
- [x] T103 [US9c] 改造用户管理页面 - 添加搜索筛选和分页功能 (`frontend/src/views/admin/Users.vue`)
- [x] T104 [US9c] 改造用户管理页面 - 添加重置密码、禁用、启用按钮
- [x] T105 [US9c] 添加 AuditLog 记录 - 密码重置、账号禁用/启用

**检查点**: US9c 完成 - 管理员可以查询用户、重置密码和禁用账号

---

## Phase 11: US10 - 管理员角色与权限管理 (优先级: P2)

**目标**: 管理员可以在前端配置角色和权限

**独立测试**: 访问 /admin/roles，编辑角色权限，验证变更

- [x] T058 [US10] 创建角色管理页面 (`frontend/src/views/admin/Roles.vue`)
- [x] T059 [US10] 创建权限表格组件 (`frontend/src/components/admin/PermissionTable.vue`)
- [x] T060 [US10] 扩展 role API (`frontend/src/api/role.ts`)
- [x] T061 [US10] 添加路由 - `/admin/roles`
- [x] T062 [US10] 创建审计日志查看 (`frontend/src/views/admin/AuditLogs.vue`)
- [x] T062b [US10] 添加审计日志路由 - `/admin/audit-logs`

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

## Phase 15: US13 - 项目列表增强 (优先级: P2)

**目标**: 项目列表支持分页、过滤，任务卡片显示所有者信息

**前置**: US11 (Phase 12)

**独立测试**: 访问项目列表，验证分页切换、名称搜索、状态筛选、任务卡显示所有者

- [x] T106 [US13] 扩展项目列表 API - 添加分页参数 `page`, `page_size` 和过滤参数 `keyword`（模糊搜索）, `is_archived`, `owner_name`（所有者名称模糊搜索），多条件 AND 组合 (`backend/internal/handlers/project.go`)
- [x] T107 [US13] 更新 ProjectRepository - 添加分页查询方法 (`backend/internal/repository/project_repo.go`)
- [x] T108 [US13] 更新项目列表前端 - 添加分页器组件 (`frontend/src/views/projects/ProjectList.vue`)
- [x] T109 [US13] 更新项目列表前端 - 添加搜索框、状态筛选和所有者名称筛选 (`frontend/src/views/projects/ProjectList.vue`)
- [x] T110 [US13] 更新任务卡片组件 - 显示项目所有者头像和名称 (`frontend/src/views/projects/ProjectList.vue` 项目卡片)
- [ ] T111 [US13] 添加 AuditLog 记录 - 分页查询操作（N/A，可选）

**检查点**: US13 完成 - 项目列表支持分页和过滤，任务卡显示所有者

---

## Phase 16: US14 - 任务分享功能 (优先级: P2)

**目标**: 用户可以分享任务链接，其他成员通过链接直接访问任务详情

**前置**: Phase 8 (US6+US7 项目管理与成员)

**独立测试**: 在任务详情页点击分享，复制链接到新标签页打开，验证任务详情显示和非成员访问返回403

- [x] T112 [US14] 添加任务详情 API - GET /api/projects/:projectId/tasks/:taskId，验证项目成员权限 (`backend/internal/handlers/task.go`)
- [x] T113 [US14] 添加 TaskService.GetTaskByProject 方法 - 联合查询验证项目和任务存在性及成员资格 (`backend/internal/services/task_service.go`)
- [x] T114 [US14] 前端路由添加任务详情路由 - `/projects/:projectId/tasks/:taskId` (`frontend/src/router/index.ts`)
- [x] T115 [US14] TaskDetail 支持 URL 直接访问 - 接收 route params 加载任务 (`frontend/src/components/task/TaskDetail.vue`)
- [x] T116 [US14] 任务详情页添加分享按钮 - 点击复制链接到剪贴板 (`frontend/src/components/task/TaskDetail.vue`)
- [x] T117 [US14] 分享链接非成员访问返回 403 提示 (`frontend/src/components/task/TaskDetail.vue`)

**检查点**: US14 完成 - 任务可以分享，成员可通过链接访问

- [x] Phase 16 完成 - 任务分享功能（US14）

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
- Phase 15 → 依赖 Phase 12 (US11+US12)

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
| US9b (P2) | US1 | 需要认证基础 |
| US9c (P2) | US9 | 需要用户管理基础 |
| US10 (P2) | US4 | 需要权限 API |
| US11+US12 (P1) | US6+US7 | 需要项目管理 API |
| US13 (P2) | US11+US12 | 需要项目管理前端基础 |
| US14 (P2) | US6+US7 | 需要任务基础和项目管理 |

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
| Phase 4 (US2) | 6 | 用户资料管理（含前端） |
| Phase 5 (US3) | 5 | 角色管理 |
| Phase 6 (US4) | 5 | 权限配置 |
| Phase 7 (US5) | 5 | 用户角色分配 |
| Phase 8 (US6+US7) | 5 | 项目管理与成员 |
| Phase 9 (US8) | 7 | 前端登录注册 |
| Phase 9b (US9b) | 8 | 批量导入用户（含前端） |
| Phase 10 (US9) | 5 | 管理员用户管理 |
| Phase 10b (US9c) | 8 | 用户查询、重置密码、禁用 |
| Phase 11 (US10) | 6 | 管理员角色权限（含审计日志路由） |
| Phase 12 (US11+US12) | 10 | 项目管理前端+权限 |
| Phase 13 | 5 | 改造现有 API |
| Phase 14 | 6 | 收尾 |
| Phase 15 (US13) | 6 | 项目列表增强（分页/过滤/所有者显示） |
| **总计** | **109** | |

---

## MVP 策略 (US1 仅)

1. 完成 Phase 1: 初始化
2. 完成 Phase 2: 基础层（关键）
3. 完成 Phase 3: US1 后端注册登录
4. **停止并验证**: 测试注册、登录、登出
5. 准备部署/demo

---

## 检查点清单

- [x] Phase 1 完成
- [x] Phase 2 完成 - 6 个模型、中间件就绪
- [x] US1 完成 - 后端注册/登录/登出
- [x] US2 完成 - 用户资料（含前端）
- [x] US3 完成 - 角色 CRUD
- [x] US4 完成 - 权限配置
- [x] US5 完成 - 用户角色分配
- [x] US6+US7 完成 - 项目管理+成员
- [x] US8 完成 - 前端登录注册
- [x] US9 完成 - 管理员用户管理
- [x] US9b 完成 - 批量导入用户
- [x] US9c 完成 - 用户查询、重置密码、禁用
- [x] US10 完成 - 管理员角色权限
- [x] US11+US12 完成 - 项目管理前端+按钮权限
- [x] Phase 13 完成 - 现有 API 添加权限
- [x] Phase 14 完成 - 文档完善
- [x] Phase 15 完成 - 项目列表增强（分页/过滤/所有者显示）

# Implementation Plan: 用户管理与基于角色的访问控制

**Branch**: `002-user-rbac` | **Date**: 2026-04-10 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-user-rbac/spec.md`

## Summary

实现完整的用户认证系统和基于角色的访问控制（RBAC）。后端采用 Go + Gin + GORM + SQLite，前端采用 Vue 3 + Element Plus + Pinia + Vue Router。会话管理使用服务器端 Session + HTTP-only Cookie。权限模型基于项目成员资格（ProjectMembership），支持多项目角色。

## Technical Context

**Language/Version**: Go 1.21+
**Primary Dependencies**: Gin (Web框架), GORM (ORM), gorilla/sessions (会话管理), bcrypt (密码哈希), excelize (Excel 解析)
**Storage**: SQLite (现有数据库废弃重建，使用 GORM AutoMigrate)
**Testing**: go test (Go 标准测试框架)
**Target Platform**: Linux Server (Web 应用)
**Project Type**: Web Service + SPA (前后端分离)
**Performance Goals**: 登录响应 <2秒，权限拒绝 <500ms (p95)，100% 安全事件记录
**Constraints**: <200ms p95 API 响应，<100MB 内存
**Scale/Scope**: 小团队(5人)，MVP阶段，单体架构

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### I. 安全优先 (Security First) - ✅ 通过
- 所有用户输入通过 Gin binding tags 验证（type, range, format, length）
- 使用 GORM 的 Parameterized Queries 防止 SQL 注入
- 密码使用 bcrypt 哈希存储（cost=12，不使用明文）
- 敏感操作（登录、登出、权限更改、密码重置、账号禁用/启用）记录 AuditLog
- 会话令牌存储在 HTTP-only Cookie 中，防止 XSS 攻击

### II. 输入验证 (Input Validation) - ✅ 通过
- 请求参数使用 Gin 的 ShouldBindJSON 进行类型检查和验证
- 邮箱格式使用正则白名单验证
- 密码长度限制 8-72 字符
- GORM 模型层定义 validate tags

### III. 微服务架构 (Microservices Architecture) - ✅ 通过（v1.0 特例）
- v1.0 MVP 采用单体架构
- 服务边界清晰划定：用户模块、角色权限模块、项目模块、任务模块
- 后续扩展时可拆分微服务

### IV. 完全文档化 (Code Documentation) - ✅ 通过
- 所有公共 API、函数、类包含中文注释
- 数据库 Schema 字段包含说明
- 复杂业务逻辑包含行内注释

### V. 中文内容 (Chinese Content) - ✅ 通过
- 代码注释使用简体中文
- README 和文档使用简体中文
- 日志消息使用简体中文

### VI. 分层设计 (Layered Architecture) - ✅ 通过
- Handler (表现层) → Service (业务逻辑层) → Repository (数据访问层)
- Handler 调用 Service，Service 调用 Repository
- Repository 只能访问数据库
- 无跨层调用

## Project Structure

### Documentation (this feature)

```text
specs/002-user-rbac/
├── plan.md              # 本文件
├── research.md          # Phase 0 研究输出
├── data-model.md        # Phase 1 数据模型输出
├── quickstart.md        # Phase 1 快速开始指南
├── contracts/           # Phase 1 接口契约
│   └── api-contracts.md
└── tasks.md             # Phase 2 任务清单
```

### Source Code (repository root)

```text
backend/
├── cmd/server/main.go          # 服务入口
├── internal/
│   ├── models/                 # 数据模型（User, Role, Permission, ProjectMembership, AuditLog）
│   ├── repository/             # 数据访问层（DB 操作）
│   ├── services/               # 业务逻辑层（认证服务、权限服务）
│   ├── handlers/               # 表现层（API handlers）
│   └── middleware/             # 中间件（Auth, Permission）
├── pkg/
│   └── response/               # 统一响应格式
└── config/                     # 配置

frontend/
├── src/
│   ├── api/                    # API 调用（user.ts, auth.ts, project.ts）
│   ├── components/             # 公共组件
│   ├── views/                  # 页面（Login, Register, Profile, Admin/, Project/）
│   ├── router/                 # 路由配置（含路由守卫）
│   ├── stores/                 # Pinia 状态管理（authStore, userStore）
│   └── composables/            # 组合式函数
```

**Structure Decision**: Web application (前后端分离)。后端采用分层架构（Handler → Service → Repository），前端采用 Vue 3 SPA 结构。现有 backend/ 和 frontend/ 目录结构保持不变，新增模块按分层原则添加。

## Phase 0: Research

### 研究主题

#### 1. Go 会话管理方案
**问题**: 服务器端会话管理在 Go/Gin 中如何实现？
**决策**: 使用 gorilla/sessions 库 + 内存存储
**理由**: gorilla/sessions 是 Go 最成熟的会话管理库，支持 HTTP-only Cookie。v1.0 MVP 使用内存存储（map），后续可升级为 Redis 分布式存储。

#### 2. 密码哈希算法
**问题**: Go 中密码存储使用 bcrypt 还是 argon2？
**决策**: bcrypt
**理由**: bcrypt 是行业标准，Golang crypto 包原生支持。argon2 虽然更安全但配置复杂，对于内部工具 bcrypt 足够。

#### 3. RBAC 权限检查模式
**问题**: 如何高效实现基于项目成员的权限检查？
**决策**: 中间件层统一拦截 + Service 层权限校验
**理由**: Gin 中间件在路由层统一拦截请求，从 Session 获取用户 ID，从 URL 参数或 Body 获取项目 ID，查询 ProjectMembership 表验证权限。

#### 4. Vue 3 路由守卫模式
**问题**: 前端如何实现认证守卫？
**决策**: Vue Router 4 的 navigation guards + Pinia auth store
**理由**: Vue Router 4 提供 beforeEach 守卫，在全局层面检查 auth store 中的登录状态，未登录重定向到登录页。

#### 5. Cookie 安全配置
**问题**: 会话 Cookie 如何配置才安全？
**决策**: HttpOnly=true, Secure=true(生产), SameSite=Strict
**理由**: HttpOnly 防止 XSS 读取，Secure 在 HTTPS 下才传输，SameSite 防止 CSRF。开发环境 Secure=false。

#### 6. 用户禁用实现方案
**问题**: 如何实现用户账号的软禁用？
**决策**: 在 User 模型中添加 is_disabled 布尔字段，登录时检查该字段
**理由**: 软禁用保留数据完整性，可随时恢复，比硬删除更安全。符合 FR-036 要求。

#### 7. 密码重置机制
**问题**: 管理员重置密码后如何传递新密码？
**决策**: 重置为统一临时密码（admin123），用户首次登录强制修改
**理由**: 符合 FR-032 批量导入的现有模式，统一管理更方便。当前无邮件基础设施。

#### 8. Excel 文件解析
**问题**: Go 如何解析 .xlsx 格式的批量导入文件？
**决策**: 使用 excelize 库
**理由**: excelize 是 Go 语言处理 Office Excel 文件的开源库，支持 .xlsx 格式的读取和写入。

### 技术选型总结

| 组件 | 技术选型 | 备选方案 |
|------|---------|---------|
| 后端框架 | Gin | Echo, Fiber |
| ORM | GORM | sqlx, raw SQL |
| 会话管理 | gorilla/sessions | gin-contrib/sessions |
| 密码哈希 | bcrypt | argon2 |
| Excel 解析 | excelize | xlsx, 360-info/go-zoom |
| 前端框架 | Vue 3 | React |
| UI 组件库 | Element Plus | Ant Design Vue |
| 状态管理 | Pinia | Vuex |
| 路由 | Vue Router 4 | - |

## Phase 1: Design & Contracts

### 数据模型

#### User 模型更新
```go
type User struct {
    ID            uint      `gorm:"primaryKey"`
    Email         string    `gorm:"uniqueIndex;not null"`
    PasswordHash  string    `gorm:"not null"`
    DisplayName   string
    AvatarURL     string
    EmailVerified bool      `gorm:"default:false"`
    IsDisabled    bool      `gorm:"default:false"`  // 新增：账号禁用标志
    LastLoginAt   *time.Time
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

#### 新增 API 端点

| 端点 | 方法 | 描述 |
|------|------|------|
| `/api/admin/users` | GET | 查询用户列表，支持筛选和分页 |
| `/api/admin/users/:id/reset-password` | POST | 重置用户密码为临时密码 |
| `/api/admin/users/:id/disable` | POST | 禁用用户账号 |
| `/api/admin/users/:id/enable` | POST | 启用已禁用的用户账号 |

#### 查询参数设计
```
GET /api/admin/users?email=xxx&display_name=xxx&role=admin&is_disabled=false&page=1&page_size=20
```

### 快速开始

详见 [quickstart.md](./quickstart.md)

## 复杂度追踪

无违规项。

## 实施检查点

- [x] Phase 1 完成 - 项目初始化
- [x] Phase 2 完成 - 基础层（6个模型、中间件）
- [x] Phase 3 完成 - US1 用户注册登录
- [x] Phase 4 完成 - US2 用户资料管理
- [x] Phase 4b 完成 - US2 个人资料前端
- [x] Phase 5 完成 - US3 角色管理
- [x] Phase 6 完成 - US4 权限配置
- [x] Phase 7 完成 - US5 用户角色分配
- [x] Phase 8 完成 - US6+US7 项目管理与成员
- [x] Phase 9 完成 - US8 前端登录注册
- [x] Phase 9b 完成 - US9b 批量导入用户
- [x] Phase 10 完成 - US9 管理员用户管理
- [x] Phase 10b 完成 - US9c 用户查询、重置密码、禁用
- [x] Phase 11 完成 - US10 管理员角色权限
- [x] Phase 12 完成 - US11+US12 项目管理前端+权限
- [x] Phase 13 完成 - 改造现有 API
- [x] Phase 14 完成 - 收尾

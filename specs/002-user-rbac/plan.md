# Implementation Plan: 用户管理与基于角色的访问控制

**Branch**: `002-user-rbac` | **Date**: 2026-04-10 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-user-rbac/spec.md`

## Summary

实现完整的用户认证系统和基于角色的访问控制（RBAC）。后端采用 Go + Gin + GORM + SQLite，前端采用 Vue 3 + Element Plus + Pinia + Vue Router。会话管理使用服务器端 Session + HTTP-only Cookie。权限模型基于项目成员资格（ProjectMembership），支持多项目角色。

## Technical Context

**Language/Version**: Go 1.21+
**Primary Dependencies**: Gin (Web框架), GORM (ORM), gorilla/sessions (会话管理), bcrypt (密码哈希)
**Storage**: SQLite (现有数据库)
**Testing**: go test (Go 标准测试框架)
**Target Platform**: Linux Server (Web 应用)
**Project Type**: Web Service + SPA (前后端分离)
**Performance Goals**: 登录响应 <5秒，权限拒绝 <1秒，100% 安全事件记录
**Constraints**: <200ms p95 API 响应，<100MB 内存
**Scale/Scope**: 小团队(5人)，MVP阶段，单体架构

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### I. 安全优先 (Security First) - ✅ 通过
- 所有用户输入通过 Gin binding tags 验证（type, range, format, length）
- 使用 GORM 的 Parameterized Queries 防止 SQL 注入
- 密码使用 bcrypt 哈希存储（不使用明文）
- 敏感操作（登录、登出、权限变更）记录 AuditLog
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
└── tasks.md             # Phase 2 任务清单（/speckit.tasks 命令输出，非本命令生成）
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

### 技术选型总结

| 组件 | 技术选型 | 备选方案 |
|------|---------|---------|
| 后端框架 | Gin | Echo, Fiber |
| ORM | GORM | sqlx, raw SQL |
| 会话管理 | gorilla/sessions | gin-contrib/sessions |
| 密码哈希 | bcrypt | argon2 |
| 前端框架 | Vue 3 | React |
| UI 组件库 | Element Plus | Ant Design Vue |
| 状态管理 | Pinia | Vuex |
| 路由 | Vue Router 4 | - |

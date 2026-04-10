# Implementation Plan: 用户管理与 RBAC

**Branch**: `002-user-rbac` | **Date**: 2026-04-10 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-user-rbac/spec.md`

## Summary

为 Taskify 实现完整的用户管理与基于角色的访问控制（RBAC）系统。包括用户注册/登录、会话管理、角色权限定义、项目级成员资格控制、以及对现有 API 的认证改造。

**技术方案**: 基于 Go 1.21 + Gin + GORM + SQLite 的后端服务。使用服务器端会话（HTTP-only Cookie）， bcrypt 密码哈希，遵循分层架构设计（Handler → Service → Repository）。

## Technical Context

**Language/Version**: Go 1.21
**Primary Dependencies**: Gin, GORM, SQLite, gorilla/sessions (会话管理), bcrypt (密码哈希)
**Storage**: SQLite（现有数据库废弃重建，使用 GORM AutoMigrate）
**Testing**: Go testing / httptest
**Target Platform**: Linux Server
**Project Type**: Web Service (RESTful API)
**Performance Goals**: 登录 < 5s, 权限检查 < 1s
**Constraints**: 需兼容现有 API 端点改造
**Scale/Scope**: 小团队（5人）内部工具，支持多项目成员资格

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则 | 状态 | 说明 |
|------|------|------|
| VI. 分层设计 | PASS | Handler → Service → Repository 层级清晰 |
| I. 安全优先 | PASS | 密码 bcrypt哈希、会话Cookie HttpOnly、审计日志 |
| II. 输入验证 | PASS | GORM + 自定义 Validator |
| III. 微服务架构 | JUSTIFIED | MVP 单体架构符合 v1.0 章程，服务边界已划定 |
| IV. 完全文档化 | PASS | 中文注释，公共 API 需 JSDoc 风格 |
| V. 中文内容 | PASS | 所有文档、注释使用简体中文 |

**无架构违规，无需 Complexity Tracking**

## Project Structure

### Documentation (this feature)

```text
specs/002-user-rbac/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
backend/
├── src/
│   ├── models/           # GORM 模型定义
│   │   ├── user.go
│   │   ├── role.go
│   │   ├── permission.go
│   │   ├── project_membership.go
│   │   └── audit_log.go
│   ├── services/         # 业务逻辑层
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── role_service.go
│   │   ├── permission_service.go
│   │   └── project_service.go
│   ├── handlers/        # 表现层（已有，逐步改造）
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── role.go
│   │   └── project.go
│   ├── middleware/      # 中间件
│   │   ├── auth.go      # 认证中间件
│   │   └── rbac.go      # 权限检查中间件
│   ├── repository/      # 数据访问层
│   │   ├── user_repo.go
│   │   ├── role_repo.go
│   │   └── project_repo.go
│   └── database/        # 数据库初始化
│       ├── db.go
│       └── migrations/
└── tests/
    ├── unit/
    └── integration/

frontend/
└── [暂不在本次范围]
```

**Structure Decision**: 采用单体架构，符合 v1.0 MVP 章程。分层清晰：Handler → Service → Repository。

## Phase 0: Research (DONE)

### 调研结果摘要

| # | 问题 | 决定 | 位置 |
|---|------|------|------|
| Q1 | 会话存储 | 内存存储 (gorilla/sessions) | research.md |
| Q2 | bcrypt cost | 12 | research.md |
| Q3 | 角色权限设计 | 系统级(管理员) + 项目级(Owner/Member/Guest) | research.md |
| Q4 | 审计日志 | 硬删除 + 定期归档 | research.md |

**✅ Phase 0 完成 - research.md 已生成**

## Phase 1: Design & Contracts (DONE)

### 生成的 Phase 1 工件

| 文件 | 状态 | 说明 |
|------|------|------|
| data-model.md | ✅ | 7个实体定义，完整 ER 图 |
| contracts/api.md | ✅ | 完整 REST API 契约 |
| quickstart.md | ✅ | 开发者快速开始指南 |

### Constitution Check (Phase 1 后复检)

| 原则 | 状态 | 说明 |
|------|------|------|
| VI. 分层设计 | ✅ PASS | Phase 1 验证：Handler → Service → Repository 分层清晰，无跨层调用 |
| I. 安全优先 | ✅ PASS | bcrypt cost=12, HttpOnly Cookie, 审计日志 |
| II. 输入验证 | ✅ PASS | GORM validation tags + 自定义 Validator |
| III. 微服务架构 | ✅ JUSTIFIED | MVP 单体架构，数据模型支持未来扩展 |
| IV. 完全文档化 | ✅ PASS | data-model.md, contracts/api.md, quickstart.md |
| V. 中文内容 | ✅ PASS | 所有文档使用简体中文 |

**✅ 所有 Constitution Gates 通过 - Phase 1 完成**

## Phase 2: Implementation Tasks (待执行)

> Phase 2 由 `/speckit-tasks` 命令触发，生成 tasks.md

**预估任务清单**：

1. 创建数据库模型层 (models/)
2. 创建 Repository 数据访问层
3. 创建 Service 业务逻辑层
4. 创建 Auth 认证 Handler
5. 创建 User/Role/Permission Handlers
6. 创建 Auth/RBAC 中间件
7. 改造现有 Handlers (project/task/comment)
8. 数据库迁移与预置数据填充
9. 单元测试编写
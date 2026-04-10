# Phase 1 Quickstart: 用户管理与 RBAC

## 快速开始

### 1. 环境要求

- Go 1.21+
- SQLite3

### 2. 依赖安装

```bash
cd backend
go get github.com/gorilla/sessions
go get golang.org/x/crypto/bcrypt
go get github.com/gin-contrib/sessions
```

### 3. 数据库迁移

```bash
# 首次启动时，GORM AutoMigrate 会自动创建所有表
# 迁移顺序：User → Role → Permission → RolePermission → ProjectMembership → AuditLog
```

### 4. 启动服务

```bash
cd backend/cmd/server
go run main.go
```

服务将在 `http://localhost:8080` 启动。

### 5. 初始化数据

首次启动时，系统会自动创建预定义角色和权限：

```sql
-- 预定义角色
INSERT INTO roles (name, description, is_system, scope) VALUES
('admin', '系统管理员', true, 'system'),
('owner', '项目所有者', true, 'project'),
('member', '项目成员', true, 'project'),
('guest', '项目访客', true, 'project');

-- 预定义权限
INSERT INTO permissions (resource, action, description) VALUES
('users', 'view', '查看用户'),
('users', 'manage', '管理用户'),
('projects', 'view', '查看项目'),
('projects', 'create', '创建项目'),
('projects', 'edit', '编辑项目'),
('projects', 'delete', '删除项目'),
('projects', 'manage', '管理项目'),
('tasks', 'view', '查看任务'),
('tasks', 'create', '创建任务'),
('tasks', 'edit', '编辑任务'),
('tasks', 'move', '移动任务'),
('tasks', 'delete', '删除任务'),
('comments', 'view', '查看评论'),
('comments', 'create', '创建评论'),
('comments', 'edit', '编辑评论'),
('comments', 'delete', '删除评论'),
('roles', 'view', '查看角色'),
('roles', 'create', '创建角色'),
('roles', 'edit', '编辑角色'),
('roles', 'delete', '删除角色'),
('roles', 'assign', '分配角色');
```

### 6. API 测试

**注册用户**：
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Test1234","display_name":"测试用户"}'
```

**登录**：
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{"email":"test@example.com","password":"Test1234"}'
```

**获取当前用户**：
```bash
curl http://localhost:8080/api/users/me -b cookies.txt
```

## 分层架构说明

```
┌─────────────────────────────────────────────────────────────┐
│  Handler (handlers/)                                         │
│  - 解析请求参数                                               │
│  - 调用 Service 层                                            │
│  - 返回统一响应格式                                            │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│  Service (services/)                                         │
│  - 业务逻辑处理                                               │
│  - 事务管理                                                   │
│  - 权限检查                                                   │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│  Repository (repository/)                                    │
│  - 数据访问层                                                 │
│  - GORM 操作                                                  │
│  - 禁止跨层调用                                               │
└─────────────────────────────────────────────────────────────┘
```

## 新增文件清单

### Models
- `backend/src/models/user.go` - 用户模型
- `backend/src/models/role.go` - 角色模型
- `backend/src/models/permission.go` - 权限模型
- `backend/src/models/project_membership.go` - 项目成员模型
- `backend/src/models/audit_log.go` - 审计日志模型

### Services
- `backend/src/services/auth_service.go` - 认证服务
- `backend/src/services/user_service.go` - 用户服务
- `backend/src/services/role_service.go` - 角色服务
- `backend/src/services/permission_service.go` - 权限服务
- `backend/src/services/project_service.go` - 项目服务（改造）

### Handlers
- `backend/src/handlers/auth.go` - 认证处理器
- `backend/src/handlers/user.go` - 用户处理器
- `backend/src/handlers/role.go` - 角色处理器
- `backend/src/handlers/project.go` - 项目处理器（改造）

### Middleware
- `backend/src/middleware/auth.go` - 认证中间件
- `backend/src/middleware/rbac.go` - RBAC 中间件

### Repository
- `backend/src/repository/user_repo.go` - 用户仓库
- `backend/src/repository/role_repo.go` - 角色仓库
- `backend/src/repository/project_repo.go` - 项目仓库

### Database
- `backend/src/database/migrations.go` - 数据库迁移逻辑
- `backend/src/database/seed.go` - 预定义数据填充

## 配置项

| 环境变量 | 默认值 | 说明 |
|---------|-------|------|
| SESSION_SECRET | (必需) | 会话加密密钥 |
| SESSION_MAX_AGE | 86400 | 会话有效期（秒） |
| BCRYPT_COST | 12 | bcrypt 哈希强度 |
| COOKIE_SECURE | false | HTTPS only cookie（生产环境设为 true） |
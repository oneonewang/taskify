# Phase 1 Data Model: 用户管理与 RBAC

## 实体关系图

```
┌─────────────┐     ┌─────────────────────┐     ┌─────────────┐
│    User     │────<│ ProjectMembership   >────│   Project   │
└─────────────┘     └─────────────────────┘     └─────────────┘
       │                    │
       │                    │
       │                    ▼
       │            ┌─────────────┐
       └───────────>│    Role     │
                    └─────────────┘
                           │
                           │
                           ▼
                    ┌─────────────┐
                    │   Permission│
                    └─────────────┘

       ┌─────────────┐
       │  AuditLog   │
       └─────────────┘
```

## 数据模型

### 1. User（用户）

```go
type User struct {
    ID             uint      `gorm:"primaryKey"`           // 唯一标识符
    Email          string    `gorm:"uniqueIndex;size:255"` // 邮箱地址
    DisplayName    string    `gorm:"size:100"`             // 显示名称
    AvatarURL      string    `gorm:"size:500"`             // 头像 URL
    PasswordHash   string    `gorm:"size:255"`             // bcrypt 哈希密码
    EmailVerified  bool      `gorm:"default:false"`        // 邮箱已验证
    LastLoginAt    *time.Time                           // 最近登录时间戳
    CreatedAt      time.Time                             // 创建时间戳
    UpdatedAt      time.Time                             // 更新时间戳
}
```

**验证规则**：
- Email: 有效邮箱格式，唯一索引
- DisplayName: 1-100 字符
- PasswordHash: bcrypt cost=12，最少 8 字符

---

### 2. Role（角色）

```go
type Role struct {
    ID          uint      `gorm:"primaryKey"`            // 唯一标识符
    Name        string    `gorm:"uniqueIndex;size:50"`  // 角色名称
    Description string    `gorm:"size:255"`             // 角色描述
    IsSystem    bool      `gorm:"default:false"`        // 是否系统预定义
    Scope       string    `gorm:"size:20;default:'system'"` // 作用域: system/project
    CreatedAt   time.Time                             // 创建时间戳
    UpdatedAt   time.Time                             // 更新时间戳
}
```

**预定义角色**：

| Name | Scope | IsSystem | Description |
|------|-------|----------|-------------|
| admin | system | true | 系统管理员 |
| owner | project | true | 项目所有者 |
| member | project | true | 项目成员 |
| guest | project | true | 项目访客 |

---

### 3. Permission（权限）

```go
type Permission struct {
    ID           uint      `gorm:"primaryKey"`           // 唯一标识符
    Resource     string    `gorm:"index;size:50"`      // 资源类型
    Action       string    `gorm:"index;size:50"`      // 操作名称
    Description  string    `gorm:"size:255"`           // 权限描述
    CreatedAt    time.Time                             // 创建时间戳
}
```

**权限格式**: `{resource}.{action}` (如 `tasks.create`)

**预定义权限**：

| Resource | Actions | Description |
|----------|---------|-------------|
| users | view, manage | 用户管理 |
| projects | view, create, edit, delete, manage | 项目管理 |
| tasks | view, create, edit, move, delete | 任务管理 |
| comments | view, create, edit, delete | 评论管理 |
| roles | view, create, edit, delete, assign | 角色管理 |

---

### 4. RolePermission（角色权限关联）

```go
type RolePermission struct {
    RoleID       uint `gorm:"primaryKey"`
    PermissionID uint `gorm:"primaryKey"`
}
```

---

### 5. ProjectMembership（项目成员资格）

```go
type ProjectMembership struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"index"`
    ProjectID uint      `gorm:"index"`
    RoleID    uint      `gorm:"index"`
    JoinedAt  time.Time // 加入时间
}
```

**唯一约束**: (UserID, ProjectID) - 每个用户在每个项目只能有一个角色

---

### 6. AuditLog（审计日志）

```go
type AuditLog struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"index"`
    EventType string    `gorm:"index;size:50"` // login, logout, permission_change, access_denied
    Details   string    `gorm:"type:text"`
    IPAddress string    `gorm:"size:45"` // IPv6 兼容
    CreatedAt time.Time
}
```

**记录事件类型**：
- `login` - 用户登录
- `logout` - 用户登出
- `permission_change` - 权限变更
- `access_denied` - 访问被拒绝
- `project_created` - 项目创建
- `project_deleted` - 项目删除

---

### 7. Session（会话）

```go
// 使用 gorilla/sessions 存储在内存中
// Cookie 名: session_id
// HttpOnly: true
// Secure: 根据环境（生产环境为 true）
// MaxAge: 86400 (24小时)
```

---

## 状态转换

### 用户状态
```
[注册] → [已创建(未验证)] → [已验证] → [活跃]
                              ↓
                        [已禁用(管理员)]
```

### 项目状态
```
[创建] → [活动中] → [已归档] → [已删除]
```

### 角色分配
```
[添加成员] → [活跃成员] → [移除成员]
                              ↓
                         [角色变更]
```

---

## 索引设计

| 表 | 索引类型 | 字段 |
|----|---------|------|
| users | UNIQUE | email |
| roles | UNIQUE | name |
| project_memberships | UNIQUE | (user_id, project_id) |
| project_memberships | INDEX | user_id |
| project_memberships | INDEX | project_id |
| project_memberships | INDEX | (project_id, user_id, role_id) | <!-- 新增：优化 JOIN 查询 -->
| audit_logs | INDEX | user_id |
| audit_logs | INDEX | event_type |
| audit_logs | INDEX | created_at |
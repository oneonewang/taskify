# Phase 0 Research: 用户管理与 RBAC

## Q1: 会话存储方案

### Decision
内存存储（默认），支持未来扩展到 Redis。

### Rationale
- 当前是 MVP 阶段，内存存储实现最简单
- SQLite 不适合高并发会话读写（文件锁竞争）
- Redis 需要额外基础设施，不适合当前阶段
- Go 标准库 `net/http` 的内置会话机制配合 gorilla/sessions 可满足需求

### Alternatives Considered
| 方案 | 拒绝原因 |
|------|----------|
| Redis | 需要额外服务，MVP 阶段不必要 |
| SQLite | 高并发下文件锁成为瓶颈 |
| PostgreSQL | 同上，且增加数据库复杂度 |

---

## Q2: 密码哈希 cost 参数

### Decision
bcrypt cost = 12

### Rationale
- bcrypt cost 12 是行业标准推荐值（OWASP 2015 建议）
- 平衡安全性与性能：约 250ms 哈希时间
- cost 每增加 1，计算时间翻倍
- 最低不应低于 10（2025 年标准）

### Alternatives Considered
| cost | 拒绝原因 |
|------|----------|
| 10 | 2025 年 NIST 标准认为不够 |
| 14+ | 计算时间过长（>1s），影响用户体验 |

---

## Q3: 预定义角色权限设计

### Decision
采用系统级 + 项目级双层角色：

**系统级预定义角色**：
| 角色 | 权限描述 |
|------|----------|
| 系统管理员 | 全局管理，可管理所有用户和角色 |

**项目级预定义角色**：
| 角色 | 权限描述 |
|------|----------|
| 项目 Owner | 完整项目控制权，管理成员 |
| 成员 (Member) | 创建/编辑/删除任务和评论 |
| 访客 (Guest) | 仅查看 |

**权限粒度设计**：
```
resources: users, projects, tasks, comments, roles
actions: view, create, edit, delete, manage
```

权限格式：`{resource}.{action}` (如 `tasks.create`, `projects.delete`)

### Rationale
- 符合 spec 中澄清记录的设计决策
- 项目级角色通过 project_members 表关联用户和项目
- 支持最小权限原则

---

## Q4: 审计日志表设计

### Decision
硬删除（HARD DELETE），但使用独立 archive 表保留历史。

### Rationale
- 审计日志通常需要长期保留
- 使用 GORM 的 `deleted_at` 软删除会导致查询复杂
- 采用独立 archive 表方案：主表定期归档到 archive

### Schema
```go
type AuditLog struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"index"`
    EventType string    `gorm:"index"` // login, logout, permission_change, access_denied
    Details   string    `gorm:"type:text"`
    IPAddress string    `gorm:"size:45"`
    CreatedAt time.Time
}
```

---

## 总结

| 问题 | 决定 |
|------|------|
| 会话存储 | 内存存储 (gorilla/sessions) |
| bcrypt cost | 12 |
| 角色设计 | 系统级(管理员) + 项目级(Owner/Member/Guest) |
| 审计日志 | 硬删除 + 定期归档策略 |
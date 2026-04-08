# 研究报告: Taskify 技术选型

**日期**: 2026-04-08
**Branch**: 001-kanban-platform

---

## 1. 实时更新方案: SSE vs WebSocket

### Decision: **SSE (Server-Sent Events)**

### Rationale
- 单向数据流(服务器推送)完美契合看板场景
- 浏览器原生支持，自动重连
- 比WebSocket更简单，资源消耗更低
- 防火墙友好，使用标准HTTP

### Alternatives Considered
| 方案 | 状态 | 原因 |
|------|------|------|
| WebSocket | 拒绝 | 双向通信过于复杂，适合聊天等场景 |
| Long Polling | 拒绝 | 延迟高，开销大 |
| GraphQL Subscriptions | 拒绝 | 过度设计 |
| Firebase/Realtime DB | 拒绝 | 外部依赖，增加延迟 |

### Implementation Notes
- Go后端: goroutine安全的SSEManager，支持多客户端管理
- Vue前端: EventSource API + Vue Composable封装
- 心跳机制: 30秒间隔防止连接超时
- 事件类型: `task_created`, `task_updated`, `task_deleted`, `task_moved`

---

## 2. 前端拖放方案: vue-draggable-plus

### Decision: **vue-draggable-plus + Pinia + 乐观更新**

### Rationale
- Vue 3 Composition API原生支持
- SortableJS封装，最成熟稳定的拖放库
- 支持跨列拖放(看板核心需求)
- 支持触摸设备和动画

### Alternatives Considered
| 方案 | 状态 | 原因 |
|------|------|------|
| vuedraggable (v2) | 拒绝 | Vue 2专用，已停止维护 |
| vue.draggable.next | 拒绝 | 作者已归档，推荐vue-draggable-plus |
| Native HTML5 DnD | 拒绝 | 触摸支持差，边缘情况复杂 |
| dnd-kit | 拒绝 | React专用，无官方Vue支持 |

### Implementation Notes
- `VueDraggable` 组件，`group="kanban"` 支持跨列拖放
- Pinia store管理状态，乐观更新模式
- API失败时自动回滚本地状态
- 300ms防抖处理快速连续拖放

---

## 3. Go后端测试策略

### Decision: **go test + httptest + SQLite内存数据库**

### Rationale
- Go标准库testing.T + testify/assert提供清晰断言
- httptest测试完整HTTP管道(路由、中间件、绑定、验证)
- SQLite内存数据库提供测试隔离，最快速度

### Alternatives Considered
| 方案 | 状态 | 原因 |
|------|------|------|
| testcontainers | 拒绝 | 需要Docker，测试速度慢 |
| go-sqlite3 (cgo) | 拒绝 | C依赖，CI/CD复杂 |
| bbolt/embedded DB | 拒绝 | SQL语法差异 |

### Implementation Notes
- 单元测试: Mock Repository，服务逻辑隔离测试
- Handler测试: httptest + gin.CreateTestContext
- 集成测试: testify/suite组织相关测试
- 覆盖率目标: Handlers 90%+, Services 80%+, Repos 70%+

---

## 4. 技术栈汇总

| 层级 | 技术选型 | 理由 |
|------|----------|------|
| 后端框架 | Gin | 轻量、高性能、Go生态主流 |
| ORM | Gorm | Go最流行ORM，支持SQLite/PostgreSQL |
| 数据库(开发) | SQLite | 零配置，文件数据库 |
| 数据库(生产) | PostgreSQL | 更好的并发支持 |
| 前端框架 | Vue.js 3 | 渐进式框架，Composition API |
| UI组件库 | Element Plus | Vue 3生态成熟组件库 |
| 拖放库 | vue-draggable-plus | Vue 3最佳拖放方案 |
| 状态管理 | Pinia | Vue 3官方推荐 |
| 实时通信 | SSE | 简单、可靠、适合单向推送 |
| E2E测试 | Playwright | 现代化E2E测试框架 |
| 单元测试 | go test + testify | Go标准+主流断言库 |

---

## 5. 性能目标

| 指标 | 目标 | 说明 |
|------|------|------|
| API响应时间 | <100ms (p95) | 简单CRUD操作 |
| 看板拖放帧率 | 60fps | 流畅用户体验 |
| SSE延迟 | <500ms | 实时更新感知 |
| 并发用户 | 5-10人 | 当前设计目标 |
| 最大任务数 | ~500 | 预估规模 |

---

## 6. 已知约束

| 约束 | 来源 | 说明 |
|------|------|------|
| 本地会话持久化 | spec | 刷新保持，关闭重置 |
| 无密码认证 | spec | 通过预定义用户列表选择身份 |
| 单体应用 | spec | 非微服务，章程违规已记录 |

---

## 7. 未解决问题

无。所有技术选型已确定。

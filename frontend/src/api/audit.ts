import axios from 'axios'

const apiClient = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  },
  withCredentials: true
})

// API响应格式
export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

// 审计日志类型
export interface AuditLog {
  id: number
  user_id: number
  user_email?: string
  event_type: string
  details: string
  ip_address: string
  created_at: string
}

// 分页类型
export interface Pagination {
  total: number
  page: number
  page_size: number
}

// 审计日志列表响应
export interface AuditLogListResponse {
  logs: AuditLog[]
  pagination: Pagination
}

// 查询参数
export interface AuditLogQuery {
  user_id?: number
  event_type?: string
  from?: string
  to?: string
  page?: number
  page_size?: number
}

// 获取审计日志
export async function getAuditLogs(query: AuditLogQuery = {}): Promise<ApiResponse<AuditLogListResponse>> {
  const res = await apiClient.get('/admin/audit-logs', { params: query })
  return res.data
}

// 事件类型选项
export const EVENT_TYPES = [
  { value: 'login', label: '登录' },
  { value: 'logout', label: '登出' },
  { value: 'permission_change', label: '权限变更' },
  { value: 'access_denied', label: '访问被拒绝' },
  { value: 'project_created', label: '项目创建' },
  { value: 'project_deleted', label: '项目删除' },
  { value: 'project_archived', label: '项目归档' },
  { value: 'member_added', label: '成员添加' },
  { value: 'member_removed', label: '成员移除' },
  { value: 'member_role_changed', label: '成员角色变更' }
]

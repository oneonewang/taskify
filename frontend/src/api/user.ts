import axios from 'axios'

const apiClient = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  },
  withCredentials: true // 发送cookies
})

// API响应格式
export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

// 用户类型
export interface User {
  id: number
  email: string
  display_name: string
  avatar_url: string
  email_verified: boolean
  is_disabled: boolean
  last_login_at?: string
  created_at: string
}

// 用户查询参数
export interface UserQueryParams {
  email?: string
  display_name?: string
  is_disabled?: string
  page?: number
  page_size?: number
}

// 用户列表响应
export interface UserListResponse {
  users: User[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

// 获取当前用户
export async function getCurrentUser(): Promise<ApiResponse<User>> {
  const res = await apiClient.get('/users/me')
  return res.data
}

// 获取用户列表（支持筛选和分页）
export async function getUsers(params?: UserQueryParams): Promise<ApiResponse<UserListResponse>> {
  const res = await apiClient.get('/admin/users', { params })
  return res.data
}

// 重置用户密码
export async function resetUserPassword(userId: number): Promise<ApiResponse<null>> {
  const res = await apiClient.post(`/admin/users/${userId}/reset-password`)
  return res.data
}

// 禁用用户
export async function disableUser(userId: number): Promise<ApiResponse<null>> {
  const res = await apiClient.post(`/admin/users/${userId}/disable`)
  return res.data
}

// 启用用户
export async function enableUser(userId: number): Promise<ApiResponse<null>> {
  const res = await apiClient.post(`/admin/users/${userId}/enable`)
  return res.data
}

// 导入用户结果
export interface ImportResult {
  total: number
  success: number
  failed: number
  errors?: string[]
}

// 下载导入模板
export function getImportTemplateUrl(): string {
  return '/api/admin/users/import/template'
}

// 导入用户
export async function importUsers(file: File): Promise<ApiResponse<ImportResult>> {
  const formData = new FormData()
  formData.append('file', file)
  const res = await apiClient.post('/admin/users/import', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
  return res.data
}
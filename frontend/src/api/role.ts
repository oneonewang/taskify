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

// 权限类型
export interface Permission {
  id: number
  resource: string
  action: string
  description: string
}

// 角色类型
export interface Role {
  id: number
  name: string
  description: string
  is_system: boolean
  scope: string
  permissions?: string[]
}

// 获取所有角色
export async function getRoles(): Promise<ApiResponse<{ roles: Role[] }>> {
  const res = await apiClient.get('/admin/roles')
  return res.data
}

// 获取角色详情
export async function getRole(id: number): Promise<ApiResponse<Role>> {
  const res = await apiClient.get(`/admin/roles/${id}`)
  return res.data
}

// 创建角色
export async function createRole(data: { name: string; description?: string; scope?: string }): Promise<ApiResponse<Role>> {
  const res = await apiClient.post('/admin/roles', data)
  return res.data
}

// 更新角色
export async function updateRole(id: number, data: { name?: string; description?: string }): Promise<ApiResponse<Role>> {
  const res = await apiClient.put(`/admin/roles/${id}`, data)
  return res.data
}

// 删除角色
export async function deleteRole(id: number): Promise<ApiResponse<null>> {
  const res = await apiClient.delete(`/admin/roles/${id}`)
  return res.data
}

// 设置角色权限
export async function setRolePermissions(roleId: number, permissionIds: number[]): Promise<ApiResponse<null>> {
  const res = await apiClient.put(`/admin/roles/${roleId}/permissions`, { permission_ids: permissionIds })
  return res.data
}

// 获取所有权限
export async function getPermissions(): Promise<ApiResponse<{ permissions: Permission[] }>> {
  const res = await apiClient.get('/admin/permissions')
  return res.data
}

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

// 用户角色信息
export interface UserRole {
  system_roles: string[]
  project_memberships: ProjectMembership[]
}

// 项目成员资格
export interface ProjectMembership {
  project_id: number
  project_name: string
  role_id: number
  role_name: string
  role_display_name: string
  joined_at: string
}

// 项目成员
export interface ProjectMember {
  user_id: number
  email: string
  display_name: string
  avatar_url: string
  role_id: number
  role_name: string
  role_display_name: string
  joined_at: string
}

// 获取用户的角色
export async function getUserRoles(userId: number): Promise<ApiResponse<UserRole>> {
  const res = await apiClient.get(`/admin/users/${userId}/roles`)
  return res.data
}

// 分配系统角色
export async function assignSystemRole(userId: number, roleId: number): Promise<ApiResponse<null>> {
  const res = await apiClient.post(`/admin/users/${userId}/roles`, { role_id: roleId })
  return res.data
}

// 移除系统角色
export async function removeSystemRole(userId: number, roleId: number): Promise<ApiResponse<null>> {
  const res = await apiClient.delete(`/admin/users/${userId}/roles/${roleId}`)
  return res.data
}

// 获取用户的项目成员资格
export async function getUserProjectMemberships(userId: number): Promise<ApiResponse<UserRole>> {
  const res = await apiClient.get(`/admin/users/${userId}/project-memberships`)
  return res.data
}

// 获取项目成员列表
export async function getProjectMembers(projectId: number): Promise<ApiResponse<{ members: ProjectMember[] }>> {
  const res = await apiClient.get(`/projects/${projectId}/members`)
  return res.data
}

// 添加项目成员
export async function addProjectMember(projectId: number, userEmail: string, roleId: number): Promise<ApiResponse<null>> {
  const res = await apiClient.post(`/projects/${projectId}/members`, { user_email: userEmail, role_id: roleId })
  return res.data
}

// 更新项目成员角色
export async function updateMemberRole(projectId: number, userId: number, roleId: number): Promise<ApiResponse<null>> {
  const res = await apiClient.put(`/projects/${projectId}/members/${userId}`, { role_id: roleId })
  return res.data
}

// 移除项目成员
export async function removeProjectMember(projectId: number, userId: number): Promise<ApiResponse<null>> {
  const res = await apiClient.delete(`/projects/${projectId}/members/${userId}`)
  return res.data
}

// 获取当前用户在项目中的成员资格
export async function getMyProjectMembership(projectId: number): Promise<ApiResponse<ProjectMember | null>> {
  const res = await apiClient.get(`/projects/${projectId}/my-membership`)
  return res.data
}

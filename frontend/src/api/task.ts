import axios from 'axios'
import type { ApiResponse, Task, Assignee } from './project'

const apiClient = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  },
  withCredentials: true // 发送cookies
})

// 创建任务请求
export interface CreateTaskRequest {
  title: string
  description?: string
  assignee_id: number
}

// 更新任务请求
export interface UpdateTaskRequest {
  title?: string
  description?: string
  assignee_id?: number
}

// 更新任务状态请求
export interface UpdateTaskStatusRequest {
  status: string
  position: number
}

// 创建任务
export function createTask(projectId: number, data: CreateTaskRequest): Promise<ApiResponse<Task>> {
  return apiClient.post(`/projects/${projectId}/tasks`, data).then(res => res.data)
}

// 更新任务
export function updateTask(taskId: number, data: UpdateTaskRequest): Promise<ApiResponse<Task>> {
  return apiClient.put(`/tasks/${taskId}`, data).then(res => res.data)
}

// 更新任务状态 (看板拖放)
export function updateTaskStatus(taskId: number, data: UpdateTaskStatusRequest): Promise<ApiResponse<{
  id: number
  status: string
  position: number
  updated_at: string
}>> {
  return apiClient.put(`/tasks/${taskId}/status`, data).then(res => res.data)
}

// 删除任务
export function deleteTask(taskId: number): Promise<ApiResponse<null>> {
  return apiClient.delete(`/tasks/${taskId}`).then(res => res.data)
}

// 获取用户列表 (用于任务分配)
export function getUsers(): Promise<ApiResponse<Assignee[]>> {
  return apiClient.get('/users').then(res => res.data)
}

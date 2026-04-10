import axios from 'axios'

const apiClient = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  },
  withCredentials: true // 发送cookies
})

export interface Assignee {
  id: number
  name: string
  avatar: string
}

export interface Task {
  id: number
  title: string
  description?: string
  status: string
  position: number
  assignee: Assignee
  project_id: number
  created_at: string
  updated_at: string
}

export interface TaskCount {
  todo: number
  in_progress: number
  review: number
  done: number
}

export interface Project {
  id: number
  name: string
  description: string
  is_archived?: boolean
  created_at: string
  task_counts?: TaskCount
  tasks?: Task[]
}

export interface ProjectDetail extends Project {
  task_counts: TaskCount
}

export interface CreateProjectRequest {
  name: string
  description?: string
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

// 后端 API 响应格式 (code: 0 表示成功)
export interface BackendApiResponse<T> {
  code: number
  message: string
  data: T
}

export function getProjects(): Promise<ApiResponse<Project[]>> {
  return apiClient.get('/projects').then(res => res.data)
}

export function getProject(id: number): Promise<ApiResponse<ProjectDetail>> {
  return apiClient.get(`/projects/${id}`).then(res => res.data)
}

export function getProjectTasks(projectId: number, status?: string): Promise<ApiResponse<Task[]>> {
  const params = status ? { status } : {}
  return apiClient.get(`/projects/${projectId}/tasks`, { params }).then(res => res.data)
}

export function createProject(data: CreateProjectRequest): Promise<ApiResponse<Project>> {
  return apiClient.post('/projects', data).then(res => res.data)
}

export function updateProject(id: number, data: { name: string; description: string }): Promise<ApiResponse<Project>> {
  return apiClient.put(`/projects/${id}`, data).then(res => res.data)
}

export function archiveProject(id: number): Promise<ApiResponse<Project>> {
  return apiClient.post(`/projects/${id}/archive`).then(res => res.data)
}

export function deleteProject(id: number): Promise<ApiResponse<null>> {
  return apiClient.delete(`/projects/${id}`).then(res => res.data)
}

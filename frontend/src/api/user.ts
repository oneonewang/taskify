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
  last_login_at?: string
  created_at: string
}

// 获取当前用户
export async function getCurrentUser(): Promise<ApiResponse<User>> {
  const res = await apiClient.get('/users/me')
  return res.data
}

// 获取所有用户
export async function getUsers(): Promise<ApiResponse<User[]>> {
  const res = await apiClient.get('/admin/users')
  return res.data
}
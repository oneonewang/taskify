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
  is_admin: boolean
  last_login_at?: string
  created_at: string
}

// 注册请求
export interface RegisterRequest {
  email: string
  password: string
  display_name: string
}

// 登录请求
export interface LoginRequest {
  email: string
  password: string
}

// 修改密码请求
export interface ChangePasswordRequest {
  current_password: string
  new_password: string
}

// 注册
export async function register(data: RegisterRequest): Promise<ApiResponse<{ user: User }>> {
  const res = await apiClient.post('/auth/register', data)
  return res.data
}

// 登录
export async function login(data: LoginRequest): Promise<ApiResponse<{ user: User }>> {
  const res = await apiClient.post('/auth/login', data)
  return res.data
}

// 登出
export async function logout(): Promise<ApiResponse<null>> {
  const res = await apiClient.post('/auth/logout')
  return res.data
}

// 获取当前用户
export async function getCurrentUser(): Promise<ApiResponse<User>> {
  const res = await apiClient.get('/users/me')
  return res.data
}

// 更新当前用户资料
export async function updateProfile(data: { display_name?: string; avatar_url?: string }): Promise<ApiResponse<User>> {
  const res = await apiClient.put('/users/me', data)
  return res.data
}

// 修改密码
export async function changePassword(data: ChangePasswordRequest): Promise<ApiResponse<null>> {
  const res = await apiClient.put('/users/me/password', data)
  return res.data
}
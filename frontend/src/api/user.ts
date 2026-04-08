import axios from 'axios'

const apiClient = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

export interface User {
  id: number
  name: string
  role: 'product_manager' | 'engineer'
  avatar: string
}

export interface ApiResponse<T> {
  success: boolean
  data: T
  message?: string
  error?: {
    code: string
    message: string
    details?: any
  }
}

export function getUsers(): Promise<ApiResponse<User[]>> {
  return apiClient.get('/users').then(res => res.data)
}

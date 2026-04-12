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

// 令牌信息
export interface TokenInfo {
  id: number
  name: string
  scope: string
  token_prefix: string
  created_at: string
  last_used_at?: string
  expires_at: string
  revoked: boolean
}

// 创建令牌请求
export interface CreateTokenRequest {
  name: string
  scope: string
  expires_in_days?: number
}

// 创建令牌响应（包含明文令牌）
export interface CreateTokenResponse {
  id: number
  token: string
  name: string
  scope: string
  token_prefix: string
  created_at: string
  expires_at: string
}

// 令牌列表响应
export interface TokenListResponse {
  tokens: TokenInfo[]
}

// 创建令牌
export async function createToken(data: CreateTokenRequest): Promise<ApiResponse<CreateTokenResponse>> {
  const res = await apiClient.post('/oauth/tokens', data)
  return res.data
}

// 获取令牌列表
export async function getTokens(): Promise<ApiResponse<TokenListResponse>> {
  const res = await apiClient.get('/oauth/tokens')
  return res.data
}

// 撤销令牌
export async function revokeToken(id: number): Promise<ApiResponse<null>> {
  const res = await apiClient.delete(`/oauth/tokens/${id}`)
  return res.data
}

// 验证令牌（公开端点）
export async function validateToken(): Promise<ApiResponse<{ active: boolean; user_id?: number; username?: string; scope?: string }>> {
  const res = await apiClient.get('/oauth/token/info')
  return res.data
}
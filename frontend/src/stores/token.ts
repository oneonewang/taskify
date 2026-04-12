import { defineStore } from 'pinia'
import { ref } from 'vue'
import { createToken as apiCreateToken, getTokens as apiGetTokens, revokeToken as apiRevokeToken, type TokenInfo, type CreateTokenRequest } from '../api/token'

export const useTokenStore = defineStore('token', () => {
  const tokens = ref<TokenInfo[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const lastCreatedToken = ref<string | null>(null)

  // 获取令牌列表
  async function fetchTokens() {
    loading.value = true
    error.value = null
    try {
      const res = await apiGetTokens()
      if (res.code === 0) {
        tokens.value = res.data.tokens
      } else {
        error.value = res.message
      }
    } catch (e: any) {
      error.value = e.response?.data?.message || '获取令牌列表失败'
    } finally {
      loading.value = false
    }
  }

  // 创建令牌
  async function createToken(data: CreateTokenRequest): Promise<{ success: boolean; token?: string }> {
    loading.value = true
    error.value = null
    lastCreatedToken.value = null
    try {
      const res = await apiCreateToken(data)
      if (res.code === 0) {
        // 将新令牌添加到列表
        const newToken: TokenInfo = {
          id: res.data.id,
          name: res.data.name,
          scope: res.data.scope,
          token_prefix: res.data.token_prefix,
          created_at: res.data.created_at,
          expires_at: res.data.expires_at,
          revoked: false
        }
        tokens.value.unshift(newToken)
        // 保存明文令牌（仅此时可见）
        lastCreatedToken.value = res.data.token
        return { success: true, token: res.data.token }
      } else {
        error.value = res.message
        return { success: false }
      }
    } catch (e: any) {
      error.value = e.response?.data?.message || '创建令牌失败'
      return { success: false }
    } finally {
      loading.value = false
    }
  }

  // 撤销令牌
  async function revokeToken(id: number): Promise<boolean> {
    loading.value = true
    error.value = null
    try {
      const res = await apiRevokeToken(id)
      if (res.code === 0) {
        // 更新本地状态
        const token = tokens.value.find(t => t.id === id)
        if (token) {
          token.revoked = true
        }
        return true
      } else {
        error.value = res.message
        return false
      }
    } catch (e: any) {
      error.value = e.response?.data?.message || '撤销令牌失败'
      return false
    } finally {
      loading.value = false
    }
  }

  // 清除错误
  function clearError() {
    error.value = null
  }

  // 清除最后创建的令牌（安全考虑）
  function clearLastCreatedToken() {
    lastCreatedToken.value = null
  }

  return {
    tokens,
    loading,
    error,
    lastCreatedToken,
    fetchTokens,
    createToken,
    revokeToken,
    clearError,
    clearLastCreatedToken
  }
})
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as apiLogin, logout as apiLogout, getCurrentUser as apiGetCurrentUser, register as apiRegister, type User } from '../api/auth'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const isLoggedIn = computed(() => user.value !== null)
  const isAdmin = computed(() => user.value?.is_admin === true)

  // 初始化 - 从服务器获取当前用户
  async function init() {
    loading.value = true
    error.value = null
    try {
      const res = await apiGetCurrentUser()
      if (res.code === 0) {
        user.value = res.data
      }
    } catch (e: any) {
      user.value = null
    } finally {
      loading.value = false
    }
  }

  // 登录
  async function login(email: string, password: string) {
    loading.value = true
    error.value = null
    try {
      const res = await apiLogin({ email, password })
      if (res.code === 0) {
        user.value = res.data.user
        return true
      } else {
        error.value = res.message
        return false
      }
    } catch (e: any) {
      error.value = e.response?.data?.message || '登录失败'
      return false
    } finally {
      loading.value = false
    }
  }

  // 注册
  async function register(email: string, password: string, displayName: string) {
    loading.value = true
    error.value = null
    try {
      const res = await apiRegister({ email, password, display_name: displayName })
      if (res.code === 0) {
        user.value = res.data.user
        return true
      } else {
        error.value = res.message
        return false
      }
    } catch (e: any) {
      error.value = e.response?.data?.message || '注册失败'
      return false
    } finally {
      loading.value = false
    }
  }

  // 登出
  async function logout() {
    loading.value = true
    error.value = null
    try {
      await apiLogout()
    } catch (e: any) {
      // Ignore API errors, still clear local state
    } finally {
      user.value = null
      loading.value = false
    }
  }

  // 清除错误
  function clearError() {
    error.value = null
  }

  return {
    user,
    loading,
    error,
    isLoggedIn,
    isAdmin,
    init,
    login,
    register,
    logout,
    clearError
  }
})
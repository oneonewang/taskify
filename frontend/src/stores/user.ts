import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface User {
  id: number
  name: string
  role: 'product_manager' | 'engineer'
  avatar: string
}

export const useUserStore = defineStore('user', () => {
  // 从localStorage恢复用户选择
  const storedUserId = localStorage.getItem('currentUserId')
  const currentUserId = ref<number | null>(storedUserId ? parseInt(storedUserId) : null)

  const currentUser = ref<User | null>(null)

  const isLoggedIn = computed(() => currentUserId.value !== null)

  function setCurrentUser(user: User) {
    currentUser.value = user
    currentUserId.value = user.id
    localStorage.setItem('currentUserId', user.id.toString())
  }

  function clearCurrentUser() {
    currentUser.value = null
    currentUserId.value = null
    localStorage.removeItem('currentUserId')
  }

  return {
    currentUserId,
    currentUser,
    isLoggedIn,
    setCurrentUser,
    clearCurrentUser
  }
})

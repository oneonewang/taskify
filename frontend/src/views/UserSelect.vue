<template>
  <div class="user-select-container">
    <div class="user-select-card">
      <h1 class="title">Taskify 看板平台</h1>
      <p class="subtitle">请选择当前用户</p>

      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      <div v-else class="user-list">
        <div
          v-for="user in users"
          :key="user.id"
          class="user-item"
          @click="selectUser(user)"
        >
          <div class="avatar" :style="{ backgroundColor: user.avatar }">
            {{ user.name.charAt(0) }}
          </div>
          <div class="user-info">
            <div class="user-name">{{ user.name }}</div>
            <div class="user-role">{{ roleLabel(user.role) }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getUsers, type User } from '../api/user'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

const users = ref<User[]>([])
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    const res = await getUsers()
    if (res.success) {
      users.value = res.data
    } else {
      error.value = res.error?.message || '获取用户列表失败'
    }
  } catch (e: any) {
    error.value = e.message || '网络错误'
  } finally {
    loading.value = false
  }
})

function selectUser(user: User) {
  userStore.setCurrentUser(user)
  router.push('/projects')
}

function roleLabel(role: string): string {
  return role === 'product_manager' ? '产品经理' : '工程师'
}
</script>

<style scoped>
.user-select-container {
  width: 100%;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.user-select-card {
  background: white;
  border-radius: 16px;
  padding: 48px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  max-width: 480px;
  width: 90%;
}

.title {
  font-size: 28px;
  font-weight: 600;
  color: #333;
  text-align: center;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 16px;
  color: #666;
  text-align: center;
  margin-bottom: 32px;
}

.loading, .error {
  text-align: center;
  padding: 32px;
  color: #666;
}

.error {
  color: #e6a23c;
}

.user-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.user-item {
  display: flex;
  align-items: center;
  padding: 16px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  border: 2px solid transparent;
}

.user-item:hover {
  background: #f5f7fa;
  border-color: #667eea;
}

.avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 20px;
  font-weight: 600;
  margin-right: 16px;
}

.user-name {
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.user-role {
  font-size: 14px;
  color: #999;
  margin-top: 4px;
}
</style>

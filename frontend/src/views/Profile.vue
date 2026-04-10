<template>
  <div class="profile-page">
    <div class="page-header">
      <h1>个人资料</h1>
    </div>

    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>

    <div v-else class="profile-content">
      <!-- 用户信息卡片 -->
      <div class="profile-card">
        <div class="avatar-section">
          <div class="avatar">
            <img v-if="user?.avatar_url" :src="user.avatar_url" :alt="user.display_name" />
            <span v-else class="avatar-placeholder">{{ user?.display_name?.charAt(0) || 'U' }}</span>
          </div>
        </div>
        <div class="user-info">
          <h2>{{ user?.display_name }}</h2>
          <p class="email">{{ user?.email }}</p>
          <p class="status">
            <span :class="['status-badge', user?.email_verified ? 'verified' : 'unverified']">
              {{ user?.email_verified ? '已验证邮箱' : '未验证邮箱' }}
            </span>
          </p>
          <p class="member-since">注册时间: {{ formatDate(user?.created_at) }}</p>
          <p v-if="user?.last_login_at" class="last-login">上次登录: {{ formatDate(user.last_login_at) }}</p>
        </div>
      </div>

      <!-- 编辑资料表单 -->
      <div class="section-card">
        <h3>编辑资料</h3>
        <form @submit.prevent="handleUpdateProfile" class="profile-form">
          <div class="form-group">
            <label for="display_name">显示名称</label>
            <input
              id="display_name"
              v-model="profileForm.display_name"
              type="text"
              placeholder="请输入显示名称"
              maxlength="100"
            />
          </div>
          <div class="form-group">
            <label for="avatar_url">头像 URL</label>
            <input
              id="avatar_url"
              v-model="profileForm.avatar_url"
              type="url"
              placeholder="请输入头像图片地址"
            />
          </div>
          <div class="form-actions">
            <button type="submit" class="btn btn-primary" :disabled="savingProfile">
              {{ savingProfile ? '保存中...' : '保存资料' }}
            </button>
          </div>
        </form>
        <div v-if="profileSuccess" class="success-message">{{ profileSuccess }}</div>
        <div v-if="profileError" class="error-message">{{ profileError }}</div>
      </div>

      <!-- 修改密码表单 -->
      <div class="section-card">
        <h3>修改密码</h3>
        <form @submit.prevent="handleChangePassword" class="password-form">
          <div class="form-group">
            <label for="current_password">当前密码</label>
            <input
              id="current_password"
              v-model="passwordForm.current_password"
              type="password"
              placeholder="请输入当前密码"
              required
            />
          </div>
          <div class="form-group">
            <label for="new_password">新密码</label>
            <input
              id="new_password"
              v-model="passwordForm.new_password"
              type="password"
              placeholder="请输入新密码（至少8字符，需包含字母和数字）"
              required
              minlength="8"
            />
          </div>
          <div class="form-group">
            <label for="confirm_password">确认新密码</label>
            <input
              id="confirm_password"
              v-model="passwordForm.confirm_password"
              type="password"
              placeholder="请再次输入新密码"
              required
            />
          </div>
          <div class="form-actions">
            <button type="submit" class="btn btn-primary" :disabled="changingPassword">
              {{ changingPassword ? '修改中...' : '修改密码' }}
            </button>
          </div>
        </form>
        <div v-if="passwordSuccess" class="success-message">{{ passwordSuccess }}</div>
        <div v-if="passwordError" class="error-message">{{ passwordError }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { getCurrentUser, updateProfile, changePassword } from '../api/auth'
import type { ApiResponse, User } from '../api/auth'

const authStore = useAuthStore()

const user = ref<User | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)

// 编辑资料表单
const profileForm = reactive({
  display_name: '',
  avatar_url: ''
})

const savingProfile = ref(false)
const profileSuccess = ref('')
const profileError = ref('')

// 修改密码表单
const passwordForm = reactive({
  current_password: '',
  new_password: '',
  confirm_password: ''
})

const changingPassword = ref(false)
const passwordSuccess = ref('')
const passwordError = ref('')

onMounted(async () => {
  await loadUser()
})

async function loadUser() {
  loading.value = true
  error.value = null
  try {
    // 首先尝试从 store 获取用户信息
    if (authStore.user) {
      user.value = authStore.user
      profileForm.display_name = authStore.user.display_name || ''
      profileForm.avatar_url = authStore.user.avatar_url || ''
    }

    // 然后从服务器刷新用户信息
    const res = await getCurrentUser() as ApiResponse<User>
    if (res.code === 0) {
      user.value = res.data
      profileForm.display_name = res.data.display_name || ''
      profileForm.avatar_url = res.data.avatar_url || ''
      // 更新 store
      authStore.user = res.data
    } else {
      error.value = res.message
    }
  } catch (e: any) {
    error.value = e.response?.data?.message || '加载失败'
  } finally {
    loading.value = false
  }
}

async function handleUpdateProfile() {
  profileSuccess.value = ''
  profileError.value = ''

  if (!profileForm.display_name.trim()) {
    profileError.value = '显示名称不能为空'
    return
  }

  savingProfile.value = true
  try {
    const res = await updateProfile({
      display_name: profileForm.display_name,
      avatar_url: profileForm.avatar_url
    }) as ApiResponse<User>

    if (res.code === 0) {
      user.value = res.data
      profileSuccess.value = '资料更新成功'
      // 清除成功消息
      setTimeout(() => {
        profileSuccess.value = ''
      }, 3000)
    } else {
      profileError.value = res.message
    }
  } catch (e: any) {
    profileError.value = e.response?.data?.message || '更新失败'
  } finally {
    savingProfile.value = false
  }
}

async function handleChangePassword() {
  passwordSuccess.value = ''
  passwordError.value = ''

  if (passwordForm.new_password !== passwordForm.confirm_password) {
    passwordError.value = '两次输入的密码不一致'
    return
  }

  if (passwordForm.new_password.length < 8) {
    passwordError.value = '新密码至少8字符，需包含字母和数字'
    return
  }

  changingPassword.value = true
  try {
    const res = await changePassword({
      current_password: passwordForm.current_password,
      new_password: passwordForm.new_password
    }) as ApiResponse<null>

    if (res.code === 0) {
      passwordSuccess.value = '密码修改成功'
      // 清空表单
      passwordForm.current_password = ''
      passwordForm.new_password = ''
      passwordForm.confirm_password = ''
      // 清除成功消息
      setTimeout(() => {
        passwordSuccess.value = ''
      }, 3000)
    } else {
      passwordError.value = res.message
    }
  } catch (e: any) {
    passwordError.value = e.response?.data?.message || '修改失败'
  } finally {
    changingPassword.value = false
  }
}

function formatDate(dateStr?: string): string {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<style scoped>
.profile-page {
  max-width: 800px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h1 {
  margin: 0;
  font-size: 24px;
  color: #333;
}

.loading,
.error {
  padding: 40px;
  text-align: center;
}

.error {
  color: #dc3545;
}

.profile-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.profile-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  padding: 24px;
  display: flex;
  gap: 24px;
  align-items: center;
}

.avatar-section {
  flex-shrink: 0;
}

.avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  overflow: hidden;
  background: #e3f2fd;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  font-size: 40px;
  font-weight: bold;
  color: #1976d2;
}

.user-info h2 {
  margin: 0 0 8px;
  font-size: 20px;
  color: #333;
}

.email {
  color: #666;
  margin: 0 0 8px;
}

.status-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.status-badge.verified {
  background: #e8f5e9;
  color: #2e7d32;
}

.status-badge.unverified {
  background: #fff3e0;
  color: #ef6c00;
}

.member-since,
.last-login {
  color: #999;
  font-size: 14px;
  margin: 4px 0 0;
}

.section-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  padding: 24px;
}

.section-card h3 {
  margin: 0 0 16px;
  font-size: 18px;
  color: #333;
  padding-bottom: 12px;
  border-bottom: 1px solid #eee;
}

.profile-form,
.password-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.form-group label {
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.form-group input {
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.2s;
}

.form-group input:focus {
  outline: none;
  border-color: #1976d2;
}

.form-actions {
  margin-top: 8px;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.2s;
}

.btn-primary {
  background: #1976d2;
  color: #fff;
}

.btn-primary:hover {
  background: #1565c0;
}

.btn-primary:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.success-message {
  margin-top: 12px;
  padding: 12px;
  background: #e8f5e9;
  color: #2e7d32;
  border-radius: 4px;
  font-size: 14px;
}

.error-message {
  margin-top: 12px;
  padding: 12px;
  background: #ffebee;
  color: #c62828;
  border-radius: 4px;
  font-size: 14px;
}

@media (max-width: 600px) {
  .profile-card {
    flex-direction: column;
    text-align: center;
  }

  .profile-page {
    padding: 16px;
  }
}
</style>

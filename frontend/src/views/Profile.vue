<template>
  <div class="profile-page">
    <AppHeader title="个人资料" :show-back="true" @back="goBack" />

    <main class="profile-content">
      <div v-if="loading" class="loading-state animate-in">
        <div class="loading-spinner"></div>
        <span>加载中...</span>
      </div>
      <div v-else-if="error" class="error-state animate-in">
        <div class="error-icon">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="12"/>
            <line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
        </div>
        <span>{{ error }}</span>
      </div>

      <template v-else>
        <!-- User Info Card -->
        <div class="profile-card animate-in">
          <div class="avatar-section">
            <div class="avatar" :style="{ backgroundColor: avatarColor }">
              <img v-if="user?.avatar_url" :src="user.avatar_url" :alt="user.display_name" />
              <span v-else class="avatar-placeholder">{{ user?.display_name?.charAt(0) || 'U' }}</span>
            </div>
          </div>
          <div class="user-info">
            <h2>{{ user?.display_name }}</h2>
            <p class="email">{{ user?.email }}</p>
            <div class="status-row">
              <span :class="['status-badge', user?.email_verified ? 'verified' : 'unverified']">
                <svg v-if="user?.email_verified" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="20 6 9 17 4 12"/>
                </svg>
                {{ user?.email_verified ? '已验证邮箱' : '未验证邮箱' }}
              </span>
            </div>
            <div class="meta-info">
              <span class="meta-item">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
                  <line x1="16" y1="2" x2="16" y2="6"/>
                  <line x1="8" y1="2" x2="8" y2="6"/>
                  <line x1="3" y1="10" x2="21" y2="10"/>
                </svg>
                注册于 {{ formatDate(user?.created_at) }}
              </span>
              <span v-if="user?.last_login_at" class="meta-item">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <polyline points="12 6 12 12 16 14"/>
                </svg>
                上次登录 {{ formatDate(user.last_login_at) }}
              </span>
            </div>
          </div>
        </div>

        <!-- Edit Profile Section -->
        <div class="section-card animate-in animate-in-delay-1">
          <div class="section-header">
            <h3>编辑资料</h3>
            <p>更新您的个人信息和头像</p>
          </div>
          <form @submit.prevent="handleUpdateProfile" class="profile-form">
            <div class="form-row">
              <div class="form-group">
                <label for="display_name">显示名称</label>
                <input
                  id="display_name"
                  v-model="profileForm.display_name"
                  type="text"
                  placeholder="请输入显示名称"
                  maxlength="100"
                  class="form-input"
                />
              </div>
              <div class="form-group">
                <label for="avatar_url">头像 URL</label>
                <input
                  id="avatar_url"
                  v-model="profileForm.avatar_url"
                  type="url"
                  placeholder="请输入头像图片地址"
                  class="form-input"
                />
              </div>
            </div>
            <div class="form-actions">
              <button type="submit" class="btn btn-primary" :disabled="savingProfile">
                <svg v-if="!savingProfile" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/>
                  <polyline points="17 21 17 13 7 13 7 21"/>
                  <polyline points="7 3 7 8 15 8"/>
                </svg>
                {{ savingProfile ? '保存中...' : '保存资料' }}
              </button>
            </div>
          </form>
          <Transition name="fade">
            <div v-if="profileSuccess" class="success-message">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="20 6 9 17 4 12"/>
              </svg>
              {{ profileSuccess }}
            </div>
          </Transition>
          <Transition name="fade">
            <div v-if="profileError" class="error-message">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <line x1="12" y1="8" x2="12" y2="12"/>
                <line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
              {{ profileError }}
            </div>
          </Transition>
        </div>

        <!-- Change Password Section -->
        <div class="section-card animate-in animate-in-delay-2">
          <div class="section-header">
            <h3>修改密码</h3>
            <p>更新您的账户密码</p>
          </div>
          <form @submit.prevent="handleChangePassword" class="password-form">
            <div class="form-group">
              <label for="current_password">当前密码</label>
              <input
                id="current_password"
                v-model="passwordForm.current_password"
                type="password"
                placeholder="请输入当前密码"
                required
                class="form-input"
              />
            </div>
            <div class="form-row">
              <div class="form-group">
                <label for="new_password">新密码</label>
                <input
                  id="new_password"
                  v-model="passwordForm.new_password"
                  type="password"
                  placeholder="请输入新密码（至少8字符）"
                  required
                  minlength="8"
                  class="form-input"
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
                  class="form-input"
                />
              </div>
            </div>
            <div class="form-actions">
              <button type="submit" class="btn btn-primary" :disabled="changingPassword">
                <svg v-if="!changingPassword" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                </svg>
                {{ changingPassword ? '修改中...' : '修改密码' }}
              </button>
            </div>
          </form>
          <Transition name="fade">
            <div v-if="passwordSuccess" class="success-message">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="20 6 9 17 4 12"/>
              </svg>
              {{ passwordSuccess }}
            </div>
          </Transition>
          <Transition name="fade">
            <div v-if="passwordError" class="error-message">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <line x1="12" y1="8" x2="12" y2="12"/>
                <line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
              {{ passwordError }}
            </div>
          </Transition>
        </div>
      </template>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { getCurrentUser, updateProfile, changePassword } from '../api/auth'
import type { ApiResponse, User } from '../api/auth'
import AppHeader from '../components/AppHeader.vue'

const router = useRouter()
const authStore = useAuthStore()

const user = ref<User | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)

const profileForm = reactive({
  display_name: '',
  avatar_url: ''
})

const savingProfile = ref(false)
const profileSuccess = ref('')
const profileError = ref('')

const passwordForm = reactive({
  current_password: '',
  new_password: '',
  confirm_password: ''
})

const changingPassword = ref(false)
const passwordSuccess = ref('')
const passwordError = ref('')

const avatarColors = ['#5b5fc7', '#8b5cf6', '#ec4899', '#f97316', '#14b8a6', '#06b6d4', '#3b82f6', '#84cc16']
const avatarColor = computed(() => {
  const name = user.value?.display_name || user.value?.email || '?'
  const charCode = name.charCodeAt(0)
  return avatarColors[charCode % avatarColors.length]
})

onMounted(async () => {
  await loadUser()
})

function goBack() {
  router.push('/projects')
}

async function loadUser() {
  loading.value = true
  error.value = null
  try {
    if (authStore.user) {
      user.value = authStore.user
      profileForm.display_name = authStore.user.display_name || ''
      profileForm.avatar_url = authStore.user.avatar_url || ''
    }

    const res = await getCurrentUser() as ApiResponse<User>
    if (res.code === 0) {
      user.value = res.data
      profileForm.display_name = res.data.display_name || ''
      profileForm.avatar_url = res.data.avatar_url || ''
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
      setTimeout(() => { profileSuccess.value = '' }, 3000)
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
    passwordError.value = '新密码至少8字符'
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
      passwordForm.current_password = ''
      passwordForm.new_password = ''
      passwordForm.confirm_password = ''
      setTimeout(() => { passwordSuccess.value = '' }, 3000)
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
    day: 'numeric'
  })
}
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background: var(--color-bg-base);
}

.profile-content {
  max-width: 800px;
  margin: 0 auto;
  padding: var(--space-8) var(--space-6);
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

/* Loading/Error States */
.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-16);
  gap: var(--space-4);
  color: var(--color-text-muted);
}

.error-state {
  color: var(--color-danger);
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  background: rgba(239, 71, 111, 0.1);
  border-radius: 50%;
}

/* Profile Card */
.profile-card {
  background: var(--color-bg-surface);
  border-radius: var(--radius-xl);
  padding: var(--space-8);
  box-shadow: var(--shadow-card);
  display: flex;
  gap: var(--space-8);
  align-items: center;
  border: 1px solid var(--color-border);
}

.avatar-section {
  flex-shrink: 0;
}

.avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-md);
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  font-size: 48px;
  font-weight: bold;
  color: white;
}

.user-info {
  flex: 1;
}

.user-info h2 {
  font-family: var(--font-display);
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-text-primary);
  margin: 0 0 var(--space-2);
}

.email {
  color: var(--color-text-secondary);
  margin: 0 0 var(--space-3);
  font-size: var(--font-size-sm);
}

.status-row {
  margin-bottom: var(--space-4);
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  padding: var(--space-1) var(--space-3);
  border-radius: var(--radius-full);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
}

.status-badge.verified {
  background: rgba(78, 205, 196, 0.15);
  color: #2e8b7d;
}

.status-badge.unverified {
  background: rgba(255, 209, 102, 0.15);
  color: #b8860b;
}

.meta-info {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-4);
}

.meta-item {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

/* Section Card */
.section-card {
  background: var(--color-bg-surface);
  border-radius: var(--radius-xl);
  padding: var(--space-6);
  box-shadow: var(--shadow-card);
  border: 1px solid var(--color-border);
}

.section-header {
  margin-bottom: var(--space-6);
  padding-bottom: var(--space-4);
  border-bottom: 1px solid var(--color-border);
}

.section-header h3 {
  font-family: var(--font-display);
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  margin: 0 0 var(--space-1);
}

.section-header p {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  margin: 0;
}

/* Form Styles */
.profile-form,
.password-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-4);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.form-group label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-secondary);
}

.form-input {
  padding: var(--space-3) var(--space-4);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  font-family: var(--font-body);
  font-size: var(--font-size-base);
  color: var(--color-text-primary);
  background: var(--color-bg-surface);
  transition: all var(--transition-base);
}

.form-input:hover {
  border-color: var(--color-border-hover);
}

.form-input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(91, 95, 199, 0.1);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}

/* Buttons */
.btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-5);
  border: none;
  border-radius: var(--radius-lg);
  cursor: pointer;
  font-family: var(--font-body);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  transition: all var(--transition-base);
}

.btn-primary {
  background: var(--color-primary);
  color: white;
  box-shadow: 0 2px 8px rgba(91, 95, 199, 0.3);
}

.btn-primary:hover:not(:disabled) {
  background: var(--color-primary-hover);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(91, 95, 199, 0.4);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Messages */
.success-message,
.error-message {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  margin-top: var(--space-4);
  padding: var(--space-3) var(--space-4);
  border-radius: var(--radius-lg);
  font-size: var(--font-size-sm);
}

.success-message {
  background: rgba(78, 205, 196, 0.15);
  color: #2e8b7d;
}

.error-message {
  background: rgba(239, 71, 111, 0.1);
  color: var(--color-danger);
}

/* Fade Animation */
.fade-enter-active,
.fade-leave-active {
  transition: opacity var(--transition-base);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Responsive */
@media (max-width: 768px) {
  .profile-content {
    padding: var(--space-4);
  }

  .profile-card {
    flex-direction: column;
    text-align: center;
    padding: var(--space-6);
  }

  .avatar {
    width: 100px;
    height: 100px;
  }

  .avatar-placeholder {
    font-size: 36px;
  }

  .user-info h2 {
    font-size: var(--font-size-xl);
  }

  .meta-info {
    justify-content: center;
  }

  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>

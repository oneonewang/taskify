<template>
  <div class="register-page">
    <!-- Background decoration -->
    <div class="bg-decoration">
      <div class="bg-circle bg-circle-1"></div>
      <div class="bg-circle bg-circle-2"></div>
      <div class="bg-circle bg-circle-3"></div>
    </div>

    <!-- Register Card -->
    <div class="register-container animate-in">
      <div class="register-card">
        <!-- Brand Header -->
        <div class="card-header">
          <div class="brand-mark">
            <svg width="36" height="36" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
              <rect width="40" height="40" rx="12" fill="url(#register-brand-gradient)"/>
              <path d="M12 20L18 14L24 20L18 26L12 20Z" fill="white" fill-opacity="0.9"/>
              <path d="M18 14L26 14L26 20L18 26L18 14Z" fill="white" fill-opacity="0.6"/>
              <path d="M24 20L32 20L26 26L24 20Z" fill="white" fill-opacity="0.4"/>
              <defs>
                <linearGradient id="register-brand-gradient" x1="0" y1="0" x2="40" y2="40" gradientUnits="userSpaceOnUse">
                  <stop stop-color="#5b5fc7"/>
                  <stop offset="1" stop-color="#8b5cf6"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <h1 class="brand-title">创建账号</h1>
          <p class="brand-tagline">加入 Taskify，开始高效协作</p>
        </div>

        <!-- Register Form -->
        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-position="top"
          class="register-form"
          @submit.prevent="handleRegister"
        >
          <el-form-item label="显示名称" prop="displayName" class="form-item">
            <el-input
              v-model="form.displayName"
              placeholder="您希望如何称呼"
              :disabled="loading"
              class="custom-input"
            >
              <template #prefix>
                <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                  <circle cx="12" cy="7" r="4"/>
                </svg>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item label="邮箱" prop="email" class="form-item">
            <el-input
              v-model="form.email"
              type="email"
              placeholder="请输入工作邮箱"
              :disabled="loading"
              class="custom-input"
            >
              <template #prefix>
                <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
                  <polyline points="22,6 12,13 2,6"/>
                </svg>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item label="密码" prop="password" class="form-item">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="至少8字符，包含字母和数字"
              show-password
              :disabled="loading"
              class="custom-input"
            >
              <template #prefix>
                <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                </svg>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item label="确认密码" prop="confirmPassword" class="form-item">
            <el-input
              v-model="form.confirmPassword"
              type="password"
              placeholder="请再次输入密码"
              show-password
              :disabled="loading"
              class="custom-input"
            >
              <template #prefix>
                <svg class="input-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                  <polyline points="9 12 12 15 16 10"/>
                </svg>
              </template>
            </el-input>
          </el-form-item>

          <el-alert
            v-if="error"
            :title="error"
            type="error"
            show-icon
            :closable="false"
            class="error-alert"
          />

          <el-form-item class="submit-item">
            <el-button
              type="primary"
              native-type="submit"
              :loading="loading"
              class="submit-btn"
              @click="handleRegister"
            >
              <span v-if="!loading">注册</span>
              <span v-else>注册中...</span>
            </el-button>
          </el-form-item>
        </el-form>

        <!-- Footer -->
        <div class="register-footer">
          <span class="footer-text">已有账号？</span>
          <router-link to="/login" class="login-link">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="19" y1="12" x2="5" y2="12"/>
              <polyline points="12 19 5 12 12 5"/>
            </svg>
            立即登录
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const error = ref<string | null>(null)

const form = reactive({
  displayName: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const validateConfirmPassword = (_rule: any, value: string, callback: any) => {
  if (value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules: FormRules = {
  displayName: [
    { required: true, message: '请输入显示名称', trigger: 'blur' },
    { min: 1, max: 100, message: '显示名称长度在 1 到 100 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 8, message: '密码至少8个字符', trigger: 'blur' },
    { pattern: /^(?=.*[A-Za-z])(?=.*\d).+$/, message: '密码需包含字母和数字', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

async function handleRegister() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    loading.value = true
    error.value = null

    const success = await authStore.register(form.email, form.password, form.displayName)

    if (success) {
      ElMessage.success('注册成功')
      router.push('/projects')
    } else {
      error.value = authStore.error || '注册失败'
    }

    loading.value = false
  })
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-bg-base);
  position: relative;
  overflow: hidden;
  padding: var(--space-4);
}

/* Background Decoration */
.bg-decoration {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.bg-circle {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.5;
}

.bg-circle-1 {
  width: min(600px, 80vw);
  height: min(600px, 80vw);
  background: radial-gradient(circle, rgba(139, 92, 246, 0.15) 0%, transparent 70%);
  top: -20%;
  left: -10%;
  animation: float 8s ease-in-out infinite;
}

.bg-circle-2 {
  width: min(400px, 60vw);
  height: min(400px, 60vw);
  background: radial-gradient(circle, rgba(91, 95, 199, 0.12) 0%, transparent 70%);
  bottom: -10%;
  right: -10%;
  animation: float 10s ease-in-out infinite reverse;
}

.bg-circle-3 {
  width: min(300px, 50vw);
  height: min(300px, 50vw);
  background: radial-gradient(circle, rgba(255, 107, 107, 0.1) 0%, transparent 70%);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation: float 12s ease-in-out infinite;
}

/* Register Container */
.register-container {
  width: 100%;
  max-width: 380px;
  position: relative;
  z-index: 1;
}

/* Register Card */
.register-card {
  background: var(--color-bg-surface);
  border-radius: var(--radius-xl);
  padding: var(--space-6);
  box-shadow: var(--shadow-lg);
  border: 1px solid var(--color-border);
}

/* Card Header */
.card-header {
  text-align: center;
  margin-bottom: var(--space-5);
}

.brand-mark {
  display: inline-flex;
  margin-bottom: var(--space-3);
}

.brand-title {
  font-family: var(--font-display);
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-text-primary);
  margin: 0 0 var(--space-1);
  letter-spacing: -0.02em;
}

.brand-tagline {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
  margin: 0;
}

/* Form Styles */
.register-form {
  margin-bottom: var(--space-4);
}

.form-item {
  margin-bottom: var(--space-4);
}

.form-item :deep(.el-form-item__label) {
  font-weight: var(--font-weight-medium);
  color: var(--color-text-secondary);
  margin-bottom: var(--space-1);
  font-size: var(--font-size-sm);
}

.custom-input {
  font-family: var(--font-body);
}

.custom-input :deep(.el-input__wrapper) {
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-md);
  box-shadow: 0 0 0 1px var(--color-border);
  transition: all var(--transition-base);
}

.custom-input :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px var(--color-border-hover);
}

.custom-input :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px var(--color-primary-light);
}

.custom-input :deep(.el-input__inner) {
  font-family: var(--font-body);
  font-size: var(--font-size-sm);
}

.input-icon {
  color: var(--color-text-muted);
  margin-right: var(--space-1);
}

/* Error Alert */
.error-alert {
  margin-bottom: var(--space-4);
  border-radius: var(--radius-md);
  border: none;
}

/* Submit Button */
.submit-item {
  margin-bottom: 0;
}

.submit-btn {
  width: 100%;
  height: 44px;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  border-radius: var(--radius-md);
  background: linear-gradient(135deg, var(--color-primary) 0%, #7c5df0 100%);
  border: none;
  box-shadow: 0 2px 8px rgba(91, 95, 199, 0.3);
  transition: all var(--transition-base);
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(91, 95, 199, 0.4);
}

/* Footer */
.register-footer {
  text-align: center;
  padding-top: var(--space-4);
  border-top: 1px solid var(--color-border);
}

.footer-text {
  color: var(--color-text-muted);
  font-size: var(--font-size-xs);
}

.login-link {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: var(--color-primary);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
  text-decoration: none;
  margin-left: var(--space-1);
  transition: all var(--transition-base);
}

.login-link:hover {
  color: var(--color-primary-hover);
}

.login-link svg {
  transition: transform var(--transition-base);
}

.login-link:hover svg {
  transform: translateX(-2px);
}

/* Responsive */
@media (max-width: 480px) {
  .register-page {
    padding: var(--space-4);
  }

  .register-card {
    padding: var(--space-5);
    border-radius: var(--radius-xl);
  }

  .brand-title {
    font-size: var(--font-size-xl);
  }
}
</style>

<template>
  <div class="token-settings-page">
    <AppHeader title="令牌管理" :show-back="true" @back="goBack" />

    <main class="token-settings-content">
      <div v-if="loading" class="loading-state animate-in">
        <div class="loading-spinner"></div>
        <span>加载中...</span>
      </div>

      <template v-else>
        <!-- Introduction Card -->
        <div class="intro-card animate-in">
          <div class="intro-icon">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
              <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
            </svg>
          </div>
          <div class="intro-content">
            <h2>个人访问令牌</h2>
            <p>个人访问令牌用于 MCP 客户端认证。您可以创建多个令牌，每个令牌有不同的权限范围。</p>
          </div>
        </div>

        <!-- Token Manager Component -->
        <TokenManager class="animate-in animate-in-delay-1" />
      </template>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import TokenManager from '../components/token/TokenManager.vue'
import AppHeader from '../components/AppHeader.vue'

const router = useRouter()
const loading = ref(false)

onMounted(async () => {
  // Initial loading if needed
  loading.value = false
})

function goBack() {
  router.push('/projects')
}
</script>

<style scoped>
.token-settings-page {
  min-height: 100vh;
  background: var(--color-bg-base);
}

.token-settings-content {
  max-width: 800px;
  margin: 0 auto;
  padding: var(--space-8) var(--space-6);
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

/* Loading State */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-16);
  gap: var(--space-4);
  color: var(--color-text-muted);
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

/* Intro Card */
.intro-card {
  display: flex;
  gap: var(--space-5);
  padding: var(--space-6);
  background: linear-gradient(135deg, rgba(91, 95, 199, 0.1), rgba(91, 95, 199, 0.05));
  border-radius: var(--radius-xl);
  border: 1px solid rgba(91, 95, 199, 0.2);
}

.intro-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background: var(--color-primary);
  color: white;
  border-radius: var(--radius-lg);
  flex-shrink: 0;
}

.intro-content h2 {
  font-family: var(--font-display);
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  margin: 0 0 var(--space-2);
}

.intro-content p {
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  margin: 0;
  line-height: 1.6;
}

/* Responsive */
@media (max-width: 768px) {
  .token-settings-content {
    padding: var(--space-4);
  }

  .intro-card {
    flex-direction: column;
    text-align: center;
    align-items: center;
  }
}
</style>
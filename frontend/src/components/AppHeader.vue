<template>
  <header class="app-header">
    <div class="header-left">
      <!-- Back Button (optional) -->
      <button v-if="showBack" @click="goBack" class="back-btn">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="19" y1="12" x2="5" y2="12"/>
          <polyline points="12 19 5 12 12 5"/>
        </svg>
      </button>

      <!-- Logo / Brand -->
      <router-link to="/projects" class="brand-link">
        <div class="brand-mark">
          <svg width="28" height="28" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect width="40" height="40" rx="10" fill="url(#header-brand-gradient)"/>
            <path d="M12 20L18 14L24 20L18 26L12 20Z" fill="white" fill-opacity="0.9"/>
            <path d="M18 14L26 14L26 20L18 26L18 14Z" fill="white" fill-opacity="0.6"/>
            <defs>
              <linearGradient id="header-brand-gradient" x1="0" y1="0" x2="40" y2="40" gradientUnits="userSpaceOnUse">
                <stop stop-color="#5b5fc7"/>
                <stop offset="1" stop-color="#8b5cf6"/>
              </linearGradient>
            </defs>
          </svg>
        </div>
        <span v-if="!hideBrandText" class="brand-text">Taskify</span>
      </router-link>

      <!-- Page Title -->
      <div v-if="title" class="page-title-wrapper">
        <span class="title-separator">/</span>
        <h1 class="page-title">{{ title }}</h1>
      </div>
    </div>

    <div class="header-right">
      <!-- Project Selector (optional) -->
      <slot name="center"></slot>

      <!-- Action Buttons -->
      <slot name="actions"></slot>

      <!-- User Menu -->
      <div class="user-menu">
        <button class="user-btn" @click="toggleUserMenu">
          <span class="user-avatar" :style="{ backgroundColor: avatarColor }">
            {{ userInitial }}
          </span>
          <span v-if="!hideUserName" class="user-name">{{ displayName }}</span>
          <svg class="chevron" :class="{ open: userMenuOpen }" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="6 9 12 15 18 9"/>
          </svg>
        </button>

        <!-- Dropdown Menu -->
        <Transition name="dropdown">
          <div v-if="userMenuOpen" class="user-dropdown" @click.stop>
            <div class="dropdown-header">
              <span class="dropdown-name">{{ displayName }}</span>
              <span class="dropdown-email">{{ userEmail }}</span>
            </div>
            <div class="dropdown-divider"></div>
            <router-link to="/profile" class="dropdown-item" @click="closeUserMenu">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
              个人资料
            </router-link>
            <router-link v-if="authStore.isAdmin" to="/admin/users" class="dropdown-item" @click="closeUserMenu">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="3"/>
                <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
              </svg>
              管理后台
            </router-link>
            <div class="dropdown-divider"></div>
            <button class="dropdown-item logout" @click="handleLogout">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
                <polyline points="16 17 21 12 16 7"/>
                <line x1="21" y1="12" x2="9" y2="12"/>
              </svg>
              退出登录
            </button>
          </div>
        </Transition>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

defineProps<{
  title?: string
  showBack?: boolean
  hideBrandText?: boolean
  hideUserName?: boolean
}>()

const emit = defineEmits<{
  back: []
}>()

const router = useRouter()
const authStore = useAuthStore()

const userMenuOpen = ref(false)

const displayName = computed(() => {
  return authStore.user?.display_name || authStore.user?.email?.split('@')[0] || '用户'
})

const userEmail = computed(() => {
  return authStore.user?.email || ''
})

const userInitial = computed(() => {
  const name = authStore.user?.display_name || authStore.user?.email || '?'
  return name.charAt(0).toUpperCase()
})

const avatarColors = ['#5b5fc7', '#8b5cf6', '#ec4899', '#f97316', '#14b8a6', '#06b6d4', '#3b82f6', '#84cc16']
const avatarColor = computed(() => {
  const name = authStore.user?.display_name || authStore.user?.email || '?'
  const charCode = name.charCodeAt(0)
  return avatarColors[charCode % avatarColors.length]
})

function goBack() {
  emit('back')
  router.back()
}

function toggleUserMenu() {
  userMenuOpen.value = !userMenuOpen.value
}

function closeUserMenu() {
  userMenuOpen.value = false
}

async function handleLogout() {
  closeUserMenu()
  await authStore.logout()
  router.push('/login')
}

// Close menu when clicking outside
function handleClickOutside(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (!target.closest('.user-menu')) {
    userMenuOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.app-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-3) var(--space-6);
  background: var(--color-bg-surface);
  border-bottom: 1px solid var(--color-border);
  height: 64px;
  position: sticky;
  top: 0;
  z-index: var(--z-dropdown);
}

/* Header Left */
.header-left {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.back-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: var(--color-bg-muted);
  border: none;
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all var(--transition-base);
}

.back-btn:hover {
  background: var(--color-border);
  color: var(--color-text-primary);
}

.brand-link {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  text-decoration: none;
  transition: opacity var(--transition-base);
}

.brand-link:hover {
  opacity: 0.85;
}

.brand-text {
  font-family: var(--font-display);
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-bold);
  color: var(--color-text-primary);
}

.page-title-wrapper {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.title-separator {
  color: var(--color-text-muted);
  font-size: var(--font-size-lg);
}

.page-title {
  font-family: var(--font-display);
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  margin: 0;
}

/* Header Right */
.header-right {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

/* User Menu */
.user-menu {
  position: relative;
}

.user-btn {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2);
  background: transparent;
  border: none;
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--transition-base);
}

.user-btn:hover {
  background: var(--color-bg-muted);
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
}

.user-name {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-primary);
}

.chevron {
  color: var(--color-text-muted);
  transition: transform var(--transition-base);
}

.chevron.open {
  transform: rotate(180deg);
}

/* User Dropdown */
.user-dropdown {
  position: absolute;
  top: calc(100% + var(--space-2));
  right: 0;
  min-width: 220px;
  background: var(--color-bg-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-xl);
  padding: var(--space-2);
  z-index: var(--z-dropdown);
}

.dropdown-header {
  padding: var(--space-3);
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.dropdown-name {
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  font-size: var(--font-size-sm);
}

.dropdown-email {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
}

.dropdown-divider {
  height: 1px;
  background: var(--color-border);
  margin: var(--space-2) 0;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  width: 100%;
  padding: var(--space-3);
  background: none;
  border: none;
  border-radius: var(--radius-md);
  font-family: var(--font-body);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  text-decoration: none;
  cursor: pointer;
  transition: all var(--transition-base);
}

.dropdown-item:hover {
  background: var(--color-bg-muted);
  color: var(--color-text-primary);
}

.dropdown-item.logout {
  color: var(--color-danger);
}

.dropdown-item.logout:hover {
  background: rgba(239, 71, 111, 0.1);
  color: var(--color-danger);
}

/* Dropdown Animation */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all var(--transition-fast);
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* Responsive */
@media (max-width: 768px) {
  .app-header {
    padding: var(--space-3) var(--space-4);
  }

  .brand-text,
  .user-name,
  .title-separator,
  .page-title-wrapper {
    display: none;
  }
}
</style>

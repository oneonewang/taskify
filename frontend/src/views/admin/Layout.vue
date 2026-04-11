<template>
  <div class="admin-layout">
    <!-- Sidebar -->
    <aside class="admin-sidebar">
      <!-- Sidebar Header -->
      <div class="sidebar-header">
        <router-link to="/projects" class="brand-link">
          <div class="brand-mark">
            <svg width="32" height="32" viewBox="0 0 40 40" fill="none" xmlns="http://www.w3.org/2000/svg">
              <rect width="40" height="40" rx="10" fill="url(#admin-brand-gradient)"/>
              <path d="M12 20L18 14L24 20L18 26L12 20Z" fill="white" fill-opacity="0.9"/>
              <path d="M18 14L26 14L26 20L18 26L18 14Z" fill="white" fill-opacity="0.6"/>
              <defs>
                <linearGradient id="admin-brand-gradient" x1="0" y1="0" x2="40" y2="40" gradientUnits="userSpaceOnUse">
                  <stop stop-color="#5b5fc7"/>
                  <stop offset="1" stop-color="#8b5cf6"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <div class="brand-text">
            <span class="brand-name">Taskify</span>
            <span class="brand-badge">管理后台</span>
          </div>
        </router-link>
      </div>

      <!-- Navigation -->
      <nav class="sidebar-nav">
        <div class="nav-section">
          <span class="nav-section-title">系统管理</span>
          <router-link to="/admin/users" class="nav-item" :class="{ active: isActive('/admin/users') }">
            <div class="nav-icon-wrapper">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="9" cy="7" r="4"/>
                <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
              </svg>
            </div>
            <span class="nav-label">用户管理</span>
          </router-link>
          <router-link to="/admin/roles" class="nav-item" :class="{ active: isActive('/admin/roles') }">
            <div class="nav-icon-wrapper">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                <path d="M9 12l2 2 4-4"/>
              </svg>
            </div>
            <span class="nav-label">角色管理</span>
          </router-link>
          <router-link to="/admin/audit-logs" class="nav-item" :class="{ active: isActive('/admin/audit-logs') }">
            <div class="nav-icon-wrapper">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14 2 14 8 20 8"/>
                <line x1="16" y1="13" x2="8" y2="13"/>
                <line x1="16" y1="17" x2="8" y2="17"/>
                <polyline points="10 9 9 9 8 9"/>
              </svg>
            </div>
            <span class="nav-label">审计日志</span>
          </router-link>
        </div>
      </nav>

      <!-- Sidebar Footer -->
      <div class="sidebar-footer">
        <router-link to="/projects" class="back-link">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="19" y1="12" x2="5" y2="12"/>
            <polyline points="12 19 5 12 12 5"/>
          </svg>
          返回项目
        </router-link>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="admin-main">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()

function isActive(path: string): boolean {
  return route.path.startsWith(path)
}
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  background: var(--color-bg-base);
}

/* Sidebar */
.admin-sidebar {
  width: 260px;
  background: var(--color-bg-sidebar);
  color: #fff;
  display: flex;
  flex-direction: column;
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  z-index: 10;
}

.admin-sidebar::before {
  content: '';
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 1px;
  background: linear-gradient(
    180deg,
    transparent 0%,
    rgba(255, 255, 255, 0.1) 50%,
    transparent 100%
  );
}

/* Sidebar Header */
.sidebar-header {
  padding: var(--space-6);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.brand-link {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  text-decoration: none;
  transition: opacity var(--transition-base);
}

.brand-link:hover {
  opacity: 0.9;
}

.brand-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.brand-name {
  font-family: var(--font-display);
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-bold);
  color: #fff;
  letter-spacing: -0.01em;
}

.brand-badge {
  font-size: var(--font-size-xs);
  color: var(--color-primary-light);
  font-weight: var(--font-weight-medium);
}

/* Navigation */
.sidebar-nav {
  flex: 1;
  padding: var(--space-4) var(--space-3);
  overflow-y: auto;
}

.nav-section {
  margin-bottom: var(--space-6);
}

.nav-section-title {
  display: block;
  padding: var(--space-2) var(--space-3);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
  color: rgba(255, 255, 255, 0.4);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-3);
  color: rgba(255, 255, 255, 0.7);
  text-decoration: none;
  border-radius: var(--radius-lg);
  margin-bottom: var(--space-1);
  transition: all var(--transition-base);
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.06);
  color: #fff;
}

.nav-item.active {
  background: rgba(91, 95, 199, 0.3);
  color: #fff;
}

.nav-item.active .nav-icon-wrapper {
  background: var(--color-primary);
  color: #fff;
}

.nav-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: rgba(255, 255, 255, 0.08);
  border-radius: var(--radius-md);
  color: rgba(255, 255, 255, 0.7);
  transition: all var(--transition-base);
}

.nav-item:hover .nav-icon-wrapper {
  background: rgba(255, 255, 255, 0.12);
  color: #fff;
}

.nav-label {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
}

/* Sidebar Footer */
.sidebar-footer {
  padding: var(--space-4) var(--space-3);
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.back-link {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3);
  color: rgba(255, 255, 255, 0.6);
  text-decoration: none;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  border-radius: var(--radius-md);
  transition: all var(--transition-base);
}

.back-link:hover {
  background: rgba(255, 255, 255, 0.06);
  color: #fff;
}

/* Main Content */
.admin-main {
  flex: 1;
  margin-left: 260px;
  padding: var(--space-8);
  min-height: 100vh;
  background: var(--color-bg-base);
}

/* Responsive */
@media (max-width: 1024px) {
  .admin-sidebar {
    width: 72px;
  }

  .brand-text,
  .nav-section-title,
  .nav-label,
  .back-link span {
    display: none;
  }

  .brand-link {
    justify-content: center;
  }

  .nav-item {
    justify-content: center;
    padding: var(--space-3);
  }

  .nav-icon-wrapper {
    width: 40px;
    height: 40px;
  }

  .back-link {
    justify-content: center;
  }

  .admin-main {
    margin-left: 72px;
  }
}

@media (max-width: 768px) {
  .admin-sidebar {
    display: none;
  }

  .admin-main {
    margin-left: 0;
    padding: var(--space-4);
  }
}
</style>

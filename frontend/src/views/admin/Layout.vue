<template>
  <div class="admin-layout">
    <aside class="admin-sidebar">
      <div class="sidebar-header">
        <h2>管理后台</h2>
      </div>
      <nav class="sidebar-nav">
        <router-link to="/admin/users" class="nav-item">
          <span class="nav-icon">👥</span>
          <span>用户管理</span>
        </router-link>
        <router-link to="/admin/roles" class="nav-item">
          <span class="nav-icon">🔐</span>
          <span>角色管理</span>
        </router-link>
        <router-link to="/admin/audit-logs" class="nav-item">
          <span class="nav-icon">📋</span>
          <span>审计日志</span>
        </router-link>
      </nav>
      <div class="sidebar-footer">
        <router-link to="/projects" class="back-link">
          ← 返回项目
        </router-link>
      </div>
    </aside>
    <main class="admin-main">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
// Admin layout
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
}

.admin-sidebar {
  width: 240px;
  background: #1a1a2e;
  color: #fff;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #333;
}

.sidebar-header h2 {
  margin: 0;
  font-size: 18px;
}

.sidebar-nav {
  flex: 1;
  padding: 10px 0;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  color: #aaa;
  text-decoration: none;
  transition: background 0.2s;
}

.nav-item:hover,
.nav-item.router-link-active {
  background: #16213e;
  color: #fff;
}

.nav-icon {
  margin-right: 10px;
  font-size: 18px;
}

.sidebar-footer {
  padding: 20px;
  border-top: 1px solid #333;
}

.back-link {
  color: #aaa;
  text-decoration: none;
  font-size: 14px;
}

.back-link:hover {
  color: #fff;
}

.admin-main {
  flex: 1;
  padding: 24px;
  background: #f5f5f5;
  overflow-y: auto;
}
</style>

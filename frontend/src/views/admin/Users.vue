<template>
  <div class="admin-page">
    <div class="page-header animate-in">
      <div class="header-content">
        <div class="page-title-section">
          <h1>用户管理</h1>
          <p>管理系统用户和权限</p>
        </div>
        <div class="header-actions">
          <button @click="showImportDialog = true" class="btn btn-primary">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
              <polyline points="17 8 12 3 7 8"/>
              <line x1="12" y1="3" x2="12" y2="15"/>
            </svg>
            批量导入
          </button>
        </div>
      </div>

      <!-- Filter Bar -->
      <div class="filter-bar">
        <div class="filter-item">
          <input
            v-model="filters.email"
            type="text"
            placeholder="搜索邮箱..."
            class="filter-input"
            @keyup.enter="searchUsers"
          />
        </div>
        <div class="filter-item">
          <input
            v-model="filters.display_name"
            type="text"
            placeholder="搜索名称..."
            class="filter-input"
            @keyup.enter="searchUsers"
          />
        </div>
        <div class="filter-item">
          <select v-model="filters.is_disabled" class="filter-select">
            <option value="">全部状态</option>
            <option value="false">已启用</option>
            <option value="true">已禁用</option>
          </select>
        </div>
        <button @click="searchUsers" class="btn btn-primary btn-sm">搜索</button>
        <button @click="resetFilters" class="btn btn-ghost btn-sm">重置</button>
      </div>
    </div>

    <div v-if="loading" class="loading-state animate-in">
      <div class="loading-spinner"></div>
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

    <div v-else class="users-table-container animate-in">
      <table class="users-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>邮箱</th>
            <th>显示名称</th>
            <th>系统角色</th>
            <th>状态</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id" class="user-row">
            <td class="id-cell">{{ user.id }}</td>
            <td class="email-cell">{{ user.email }}</td>
            <td class="name-cell">{{ user.display_name || '-' }}</td>
            <td class="roles-cell">
              <div class="role-tags">
                <span v-for="role in getUserSystemRoles(user.id)" :key="role" class="role-tag">
                  {{ role }}
                </span>
                <span v-if="getUserSystemRoles(user.id).length === 0" class="no-role">无</span>
              </div>
            </td>
            <td class="status-cell">
              <span :class="['status-badge', user.is_disabled ? 'disabled' : 'active']">
                {{ user.is_disabled ? '已禁用' : '正常' }}
              </span>
            </td>
            <td class="actions-cell">
              <button @click="openRoleDialog(user)" class="btn btn-ghost btn-sm">分配角色</button>
              <button @click="handleResetPassword(user)" class="btn btn-warning btn-sm">重置密码</button>
              <button v-if="!user.is_disabled" @click="handleDisableUser(user)" class="btn btn-danger btn-sm">禁用</button>
              <button v-else @click="handleEnableUser(user)" class="btn btn-success btn-sm">启用</button>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Pagination -->
      <div class="pagination" v-if="totalPages > 0">
        <span class="page-info">共 {{ totalUsers }} 条</span>
        <div class="pagination-controls">
          <button @click="goToPage(currentPage - 1)" :disabled="currentPage <= 1" class="btn btn-ghost btn-sm">
            上一页
          </button>
          <span class="page-number">{{ currentPage }} / {{ totalPages }}</span>
          <button @click="goToPage(currentPage + 1)" :disabled="currentPage >= totalPages" class="btn btn-ghost btn-sm">
            下一页
          </button>
        </div>
      </div>
    </div>

    <!-- Role Dialog -->
    <div v-if="showRoleDialog" class="dialog-overlay" @click.self="closeRoleDialog">
      <div class="dialog animate-in">
        <div class="dialog-header">
          <h3>为 {{ selectedUser?.email }} 分配系统角色</h3>
          <button @click="closeRoleDialog" class="dialog-close">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="dialog-body">
          <label v-for="role in availableRoles" :key="role.id" class="role-option">
            <input type="radio" :value="role.id" v-model="selectedRoleId" />
            <div class="role-info">
              <span class="role-name">{{ role.display_name || role.name }}</span>
              <span class="role-desc">{{ role.description }}</span>
            </div>
          </label>
        </div>
        <div class="dialog-footer">
          <button @click="closeRoleDialog" class="btn btn-ghost">取消</button>
          <button @click="assignRole" class="btn btn-primary" :disabled="!selectedRoleId">确认</button>
        </div>
      </div>
    </div>

    <!-- Toast Notification -->
    <Transition name="toast">
      <div v-if="showToast" :class="['toast', toastType]">
        <svg v-if="toastType === 'success'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="20 6 9 17 4 12"/>
        </svg>
        <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <line x1="12" y1="8" x2="12" y2="12"/>
          <line x1="12" y1="16" x2="12.01" y2="16"/>
        </svg>
        {{ toastMessage }}
      </div>
    </Transition>

    <!-- Import Dialog -->
    <ImportUsersDialog
      :visible="showImportDialog"
      @close="showImportDialog = false"
      @success="loadData"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getUsers as getUsersApi, resetUserPassword, disableUser, enableUser } from '../../api/user'
import { getRoles } from '../../api/role'
import { assignSystemRole as assignSystemRoleApi, getUserRoles as getUserRolesApi } from '../../api/membership'
import type { ApiResponse, User, UserListResponse } from '../../api/user'
import ImportUsersDialog from '../../components/admin/ImportUsersDialog.vue'

interface Role {
  id: number
  name: string
  description: string
  is_system: boolean
  scope: string
  display_name?: string
}

const users = ref<User[]>([])
const roles = ref<Role[]>([])
const userRolesMap = ref<Map<number, string[]>>(new Map())
const loading = ref(false)
const error = ref<string | null>(null)
const showRoleDialog = ref(false)
const selectedUser = ref<User | null>(null)
const selectedRoleId = ref<number | null>(null)
const availableRoles = ref<Role[]>([])
const showImportDialog = ref(false)

const currentPage = ref(1)
const pageSize = ref(20)
const totalUsers = ref(0)
const totalPages = ref(0)

const filters = ref({
  email: '',
  display_name: '',
  is_disabled: ''
})

const showToast = ref(false)
const toastMessage = ref('')
const toastType = ref<'success' | 'error'>('success')

onMounted(async () => {
  await loadData()
})

async function loadData() {
  loading.value = true
  error.value = null
  try {
    const [usersRes, rolesRes] = await Promise.all([
      getUsersApi({
        email: filters.value.email || undefined,
        display_name: filters.value.display_name || undefined,
        is_disabled: filters.value.is_disabled || undefined,
        page: currentPage.value,
        page_size: pageSize.value
      }) as Promise<ApiResponse<UserListResponse>>,
      getRoles() as Promise<ApiResponse<{ roles: Role[] }>>
    ])

    if (usersRes.code === 0) {
      users.value = usersRes.data.users
      totalUsers.value = usersRes.data.total
      totalPages.value = usersRes.data.total_pages
    } else {
      error.value = usersRes.message
    }

    if (rolesRes.code === 0) {
      roles.value = rolesRes.data.roles.filter((r: Role) => r.scope === 'system')
      availableRoles.value = roles.value
    }

    for (const user of users.value) {
      const res = await getUserRolesApi(user.id) as any
      if (res.code === 0) {
        userRolesMap.value.set(user.id, res.data.system_roles || [])
      }
    }
  } catch (e: any) {
    error.value = e.response?.data?.message || '加载失败'
  } finally {
    loading.value = false
  }
}

function searchUsers() {
  currentPage.value = 1
  loadData()
}

function resetFilters() {
  filters.value = { email: '', display_name: '', is_disabled: '' }
  searchUsers()
}

function goToPage(page: number) {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  loadData()
}

function getUserSystemRoles(userId: number): string[] {
  return userRolesMap.value.get(userId) || []
}

function openRoleDialog(user: User) {
  selectedUser.value = user
  selectedRoleId.value = null
  showRoleDialog.value = true
}

function closeRoleDialog() {
  showRoleDialog.value = false
  selectedUser.value = null
  selectedRoleId.value = null
}

async function assignRole() {
  if (!selectedUser.value || !selectedRoleId.value) return
  try {
    const res = await assignSystemRoleApi(selectedUser.value.id, selectedRoleId.value) as ApiResponse<null>
    if (res.code === 0) {
      await loadData()
      closeRoleDialog()
      showNotification('角色分配成功', 'success')
    } else {
      showNotification(res.message, 'error')
    }
  } catch (e: any) {
    showNotification(e.response?.data?.message || '分配失败', 'error')
  }
}

async function handleResetPassword(user: User) {
  if (!confirm(`确定要重置用户 ${user.email} 的密码吗？`)) return
  try {
    const res = await resetUserPassword(user.id) as ApiResponse<null>
    if (res.code === 0) {
      showNotification(`密码已重置为: ${res.message.replace('密码已重置为: ', '')}`, 'success')
    } else {
      showNotification(res.message, 'error')
    }
  } catch (e: any) {
    showNotification(e.response?.data?.message || '重置密码失败', 'error')
  }
}

async function handleDisableUser(user: User) {
  if (!confirm(`确定要禁用用户 ${user.email} 吗？`)) return
  try {
    const res = await disableUser(user.id) as ApiResponse<null>
    if (res.code === 0) {
      showNotification('用户已禁用', 'success')
      await loadData()
    } else {
      showNotification(res.message, 'error')
    }
  } catch (e: any) {
    showNotification(e.response?.data?.message || '禁用用户失败', 'error')
  }
}

async function handleEnableUser(user: User) {
  if (!confirm(`确定要启用用户 ${user.email} 吗？`)) return
  try {
    const res = await enableUser(user.id) as ApiResponse<null>
    if (res.code === 0) {
      showNotification('用户已启用', 'success')
      await loadData()
    } else {
      showNotification(res.message, 'error')
    }
  } catch (e: any) {
    showNotification(e.response?.data?.message || '启用用户失败', 'error')
  }
}

function showNotification(message: string, type: 'success' | 'error') {
  toastMessage.value = message
  toastType.value = type
  showToast.value = true
  setTimeout(() => { showToast.value = false }, 3000)
}
</script>

<style scoped>
.admin-page {
  padding: var(--space-6);
  max-width: 1400px;
  margin: 0 auto;
}

/* Page Header */
.page-header {
  margin-bottom: var(--space-6);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--space-5);
}

.page-title-section h1 {
  font-family: var(--font-display);
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-text-primary);
  margin: 0 0 var(--space-1);
}

.page-title-section p {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  margin: 0;
}

/* Filter Bar */
.filter-bar {
  display: flex;
  gap: var(--space-3);
  flex-wrap: wrap;
  background: var(--color-bg-surface);
  padding: var(--space-4);
  border-radius: var(--radius-xl);
  border: 1px solid var(--color-border);
}

.filter-item {
  flex: 0 0 auto;
}

.filter-input,
.filter-select {
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--font-size-sm);
  min-width: 150px;
  transition: all var(--transition-base);
}

.filter-input:focus,
.filter-select:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(91, 95, 199, 0.1);
}

/* Buttons */
.btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  font-family: var(--font-body);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  transition: all var(--transition-base);
}

.btn-sm {
  padding: var(--space-1) var(--space-3);
  font-size: var(--font-size-xs);
}

.btn-primary {
  background: var(--color-primary);
  color: white;
}

.btn-primary:hover {
  background: var(--color-primary-hover);
}

.btn-ghost {
  background: transparent;
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
}

.btn-ghost:hover {
  background: var(--color-bg-muted);
  color: var(--color-text-primary);
}

.btn-warning {
  background: #f5a623;
  color: white;
}

.btn-danger {
  background: var(--color-danger);
  color: white;
}

.btn-success {
  background: var(--color-success);
  color: white;
}

/* Loading/Error */
.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-16);
  gap: var(--space-4);
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

/* Users Table */
.users-table-container {
  background: var(--color-bg-surface);
  border-radius: var(--radius-xl);
  border: 1px solid var(--color-border);
  overflow: hidden;
}

.users-table {
  width: 100%;
  border-collapse: collapse;
}

.users-table th,
.users-table td {
  padding: var(--space-4);
  text-align: left;
  border-bottom: 1px solid var(--color-border);
}

.users-table th {
  background: var(--color-bg-muted);
  font-weight: var(--font-weight-semibold);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.user-row:hover {
  background: var(--color-bg-muted);
}

.id-cell {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  width: 60px;
}

.email-cell {
  font-weight: var(--font-weight-medium);
}

.name-cell {
  color: var(--color-text-secondary);
}

.role-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-1);
}

.role-tag {
  display: inline-block;
  padding: 2px var(--space-2);
  background: rgba(91, 95, 199, 0.1);
  color: var(--color-primary);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
}

.no-role {
  color: var(--color-text-muted);
  font-size: var(--font-size-xs);
}

.status-badge {
  display: inline-block;
  padding: 2px var(--space-2);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
}

.status-badge.active {
  background: rgba(78, 205, 196, 0.15);
  color: #2e8b7d;
}

.status-badge.disabled {
  background: rgba(239, 71, 111, 0.1);
  color: var(--color-danger);
}

.actions-cell {
  white-space: nowrap;
}

/* Pagination */
.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-4);
  border-top: 1px solid var(--color-border);
}

.page-info {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.page-number {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

/* Dialog */
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--z-modal);
  backdrop-filter: blur(4px);
}

.dialog {
  background: var(--color-bg-surface);
  border-radius: var(--radius-xl);
  padding: var(--space-6);
  min-width: 400px;
  max-width: 90%;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: var(--shadow-xl);
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-5);
}

.dialog-header h3 {
  font-family: var(--font-display);
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  margin: 0;
}

.dialog-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: none;
  border: none;
  border-radius: var(--radius-md);
  color: var(--color-text-muted);
  cursor: pointer;
  transition: all var(--transition-base);
}

.dialog-close:hover {
  background: var(--color-bg-muted);
  color: var(--color-text-primary);
}

.dialog-body {
  margin-bottom: var(--space-5);
}

.role-option {
  display: flex;
  align-items: flex-start;
  gap: var(--space-3);
  padding: var(--space-3);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--transition-base);
}

.role-option:hover {
  background: var(--color-bg-muted);
}

.role-option input {
  margin-top: var(--space-1);
}

.role-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.role-name {
  font-weight: var(--font-weight-medium);
  color: var(--color-text-primary);
}

.role-desc {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  padding-top: var(--space-4);
  border-top: 1px solid var(--color-border);
}

/* Toast */
.toast {
  position: fixed;
  bottom: var(--space-6);
  right: var(--space-6);
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-5);
  border-radius: var(--radius-lg);
  color: white;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  z-index: var(--z-toast);
  box-shadow: var(--shadow-lg);
}

.toast.success {
  background: var(--color-success);
}

.toast.error {
  background: var(--color-danger);
}

.toast-enter-active,
.toast-leave-active {
  transition: all var(--transition-base);
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>

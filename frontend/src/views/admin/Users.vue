<template>
  <div class="users-page">
    <div class="page-header">
      <h1>用户管理</h1>
      <button @click="showImportDialog = true" class="btn btn-primary">批量导入</button>
    </div>

    <!-- 搜索筛选区域 -->
    <div class="filter-bar">
      <div class="filter-item">
        <input
          v-model="filters.email"
          type="text"
          placeholder="按邮箱搜索"
          class="filter-input"
        />
      </div>
      <div class="filter-item">
        <input
          v-model="filters.display_name"
          type="text"
          placeholder="按显示名称搜索"
          class="filter-input"
        />
      </div>
      <div class="filter-item">
        <select v-model="filters.is_disabled" class="filter-select">
          <option value="">全部状态</option>
          <option value="false">已启用</option>
          <option value="true">已禁用</option>
        </select>
      </div>
      <button @click="searchUsers" class="btn btn-primary">搜索</button>
      <button @click="resetFilters" class="btn btn-secondary">重置</button>
    </div>

    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>

    <div v-else class="users-table-container">
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
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.id }}</td>
            <td>{{ user.email }}</td>
            <td>{{ user.display_name }}</td>
            <td>
              <span v-for="role in getUserSystemRoles(user.id)" :key="role" class="role-tag">
                {{ role }}
              </span>
              <span v-if="getUserSystemRoles(user.id).length === 0" class="no-role">无</span>
            </td>
            <td>
              <span :class="['status-tag', user.is_disabled ? 'status-disabled' : 'status-active']">
                {{ user.is_disabled ? '已禁用' : '正常' }}
              </span>
            </td>
            <td class="action-cell">
              <button @click="openRoleDialog(user)" class="btn btn-sm">分配角色</button>
              <button @click="handleResetPassword(user)" class="btn btn-sm btn-warning">重置密码</button>
              <button v-if="!user.is_disabled" @click="handleDisableUser(user)" class="btn btn-sm btn-danger">禁用</button>
              <button v-else @click="handleEnableUser(user)" class="btn btn-sm btn-success">启用</button>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- 分页控件 -->
      <div class="pagination" v-if="totalPages > 0">
        <button
          @click="goToPage(currentPage - 1)"
          :disabled="currentPage <= 1"
          class="btn btn-sm"
        >
          上一页
        </button>
        <span class="page-info">第 {{ currentPage }} / {{ totalPages }} 页，共 {{ totalUsers }} 条</span>
        <button
          @click="goToPage(currentPage + 1)"
          :disabled="currentPage >= totalPages"
          class="btn btn-sm"
        >
          下一页
        </button>
      </div>
    </div>

    <!-- 分配角色对话框 -->
    <div v-if="showRoleDialog" class="dialog-overlay" @click.self="closeRoleDialog">
      <div class="dialog">
        <h3>为 {{ selectedUser?.email }} 分配系统角色</h3>
        <div class="role-list">
          <label v-for="role in availableRoles" :key="role.id" class="role-option">
            <input type="radio" :value="role.id" v-model="selectedRoleId" />
            <span>{{ role.display_name }}</span>
          </label>
        </div>
        <div class="dialog-actions">
          <button @click="closeRoleDialog" class="btn btn-secondary">取消</button>
          <button @click="assignRole" class="btn btn-primary" :disabled="!selectedRoleId">确认</button>
        </div>
      </div>
    </div>

    <!-- 操作结果提示 -->
    <div v-if="showToast" :class="['toast', toastType]">
      {{ toastMessage }}
    </div>

    <!-- 批量导入对话框 -->
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

// 分页相关
const currentPage = ref(1)
const pageSize = ref(20)
const totalUsers = ref(0)
const totalPages = ref(0)

// 搜索筛选
const filters = ref({
  email: '',
  display_name: '',
  is_disabled: ''
})

// 提示消息
const showToast = ref(false)
const toastMessage = ref('')
const toastType = ref('success')

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
      getRoles() as Promise<ApiResponse<Role[]>>
    ])

    if (usersRes.code === 0) {
      users.value = usersRes.data.users
      totalUsers.value = usersRes.data.total
      totalPages.value = usersRes.data.total_pages
    } else {
      error.value = usersRes.message
    }

    if (rolesRes.code === 0) {
      roles.value = rolesRes.data.filter((r: Role) => r.scope === 'system')
      availableRoles.value = roles.value
    }

    // 加载每个用户的角色
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
  filters.value = {
    email: '',
    display_name: '',
    is_disabled: ''
  }
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
  setTimeout(() => {
    showToast.value = false
  }, 3000)
}
</script>

<style scoped>
.users-page {
  max-width: 1400px;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h1 {
  margin: 0;
  font-size: 24px;
  color: #333;
}

.filter-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.filter-item {
  flex: 0 0 auto;
}

.filter-input,
.filter-select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  min-width: 150px;
}

.filter-input:focus,
.filter-select:focus {
  outline: none;
  border-color: #1976d2;
}

.loading,
.error {
  padding: 40px;
  text-align: center;
}

.error {
  color: #dc3545;
}

.users-table-container {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  overflow: hidden;
}

.users-table {
  width: 100%;
  border-collapse: collapse;
}

.users-table th,
.users-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.users-table th {
  background: #f8f9fa;
  font-weight: 600;
  color: #333;
}

.users-table tbody tr:hover {
  background: #f5f5f5;
}

.role-tag {
  display: inline-block;
  padding: 2px 8px;
  margin-right: 4px;
  background: #e3f2fd;
  color: #1976d2;
  border-radius: 4px;
  font-size: 12px;
}

.no-role {
  color: #999;
  font-size: 12px;
}

.status-tag {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.status-active {
  background: #e8f5e9;
  color: #2e7d32;
}

.status-disabled {
  background: #ffebee;
  color: #c62828;
}

.action-cell {
  white-space: nowrap;
}

.btn {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-sm {
  padding: 4px 8px;
  font-size: 12px;
  margin-right: 4px;
}

.btn-primary {
  background: #1976d2;
  color: #fff;
}

.btn-secondary {
  background: #6c757d;
  color: #fff;
}

.btn-warning {
  background: #ff9800;
  color: #fff;
}

.btn-danger {
  background: #dc3545;
  color: #fff;
}

.btn-success {
  background: #28a745;
  color: #fff;
}

.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog {
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  min-width: 320px;
  max-width: 400px;
}

.dialog h3 {
  margin: 0 0 16px;
  font-size: 18px;
}

.role-list {
  margin: 16px 0;
}

.role-option {
  display: flex;
  align-items: center;
  padding: 8px;
  cursor: pointer;
}

.role-option:hover {
  background: #f5f5f5;
}

.role-option input {
  margin-right: 8px;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 16px;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 16px;
  border-top: 1px solid #eee;
}

.page-info {
  color: #666;
  font-size: 14px;
}

.toast {
  position: fixed;
  bottom: 20px;
  right: 20px;
  padding: 12px 24px;
  border-radius: 4px;
  color: #fff;
  font-size: 14px;
  z-index: 2000;
  animation: fadeIn 0.3s ease;
}

.toast.success {
  background: #28a745;
}

.toast.error {
  background: #dc3545;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>

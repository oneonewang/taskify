<template>
  <div class="users-page">
    <div class="page-header">
      <h1>用户管理</h1>
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
              <button @click="openRoleDialog(user)" class="btn btn-sm">分配角色</button>
            </td>
          </tr>
        </tbody>
      </table>
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getUsers as getUsersApi } from '../../api/user'
import { getRoles } from '../../api/role'
import { assignSystemRole as assignSystemRoleApi, getUserRoles as getUserRolesApi } from '../../api/membership'
import type { ApiResponse } from '../../api/auth'
import type { User } from '../../api/auth'

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

onMounted(async () => {
  await loadData()
})

async function loadData() {
  loading.value = true
  error.value = null
  try {
    const [usersRes, rolesRes] = await Promise.all([
      getUsersApi() as Promise<ApiResponse<User[]>>,
      getRoles() as Promise<ApiResponse<Role[]>>
    ])

    if (usersRes.code === 0) {
      users.value = usersRes.data
    } else {
      error.value = usersRes.message
    }

    if (rolesRes.code === 0) {
      roles.value = rolesRes.data.filter(r => r.scope === 'system')
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
    } else {
      error.value = res.message
    }
  } catch (e: any) {
    error.value = e.response?.data?.message || '分配失败'
  }
}
</script>

<style scoped>
.users-page {
  max-width: 1200px;
}

.page-header {
  margin-bottom: 24px;
}

.page-header h1 {
  margin: 0;
  font-size: 24px;
  color: #333;
}

.loading, .error {
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
  background: #f8f9fa;
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
}

.btn-primary {
  background: #1976d2;
  color: #fff;
}

.btn-secondary {
  background: #6c757d;
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
</style>

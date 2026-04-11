<template>
  <div class="roles-page">
    <div class="page-header">
      <h1>角色管理</h1>
      <button @click="openCreateDialog" class="btn btn-primary">创建角色</button>
    </div>

    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>

    <div v-else class="roles-grid">
      <div v-for="role in roles" :key="role.id" class="role-card">
        <div class="role-header">
          <h3>{{ role.display_name || role.name }}</h3>
          <span v-if="role.is_system" class="system-badge">系统</span>
        </div>
        <p class="role-description">{{ role.description || '无描述' }}</p>
        <div class="role-meta">
          <span class="scope-tag">{{ role.scope === 'system' ? '系统级' : '项目级' }}</span>
        </div>
        <div class="role-actions">
          <button @click="openEditDialog(role)" class="btn btn-sm" :disabled="role.is_system">编辑</button>
          <button @click="openPermissionDialog(role)" class="btn btn-sm">权限</button>
          <button v-if="!role.is_system" @click="handleDeleteRole(role)" class="btn btn-sm btn-danger">删除</button>
        </div>
      </div>
    </div>

    <!-- 创建/编辑角色对话框 -->
    <div v-if="showEditDialog" class="dialog-overlay" @click.self="closeEditDialog">
      <div class="dialog">
        <h3>{{ editingRole?.id ? '编辑角色' : '创建角色' }}</h3>
        <div v-if="dialogError" class="dialog-error">{{ dialogError }}</div>
        <div class="form-group">
          <label>角色名称</label>
          <input v-model="editForm.name" type="text" :disabled="!!editingRole?.id" />
        </div>
        <div class="form-group">
          <label>描述</label>
          <textarea v-model="editForm.description" rows="3"></textarea>
        </div>
        <div class="dialog-actions">
          <button @click="closeEditDialog" class="btn btn-secondary">取消</button>
          <button @click="saveRole" class="btn btn-primary">{{ editingRole?.id ? '保存' : '创建' }}</button>
        </div>
      </div>
    </div>

    <!-- 权限配置对话框 -->
    <div v-if="showPermissionDialog" class="dialog-overlay" @click.self="closePermissionDialog">
      <div class="dialog dialog-wide">
        <h3>配置 {{ selectedRole?.name }} 的权限</h3>
        <div class="permissions-grid">
          <div v-for="group in permissionGroups" :key="group.resource" class="permission-group">
            <h4>{{ group.resource }}</h4>
            <label v-for="perm in group.permissions" :key="perm.id" class="permission-option">
              <input type="checkbox" :checked="hasPermission(perm.id)" @change="togglePermission(perm.id)" />
              <span>{{ perm.action }}</span>
            </label>
          </div>
        </div>
        <div class="dialog-actions">
          <button @click="handleSavePermissions" class="btn btn-primary">保存</button>
          <button @click="closePermissionDialog" class="btn btn-secondary">关闭</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { getRoles, createRole, updateRole, deleteRole as deleteRoleApi, setRolePermissions, getPermissions } from '../../api/role'
import type { ApiResponse } from '../../api/role'

interface Role {
  id: number
  name: string
  description: string
  is_system: boolean
  scope: string
  display_name?: string
  permissions?: string[]
}

interface Permission {
  id: number
  resource: string
  action: string
  description: string
}

const roles = ref<Role[]>([])
const permissions = ref<Permission[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const dialogError = ref<string | null>(null)
const showEditDialog = ref(false)
const showPermissionDialog = ref(false)
const editingRole = ref<Role | null>(null)
const selectedRole = ref<Role | null>(null)
const selectedPermissions = ref<Set<number>>(new Set())

const editForm = ref({
  name: '',
  description: ''
})

onMounted(async () => {
  await loadData()
})

async function loadData() {
  loading.value = true
  error.value = null
  try {
    const [rolesRes, permsRes] = await Promise.all([
      getRoles() as Promise<ApiResponse<{ roles: Role[] }>>,
      getPermissions() as Promise<ApiResponse<{ permissions: Permission[] }>>
    ])

    if (rolesRes.code === 0) {
      roles.value = rolesRes.data.roles
    } else {
      error.value = rolesRes.message
    }

    if (permsRes.code === 0) {
      permissions.value = permsRes.data.permissions
    }
  } catch (e: any) {
    error.value = e.response?.data?.message || '加载失败'
  } finally {
    loading.value = false
  }
}

const permissionGroups = computed(() => {
  const groups: { resource: string; permissions: Permission[] }[] = []
  const resourceMap = new Map<string, Permission[]>()

  for (const perm of permissions.value) {
    const list = resourceMap.get(perm.resource) || []
    list.push(perm)
    resourceMap.set(perm.resource, list)
  }

  for (const [resource, perms] of resourceMap) {
    groups.push({ resource, permissions: perms })
  }

  return groups
})

function hasPermission(permId: number): boolean {
  return selectedPermissions.value.has(permId)
}

function togglePermission(permId: number) {
  if (selectedPermissions.value.has(permId)) {
    selectedPermissions.value.delete(permId)
  } else {
    selectedPermissions.value.add(permId)
  }
}

function openCreateDialog() {
  dialogError.value = null
  editingRole.value = null
  editForm.value = { name: '', description: '' }
  showEditDialog.value = true
}

function openEditDialog(role: Role) {
  dialogError.value = null
  editingRole.value = role
  editForm.value = { name: role.name, description: role.description }
  showEditDialog.value = true
}

function closeEditDialog() {
  showEditDialog.value = false
  editingRole.value = null
  dialogError.value = null
  editForm.value = { name: '', description: '' }
}

async function saveRole() {
  dialogError.value = null
  try {
    if (editingRole.value?.id) {
      const res = await updateRole(editingRole.value.id, {
        description: editForm.value.description
      }) as ApiResponse<Role>
      if (res.code === 0) {
        await loadData()
        closeEditDialog()
      } else {
        dialogError.value = res.message
      }
    } else {
      const res = await createRole({
        name: editForm.value.name,
        description: editForm.value.description,
        scope: 'project'
      }) as ApiResponse<Role>
      if (res.code === 0) {
        await loadData()
        closeEditDialog()
      } else {
        dialogError.value = res.message
      }
    }
  } catch (e: any) {
    dialogError.value = e.response?.data?.message || '操作失败'
  }
}

async function handleDeleteRole(role: Role) {
  if (!confirm(`确定要删除角色 "${role.name}" 吗？`)) return

  try {
    const res = await deleteRoleApi(role.id) as ApiResponse<null>
    if (res.code === 0) {
      await loadData()
    } else {
      error.value = res.message
    }
  } catch (e: any) {
    error.value = e.response?.data?.message || '删除失败'
  }
}

function openPermissionDialog(role: Role) {
  selectedRole.value = role
  selectedPermissions.value = new Set(role.permissions?.map(Number) || [])
  showPermissionDialog.value = true
}

function closePermissionDialog() {
  showPermissionDialog.value = false
  selectedRole.value = null
  selectedPermissions.value = new Set()
}

async function handleSavePermissions() {
  if (!selectedRole.value) return

  try {
    const res = await setRolePermissions(selectedRole.value.id, Array.from(selectedPermissions.value)) as ApiResponse<null>
    if (res.code === 0) {
      await loadData()
      closePermissionDialog()
    } else {
      error.value = res.message
    }
  } catch (e: any) {
    error.value = e.response?.data?.message || '保存失败'
  }
}
</script>

<style scoped>
.roles-page {
  max-width: 1200px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.roles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.role-card {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.role-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.role-header h3 {
  margin: 0;
  font-size: 16px;
}

.system-badge {
  padding: 2px 6px;
  background: #e3f2fd;
  color: #1976d2;
  border-radius: 4px;
  font-size: 11px;
}

.role-description {
  margin: 0 0 12px;
  color: #666;
  font-size: 14px;
}

.role-meta {
  margin-bottom: 12px;
}

.scope-tag {
  padding: 2px 8px;
  background: #f5f5f5;
  color: #666;
  border-radius: 4px;
  font-size: 12px;
}

.role-actions {
  display: flex;
  gap: 8px;
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

.btn-danger {
  background: #dc3545;
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

.dialog-wide {
  max-width: 600px;
  max-height: 80vh;
  overflow-y: auto;
}

.dialog h3 {
  margin: 0 0 16px;
  font-size: 18px;
}

.dialog-error {
  background: #ffebee;
  color: #c62828;
  padding: 8px 12px;
  border-radius: 4px;
  margin-bottom: 16px;
  font-size: 14px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 4px;
  font-size: 14px;
  color: #666;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 16px;
}

.permissions-grid {
  max-height: 400px;
  overflow-y: auto;
}

.permission-group {
  margin-bottom: 16px;
}

.permission-group h4 {
  margin: 0 0 8px;
  font-size: 14px;
  color: #333;
  text-transform: capitalize;
}

.permission-option {
  display: flex;
  align-items: center;
  padding: 4px 0;
  cursor: pointer;
}

.permission-option input {
  margin-right: 8px;
}

.permission-option span {
  font-size: 13px;
  color: #666;
}
</style>

<template>
  <div class="admin-page">
    <div class="page-header animate-in">
      <div class="header-content">
        <div class="page-title-section">
          <h1>角色管理</h1>
          <p>管理系统角色和权限配置</p>
        </div>
        <div class="header-actions">
          <button @click="openCreateDialog" class="btn btn-primary">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"/>
              <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
            创建角色
          </button>
        </div>
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

    <div v-else class="roles-grid">
      <div v-for="(role, index) in roles" :key="role.id" class="role-card animate-in" :style="{ animationDelay: `${index * 0.05}s` }">
        <div class="role-header">
          <div class="role-icon">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
            </svg>
          </div>
          <h3>{{ role.display_name || role.name }}</h3>
          <span v-if="role.is_system" class="system-badge">系统</span>
        </div>
        <p class="role-description">{{ role.description || '无描述' }}</p>
        <div class="role-meta">
          <span class="scope-tag" :class="role.scope">
            {{ role.scope === 'system' ? '系统级' : '项目级' }}
          </span>
          <span class="permission-count">
            {{ role.permissions?.length || 0 }} 个权限
          </span>
        </div>
        <div class="role-actions">
          <button @click="openEditDialog(role)" class="btn btn-ghost btn-sm" :disabled="role.is_system">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
              <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
            </svg>
            编辑
          </button>
          <button @click="openPermissionDialog(role)" class="btn btn-primary btn-sm">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9 11 12 14 22 4"/>
              <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/>
            </svg>
            权限
          </button>
          <button v-if="!role.is_system" @click="handleDeleteRole(role)" class="btn btn-danger btn-sm">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="3 6 5 6 21 6"/>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
            </svg>
            删除
          </button>
        </div>
      </div>
    </div>

    <!-- Create/Edit Dialog -->
    <div v-if="showEditDialog" class="dialog-overlay" @click.self="closeEditDialog">
      <div class="dialog animate-in">
        <div class="dialog-header">
          <h3>{{ editingRole?.id ? '编辑角色' : '创建角色' }}</h3>
          <button @click="closeEditDialog" class="dialog-close">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div v-if="dialogError" class="dialog-error">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="12"/>
            <line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
          {{ dialogError }}
        </div>
        <div class="dialog-body">
          <div class="form-group">
            <label for="role-name">角色名称</label>
            <input
              id="role-name"
              v-model="editForm.name"
              type="text"
              :disabled="!!editingRole?.id"
              class="form-input"
              placeholder="输入角色名称"
            />
          </div>
          <div class="form-group">
            <label for="role-desc">描述</label>
            <textarea
              id="role-desc"
              v-model="editForm.description"
              rows="3"
              class="form-textarea"
              placeholder="描述角色的用途"
            ></textarea>
          </div>
        </div>
        <div class="dialog-footer">
          <button @click="closeEditDialog" class="btn btn-ghost">取消</button>
          <button @click="saveRole" class="btn btn-primary">{{ editingRole?.id ? '保存' : '创建' }}</button>
        </div>
      </div>
    </div>

    <!-- Permission Dialog -->
    <div v-if="showPermissionDialog" class="dialog-overlay" @click.self="closePermissionDialog">
      <div class="dialog dialog-wide animate-in">
        <div class="dialog-header">
          <h3>配置 {{ selectedRole?.name }} 的权限</h3>
          <button @click="closePermissionDialog" class="dialog-close">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="permissions-grid">
          <div v-for="group in permissionGroups" :key="group.resource" class="permission-group">
            <h4>{{ group.resource }}</h4>
            <label v-for="perm in group.permissions" :key="perm.id" class="permission-option">
              <input type="checkbox" :checked="hasPermission(perm.id)" @change="togglePermission(perm.id)" />
              <div class="perm-info">
                <span class="perm-action">{{ perm.action }}</span>
                <span class="perm-desc">{{ perm.description }}</span>
              </div>
            </label>
          </div>
        </div>
        <div class="dialog-footer">
          <button @click="handleSavePermissions" class="btn btn-primary">保存</button>
          <button @click="closePermissionDialog" class="btn btn-ghost">关闭</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
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
.admin-page {
  padding: var(--space-6);
  max-width: 1200px;
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

.btn-primary:hover:not(:disabled) {
  background: var(--color-primary-hover);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
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

.btn-danger {
  background: var(--color-danger);
  color: white;
}

.btn-danger:hover {
  background: #e04560;
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

/* Roles Grid */
.roles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: var(--space-5);
}

/* Role Card */
.role-card {
  background: var(--color-bg-surface);
  border-radius: var(--radius-xl);
  padding: var(--space-5);
  box-shadow: var(--shadow-card);
  border: 1px solid var(--color-border);
  transition: all var(--transition-base);
}

.role-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-card-hover);
}

.role-header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-3);
}

.role-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  background: rgba(91, 95, 199, 0.1);
  border-radius: var(--radius-lg);
  color: var(--color-primary);
}

.role-header h3 {
  flex: 1;
  font-family: var(--font-display);
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  margin: 0;
}

.system-badge {
  padding: 2px var(--space-2);
  background: rgba(91, 95, 199, 0.1);
  color: var(--color-primary);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
}

.role-description {
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  margin: 0 0 var(--space-4);
  line-height: var(--line-height-relaxed);
}

.role-meta {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-4);
}

.scope-tag {
  padding: 2px var(--space-2);
  background: var(--color-bg-muted);
  color: var(--color-text-secondary);
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
}

.scope-tag.system {
  background: rgba(91, 95, 199, 0.1);
  color: var(--color-primary);
}

.permission-count {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
}

.role-actions {
  display: flex;
  gap: var(--space-2);
  padding-top: var(--space-4);
  border-top: 1px solid var(--color-border);
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

.dialog-wide {
  min-width: 500px;
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

.dialog-error {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  background: rgba(239, 71, 111, 0.1);
  color: var(--color-danger);
  padding: var(--space-3);
  border-radius: var(--radius-md);
  margin-bottom: var(--space-4);
  font-size: var(--font-size-sm);
}

.dialog-body {
  margin-bottom: var(--space-5);
}

.form-group {
  margin-bottom: var(--space-4);
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-group label {
  display: block;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-secondary);
  margin-bottom: var(--space-2);
}

.form-input,
.form-textarea {
  width: 100%;
  padding: var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-family: var(--font-body);
  font-size: var(--font-size-base);
  transition: all var(--transition-base);
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(91, 95, 199, 0.1);
}

.form-input:disabled {
  background: var(--color-bg-muted);
  cursor: not-allowed;
}

.form-textarea {
  resize: vertical;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  padding-top: var(--space-4);
  border-top: 1px solid var(--color-border);
}

/* Permissions Grid */
.permissions-grid {
  max-height: 400px;
  overflow-y: auto;
  margin-bottom: var(--space-5);
}

.permission-group {
  margin-bottom: var(--space-5);
}

.permission-group:last-child {
  margin-bottom: 0;
}

.permission-group h4 {
  font-family: var(--font-display);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  margin: 0 0 var(--space-3);
  text-transform: capitalize;
}

.permission-option {
  display: flex;
  align-items: flex-start;
  gap: var(--space-3);
  padding: var(--space-2);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--transition-base);
}

.permission-option:hover {
  background: var(--color-bg-muted);
}

.permission-option input {
  margin-top: var(--space-1);
}

.perm-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.perm-action {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-primary);
}

.perm-desc {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
}
</style>

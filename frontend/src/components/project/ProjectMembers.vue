<template>
  <div class="project-members">
    <div class="members-header">
      <h4>项目成员</h4>
      <button v-if="canManage" @click="showAddDialog = true" class="btn btn-sm">
        添加成员
      </button>
    </div>

    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>

    <div v-else class="members-list">
      <div v-for="member in members" :key="member.user_id" class="member-item">
        <div class="member-avatar">
          {{ getInitials(member.display_name) }}
        </div>
        <div class="member-info">
          <div class="member-name">{{ member.display_name }}</div>
          <div class="member-email">{{ member.email }}</div>
        </div>
        <div class="member-role">
          <span class="role-badge">{{ member.role_display_name }}</span>
        </div>
        <div v-if="canManage" class="member-actions">
          <button @click="openEditDialog(member)" class="btn btn-xs">编辑</button>
          <button @click="removeMember(member)" class="btn btn-xs btn-danger">移除</button>
        </div>
      </div>

      <div v-if="members.length === 0" class="empty-state">
        暂无成员
      </div>
    </div>

    <!-- 添加成员对话框 -->
    <div v-if="showAddDialog" class="dialog-overlay" @click.self="closeAddDialog">
      <div class="dialog">
        <h3>添加项目成员</h3>
        <div class="form-group">
          <label>用户邮箱</label>
          <input v-model="addForm.email" type="email" placeholder="输入用户邮箱" />
        </div>
        <div class="form-group">
          <label>角色</label>
          <select v-model="addForm.roleId">
            <option value="">选择角色</option>
            <option v-for="role in projectRoles" :key="role.id" :value="role.id">
              {{ role.display_name }}
            </option>
          </select>
        </div>
        <div v-if="addError" class="error-message">{{ addError }}</div>
        <div class="dialog-actions">
          <button @click="closeAddDialog" class="btn btn-secondary">取消</button>
          <button @click="handleAddMember" class="btn btn-primary" :disabled="adding">
            {{ adding ? '添加中...' : '添加' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 编辑角色对话框 -->
    <div v-if="showEditDialog" class="dialog-overlay" @click.self="closeEditDialog">
      <div class="dialog">
        <h3>修改成员角色</h3>
        <div class="form-group">
          <label>角色</label>
          <select v-model="editForm.roleId">
            <option v-for="role in projectRoles" :key="role.id" :value="role.id">
              {{ role.display_name }}
            </option>
          </select>
        </div>
        <div class="dialog-actions">
          <button @click="closeEditDialog" class="btn btn-secondary">取消</button>
          <button @click="handleUpdateRole" class="btn btn-primary">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getProjectMembers, addProjectMember, updateMemberRole, removeProjectMember } from '../../api/membership'
import type { ApiResponse } from '../../api/membership'

interface Member {
  user_id: number
  email: string
  display_name: string
  avatar_url: string
  role_id: number
  role_name: string
  role_display_name: string
  joined_at: string
}

const props = defineProps<{
  projectId: number
  canManage: boolean
}>()

const members = ref<Member[]>([])
const projectRoles = ref<{ id: number; name: string; display_name: string }[]>([
  { id: 3, name: 'owner', display_name: '所有者' },
  { id: 4, name: 'member', display_name: '成员' },
  { id: 5, name: 'guest', display_name: '访客' }
])
const loading = ref(false)
const error = ref<string | null>(null)
const showAddDialog = ref(false)
const showEditDialog = ref(false)
const adding = ref(false)
const addError = ref<string | null>(null)

const addForm = ref({ email: '', roleId: '' as number | '' })
const editForm = ref({ userId: 0, roleId: '' as number | '' })
const editingMember = ref<Member | null>(null)

onMounted(async () => {
  await loadMembers()
})

async function loadMembers() {
  loading.value = true
  error.value = null
  try {
    const res = await getProjectMembers(props.projectId) as ApiResponse<{ members: Member[] }>
    if (res.code === 0) {
      members.value = res.data.members
    } else {
      error.value = res.message
    }
  } catch (e: any) {
    error.value = e.message || '加载失败'
  } finally {
    loading.value = false
  }
}

function getInitials(name: string): string {
  return name.slice(0, 2).toUpperCase()
}

function closeAddDialog() {
  showAddDialog.value = false
  addForm.value = { email: '', roleId: '' }
  addError.value = null
}

function openEditDialog(member: Member) {
  editingMember.value = member
  editForm.value = { userId: member.user_id, roleId: member.role_id }
  showEditDialog.value = true
}

function closeEditDialog() {
  showEditDialog.value = false
  editingMember.value = null
}

async function handleAddMember() {
  if (!addForm.value.email || !addForm.value.roleId) {
    addError.value = '请填写完整信息'
    return
  }

  adding.value = true
  addError.value = null
  try {
    const res = await addProjectMember(props.projectId, addForm.value.email, addForm.value.roleId as number) as ApiResponse<null>
    if (res.code === 0) {
      await loadMembers()
      closeAddDialog()
    } else {
      addError.value = res.message
    }
  } catch (e: any) {
    addError.value = e.message || '添加失败'
  } finally {
    adding.value = false
  }
}

async function handleUpdateRole() {
  if (!editingMember.value || !editForm.value.roleId) return

  try {
    const res = await updateMemberRole(props.projectId, editingMember.value.user_id, editForm.value.roleId as number) as ApiResponse<null>
    if (res.code === 0) {
      await loadMembers()
      closeEditDialog()
    }
  } catch (e: any) {
    console.error('更新角色失败', e)
  }
}

async function removeMember(member: Member) {
  if (!confirm(`确定要移除成员 "${member.display_name}" 吗？`)) return

  try {
    const res = await removeProjectMember(props.projectId, member.user_id) as ApiResponse<null>
    if (res.code === 0) {
      await loadMembers()
    }
  } catch (e: any) {
    console.error('移除成员失败', e)
  }
}
</script>

<style scoped>
.project-members {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
}

.members-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.members-header h4 {
  margin: 0;
  font-size: 14px;
  color: #333;
}

.members-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.member-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px;
  border-radius: 6px;
}

.member-item:hover {
  background: #f5f5f5;
}

.member-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #1976d2;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 500;
}

.member-info {
  flex: 1;
  min-width: 0;
}

.member-name {
  font-size: 14px;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.member-email {
  font-size: 12px;
  color: #999;
}

.role-badge {
  padding: 2px 8px;
  background: #e3f2fd;
  color: #1976d2;
  border-radius: 4px;
  font-size: 12px;
}

.member-actions {
  display: flex;
  gap: 4px;
}

.loading, .error, .empty-state {
  padding: 20px;
  text-align: center;
  color: #666;
  font-size: 14px;
}

.error {
  color: #dc3545;
}

.btn {
  padding: 4px 8px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.btn-sm {
  padding: 6px 12px;
  font-size: 13px;
}

.btn-xs {
  padding: 2px 6px;
  font-size: 11px;
}

.btn-primary {
  background: #1976d2;
  color: #fff;
}

.btn-secondary {
  background: #f5f5f5;
  color: #666;
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
}

.dialog h3 {
  margin: 0 0 16px;
  font-size: 16px;
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
.form-group select {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}

.error-message {
  padding: 8px;
  background: #fef2f2;
  color: #dc3545;
  border-radius: 4px;
  font-size: 13px;
  margin-bottom: 12px;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 16px;
}
</style>

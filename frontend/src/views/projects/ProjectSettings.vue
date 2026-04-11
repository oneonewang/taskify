<template>
  <div class="project-settings-page">
    <AppHeader title="项目设置" :show-back="true" @back="goBack" />

    <main class="settings-content">
      <div v-if="loading" class="loading-state animate-in">
        <div class="loading-spinner"></div>
        <span>加载中...</span>
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

      <template v-else>
        <!-- Basic Info Section -->
        <section class="settings-section animate-in">
          <div class="section-header">
            <h2>基本信息</h2>
            <p>项目名称和描述</p>
          </div>
          <div class="section-body">
            <div class="form-group">
              <label for="project-name">项目名称</label>
              <input
                id="project-name"
                v-model="form.name"
                type="text"
                maxlength="100"
                class="form-input"
                placeholder="输入项目名称"
              />
            </div>
            <div class="form-group">
              <label for="project-description">描述</label>
              <textarea
                id="project-description"
                v-model="form.description"
                rows="3"
                maxlength="500"
                class="form-textarea"
                placeholder="描述项目的目的和功能"
              ></textarea>
            </div>
            <div class="form-actions">
              <button @click="handleUpdate" class="btn btn-primary" :disabled="saving">
                <svg v-if="!saving" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/>
                  <polyline points="17 21 17 13 7 13 7 21"/>
                  <polyline points="7 3 7 8 15 8"/>
                </svg>
                {{ saving ? '保存中...' : '保存更改' }}
              </button>
            </div>
          </div>
        </section>

        <!-- Members Section -->
        <section class="settings-section animate-in animate-in-delay-1">
          <div class="section-header">
            <h2>成员管理</h2>
            <p>管理项目成员和权限</p>
          </div>
          <div class="section-body">
            <!-- Owner Transfer -->
            <div v-if="canManageMembers" class="owner-transfer">
              <div class="owner-transfer-info">
                <span class="owner-label">当前所有者：</span>
                <span class="owner-name">{{ project?.owner_name || '未知' }}</span>
              </div>
              <div class="owner-transfer-form">
                <select v-model="newOwnerId" class="form-select">
                  <option :value="0">选择新所有者...</option>
                  <option v-for="member in potentialOwners" :key="member.user_id" :value="member.user_id">
                    {{ member.display_name }}
                  </option>
                </select>
                <button
                  class="btn btn-warning"
                  :disabled="!newOwnerId || transferringOwner"
                  @click="handleTransferOwnership"
                >
                  {{ transferringOwner ? '转让中...' : '转让所有权' }}
                </button>
              </div>
            </div>
            <ProjectMembers :project-id="projectId" :can-manage="canManageMembers" />
          </div>
        </section>

        <!-- Danger Zone -->
        <section v-if="canDeleteProject" class="settings-section danger-zone animate-in animate-in-delay-2">
          <div class="section-header danger">
            <h2>危险区域</h2>
            <p>以下操作不可逆，请谨慎操作</p>
          </div>
          <div class="section-body">
            <div class="danger-action">
              <div class="danger-info">
                <div class="danger-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21 8v13H3V8"/>
                    <path d="M23 3H1v5h22V3z"/>
                    <path d="M10 12h4"/>
                  </svg>
                </div>
                <div>
                  <strong>归档项目</strong>
                  <p>归档后的项目将不再显示在项目列表中，但数据保留。</p>
                </div>
              </div>
              <button @click="handleArchive" class="btn btn-warning">
                {{ project?.is_archived ? '取消归档' : '归档项目' }}
              </button>
            </div>
            <div class="danger-action">
              <div class="danger-info">
                <div class="danger-icon danger">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3 6 5 6 21 6"/>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                    <line x1="10" y1="11" x2="10" y2="17"/>
                    <line x1="14" y1="11" x2="14" y2="17"/>
                  </svg>
                </div>
                <div>
                  <strong>删除项目</strong>
                  <p>删除项目将同时删除所有任务和评论，此操作不可撤销！</p>
                </div>
              </div>
              <button @click="handleDelete" class="btn btn-danger">删除项目</button>
            </div>
          </div>
        </section>
      </template>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getProject, updateProject, archiveProject, deleteProject } from '../../api/project'
import type { Project, ApiResponse } from '../../api/project'
import { getProjectMembers, updateMemberRole } from '../../api/membership'
import type { ApiResponse as MembershipApiResponse } from '../../api/membership'
import ProjectMembers from '../../components/project/ProjectMembers.vue'
import { usePermission } from '../../composables/usePermission'
import AppHeader from '../../components/AppHeader.vue'

const router = useRouter()
const route = useRoute()

const projectId = computed(() => Number(route.params.id))

const project = ref<Project | null>(null)
const form = ref({ name: '', description: '' })
const loading = ref(false)
const saving = ref(false)
const error = ref<string | null>(null)

const canDeleteProject = ref(false)
const canManageMembers = ref(false)

// 所有权转让
const newOwnerId = ref(0)
const transferringOwner = ref(false)
const potentialOwners = ref<{ user_id: number; display_name: string }[]>([])

async function checkPermissions() {
  const permission = usePermission()
  canDeleteProject.value = await permission.canDeleteProject(projectId.value)
  canManageMembers.value = await permission.canManageProjectMembers(projectId.value)
}

onMounted(async () => {
  await loadProject()
  await checkPermissions()
  if (canManageMembers.value) {
    await loadPotentialOwners()
  }
})

async function loadProject() {
  loading.value = true
  error.value = null
  try {
    const res = await getProject(projectId.value) as ApiResponse<Project>
    if (res.code === 0) {
      project.value = res.data
      form.value = {
        name: res.data.name,
        description: res.data.description || ''
      }
    } else {
      error.value = res.message || '加载失败'
    }
  } catch (e: any) {
    error.value = e.message || '加载失败'
  } finally {
    loading.value = false
  }
}

function goBack() {
  router.push(`/projects/${projectId.value}`)
}

async function handleUpdate() {
  if (!form.value.name.trim()) {
    error.value = '项目名称不能为空'
    return
  }

  saving.value = true
  error.value = null
  try {
    const res = await updateProject(projectId.value, {
      name: form.value.name.trim(),
      description: form.value.description.trim()
    }) as ApiResponse<Project>
    if (res.code === 0) {
      project.value = res.data
      ElMessage.success('保存成功')
    } else {
      error.value = res.message || '保存失败'
    }
  } catch (e: any) {
    error.value = e.message || '保存失败'
  } finally {
    saving.value = false
  }
}

async function handleArchive() {
  if (!project.value) return

  try {
    const res = await archiveProject(projectId.value) as ApiResponse<Project>
    if (res.code === 0) {
      project.value = res.data
    }
  } catch (e: any) {
    error.value = e.message || '操作失败'
  }
}

async function loadPotentialOwners() {
  try {
    const res = await getProjectMembers(projectId.value) as MembershipApiResponse<{ members: { user_id: number; display_name: string }[] }>
    if (res.code === 0) {
      // 排除当前所有者
      potentialOwners.value = res.data.members.filter(m => m.user_id !== project.value?.owner_id)
    }
  } catch (e) {
    console.error('加载成员失败', e)
  }
}

async function handleTransferOwnership() {
  if (!newOwnerId.value || !project.value) return

  if (!confirm(`确定要将项目所有权转让给选中成员吗？`)) return

  transferringOwner.value = true
  try {
    // 查找新所有者的角色ID
    const res = await getProjectMembers(projectId.value) as MembershipApiResponse<{ members: { user_id: number; role_id: number }[] }>
    if (res.code !== 0) {
      ElMessage.error('获取成员信息失败')
      return
    }
    const newOwnerMember = res.data.members.find(m => m.user_id === newOwnerId.value)
    if (!newOwnerMember) {
      ElMessage.error('未找到选中成员')
      return
    }

    // 1. 将新所有者角色更新为 owner
    const updateRes = await updateMemberRole(projectId.value, newOwnerId.value, 3) // 3 = owner role ID
    if (updateRes.code !== 0) {
      ElMessage.error(updateRes.message || '转让所有权失败')
      return
    }

    // 2. 将原所有者降级为 member
    if (project.value.owner_id) {
      const oldOwnerRes = await updateMemberRole(projectId.value, project.value.owner_id, 4) // 4 = member role ID
      if (oldOwnerRes.code !== 0) {
        ElMessage.error('转让成功，但降级原所有者失败')
      }
    }

    ElMessage.success('所有权转让成功')
    newOwnerId.value = 0
    await loadProject()
  } catch (e: any) {
    ElMessage.error(e.message || '转让所有权失败')
  } finally {
    transferringOwner.value = false
  }
}

async function handleDelete() {
  if (!confirm('确定要删除此项目吗？此操作不可撤销！')) return

  try {
    const res = await deleteProject(projectId.value) as ApiResponse<null>
    if (res.code === 0) {
      router.push('/projects')
    } else {
      error.value = res.message || '删除失败'
    }
  } catch (e: any) {
    error.value = e.message || '删除失败'
  }
}
</script>

<style scoped>
.project-settings-page {
  min-height: 100vh;
  background: var(--color-bg-base);
}

.settings-content {
  max-width: 800px;
  margin: 0 auto;
  padding: var(--space-8) var(--space-6);
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

/* Loading/Error States */
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

/* Settings Section */
.settings-section {
  background: var(--color-bg-surface);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-card);
  border: 1px solid var(--color-border);
  overflow: hidden;
}

.section-header {
  padding: var(--space-5) var(--space-6);
  border-bottom: 1px solid var(--color-border);
  background: var(--color-bg-muted);
}

.section-header h2 {
  font-family: var(--font-display);
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  margin: 0 0 var(--space-1);
}

.section-header p {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  margin: 0;
}

.section-header.danger {
  background: rgba(239, 71, 111, 0.05);
}

.section-header.danger h2 {
  color: var(--color-danger);
}

.section-body {
  padding: var(--space-6);
}

/* Form Styles */
.form-group {
  margin-bottom: var(--space-5);
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
  padding: var(--space-3) var(--space-4);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  font-family: var(--font-body);
  font-size: var(--font-size-base);
  color: var(--color-text-primary);
  background: var(--color-bg-surface);
  transition: all var(--transition-base);
}

.form-input:hover,
.form-textarea:hover {
  border-color: var(--color-border-hover);
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(91, 95, 199, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 100px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: var(--space-5);
  padding-top: var(--space-5);
  border-top: 1px solid var(--color-border);
}

/* Buttons */
.btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-5);
  border: none;
  border-radius: var(--radius-lg);
  cursor: pointer;
  font-family: var(--font-body);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  transition: all var(--transition-base);
}

.btn-primary {
  background: var(--color-primary);
  color: white;
  box-shadow: 0 2px 8px rgba(91, 95, 199, 0.3);
}

.btn-primary:hover:not(:disabled) {
  background: var(--color-primary-hover);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(91, 95, 199, 0.4);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-warning {
  background: #f5a623;
  color: white;
}

.btn-warning:hover {
  background: #e09612;
}

.btn-danger {
  background: var(--color-danger);
  color: white;
}

.btn-danger:hover {
  background: #e04560;
}

/* Danger Zone */
.danger-zone {
  border-color: rgba(239, 71, 111, 0.3);
}

.danger-action {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-4) 0;
  border-bottom: 1px solid var(--color-border);
}

.danger-action:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.danger-info {
  display: flex;
  align-items: flex-start;
  gap: var(--space-4);
}

.danger-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  background: rgba(255, 166, 43, 0.1);
  border-radius: var(--radius-lg);
  color: #f5a623;
  flex-shrink: 0;
}

.danger-icon.danger {
  background: rgba(239, 71, 111, 0.1);
  color: var(--color-danger);
}

.danger-info strong {
  display: block;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  margin-bottom: var(--space-1);
}

.danger-info p {
  margin: 0;
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
}

/* Responsive */
@media (max-width: 768px) {
  .settings-content {
    padding: var(--space-4);
  }

  .section-body {
    padding: var(--space-4);
  }

  .danger-action {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-4);
  }

  .btn {
    width: 100%;
    justify-content: center;
  }
}

/* Owner Transfer */
.owner-transfer {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  padding: var(--space-4);
  background: var(--color-bg-muted);
  border-radius: var(--radius-lg);
  margin-bottom: var(--space-4);
}

.owner-transfer-info {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.owner-label {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
}

.owner-name {
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
}

.owner-transfer-form {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  flex: 1;
}

.owner-transfer-form .form-select {
  flex: 1;
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-family: var(--font-body);
  font-size: var(--font-size-sm);
  background: var(--color-bg-surface);
}

.owner-transfer-form .btn-warning {
  background: #f5a623;
  color: white;
  white-space: nowrap;
}

.owner-transfer-form .btn-warning:hover:not(:disabled) {
  background: #e09612;
}

.owner-transfer-form .btn-warning:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>

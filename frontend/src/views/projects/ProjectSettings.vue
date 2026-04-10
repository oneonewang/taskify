<template>
  <div class="project-settings">
    <div class="page-header">
      <h1>项目设置</h1>
    </div>

    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>

    <div v-else class="settings-sections">
      <!-- 基本信息 -->
      <section class="settings-section">
        <h2>基本信息</h2>
        <div class="form-group">
          <label>项目名称</label>
          <input v-model="form.name" type="text" maxlength="100" />
        </div>
        <div class="form-group">
          <label>描述</label>
          <textarea v-model="form.description" rows="3" maxlength="500"></textarea>
        </div>
        <div class="form-actions">
          <button @click="handleUpdate" class="btn btn-primary" :disabled="saving">
            {{ saving ? '保存中...' : '保存更改' }}
          </button>
        </div>
      </section>

      <!-- 成员管理 -->
      <section class="settings-section">
        <h2>成员管理</h2>
        <ProjectMembers :project-id="projectId" :can-manage="canManageMembers" />
      </section>

      <!-- 归档操作 -->
      <section v-if="canDeleteProject" class="settings-section danger-zone">
        <h2>危险区域</h2>
        <div class="danger-action">
          <div>
            <strong>归档项目</strong>
            <p>归档后的项目将不再显示在项目列表中。</p>
          </div>
          <button @click="handleArchive" class="btn btn-warning">
            {{ project?.is_archived ? '取消归档' : '归档项目' }}
          </button>
        </div>
        <div class="danger-action">
          <div>
            <strong>删除项目</strong>
            <p>删除项目将同时删除所有任务和评论，此操作不可撤销！</p>
          </div>
          <button @click="handleDelete" class="btn btn-danger">删除项目</button>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getProject, updateProject, archiveProject, deleteProject } from '../../api/project'
import type { Project, ApiResponse } from '../../api/project'
import ProjectMembers from '../../components/project/ProjectMembers.vue'
import { usePermission } from '../../composables/usePermission'

const router = useRouter()
const route = useRoute()
const projectId = computed(() => Number(route.params.id))

const project = ref<Project | null>(null)
const form = ref({ name: '', description: '' })
const loading = ref(false)
const saving = ref(false)
const error = ref<string | null>(null)

// 动态检查删除权限
const canDeleteProject = ref(false)
const canManageMembers = ref(false)

async function checkPermissions() {
  const permission = usePermission()
  canDeleteProject.value = await permission.canDeleteProject(projectId.value)
  canManageMembers.value = await permission.canManageProjectMembers(projectId.value)
}

onMounted(async () => {
  await loadProject()
  await checkPermissions()
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
.project-settings {
  max-width: 800px;
  margin: 0 auto;
  padding: 24px;
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
  text-align: center;
  padding: 40px;
  color: #666;
}

.error {
  color: #dc3545;
}

.settings-sections {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.settings-section {
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.settings-section h2 {
  margin: 0 0 20px;
  font-size: 16px;
  color: #333;
  padding-bottom: 12px;
  border-bottom: 1px solid #eee;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 14px;
  color: #666;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #1976d2;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}

.danger-zone {
  border: 1px solid #dc3545;
}

.danger-zone h2 {
  color: #dc3545;
}

.danger-action {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #eee;
}

.danger-action:last-child {
  border-bottom: none;
}

.danger-action p {
  margin: 4px 0 0;
  font-size: 13px;
  color: #666;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

.btn-primary {
  background: #1976d2;
  color: #fff;
}

.btn-warning {
  background: #f5a623;
  color: #fff;
}

.btn-danger {
  background: #dc3545;
  color: #fff;
}
</style>

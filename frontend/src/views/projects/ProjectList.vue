<template>
  <div class="project-list-page">
    <div class="page-header">
      <h1>我的项目</h1>
      <button @click="showCreateDialog = true" class="btn btn-primary">
        创建项目
      </button>
    </div>

    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>

    <div v-else-if="projects.length === 0" class="empty-state">
      <p>暂无项目</p>
      <button @click="showCreateDialog = true" class="btn btn-primary">
        创建第一个项目
      </button>
    </div>

    <div v-else class="projects-grid">
      <div
        v-for="project in projects"
        :key="project.id"
        class="project-card"
        @click="goToProject(project.id)"
      >
        <div class="project-card-header">
          <h3>{{ project.name }}</h3>
          <span v-if="project.is_archived" class="archived-badge">已归档</span>
        </div>
        <p class="project-description">{{ project.description || '无描述' }}</p>
        <div class="project-meta">
          <span v-if="project.task_counts" class="task-count">
            {{ getTotalTasks(project.task_counts) }} 个任务
          </span>
          <span class="created-at">
            创建于 {{ formatDate(project.created_at) }}
          </span>
        </div>
      </div>
    </div>

    <!-- 创建项目对话框 -->
    <CreateProjectDialog
      v-model:visible="showCreateDialog"
      @success="onProjectCreated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getProjects } from '../../api/project'
import CreateProjectDialog from '../../components/project/CreateProjectDialog.vue'
import type { Project, ApiResponse } from '../../api/project'

const router = useRouter()
const route = useRoute()
const projects = ref<Project[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const showCreateDialog = ref(false)

onMounted(async () => {
  // 检查是否通过路由跳转打开创建对话框
  if (route.name === 'ProjectNew' || route.query.new === 'true') {
    showCreateDialog.value = true
  }
  await loadProjects()
})

async function loadProjects() {
  loading.value = true
  error.value = null
  try {
    const res = await getProjects() as ApiResponse<Project[]>
    if (res.code === 0) {
      projects.value = res.data
    } else {
      error.value = res.message || '加载失败'
    }
  } catch (e: any) {
    error.value = e.message || '加载失败'
  } finally {
    loading.value = false
  }
}

function getTotalTasks(taskCounts: any): number {
  if (!taskCounts) return 0
  return (taskCounts.todo || 0) + (taskCounts.in_progress || 0) +
         (taskCounts.review || 0) + (taskCounts.done || 0)
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}

function goToProject(projectId: number) {
  router.push(`/projects/${projectId}`)
}

function onProjectCreated(newProject: Project) {
  projects.value.unshift(newProject)
  showCreateDialog.value = false
}
</script>

<style scoped>
.project-list-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
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

.loading, .error, .empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #666;
}

.error {
  color: #dc3545;
}

.empty-state p {
  margin-bottom: 16px;
  color: #999;
}

.projects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.project-card {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  cursor: pointer;
  transition: box-shadow 0.2s, transform 0.2s;
}

.project-card:hover {
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  transform: translateY(-2px);
}

.project-card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.project-card-header h3 {
  margin: 0;
  font-size: 16px;
  color: #333;
}

.archived-badge {
  padding: 2px 6px;
  background: #e3f2fd;
  color: #1976d2;
  border-radius: 4px;
  font-size: 11px;
}

.project-description {
  margin: 0 0 12px;
  color: #666;
  font-size: 14px;
  line-height: 1.4;
  height: 2.8em;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.project-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #999;
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

.btn-primary:hover {
  background: #1565c0;
}
</style>

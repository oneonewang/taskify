<template>
  <div class="project-list-page">
    <AppHeader title="项目" />

    <main class="page-content">
      <!-- Tab Bar -->
      <div class="tab-bar animate-in">
        <button
          :class="['tab-btn', { active: activeTab === 'my' }]"
          @click="switchTab('my')"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
          我参与的项目
        </button>
        <button
          :class="['tab-btn', { active: activeTab === 'all' }]"
          @click="switchTab('all')"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="7" height="7"/>
            <rect x="14" y="3" width="7" height="7"/>
            <rect x="14" y="14" width="7" height="7"/>
            <rect x="3" y="14" width="7" height="7"/>
          </svg>
          全部项目
        </button>
        <div class="tab-spacer"></div>
        <button @click="showCreateDialog = true" class="btn btn-primary">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          创建项目
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="loading-state animate-in">
        <div class="loading-spinner"></div>
        <span>加载中...</span>
      </div>

      <!-- Error State -->
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

      <!-- Empty State -->
      <div v-else-if="projects.length === 0" class="empty-state animate-in">
        <div class="empty-illustration">
          <svg width="80" height="80" viewBox="0 0 80 80" fill="none">
            <rect x="10" y="20" width="60" height="50" rx="8" stroke="currentColor" stroke-width="2" stroke-dasharray="4 4"/>
            <rect x="20" y="10" width="40" height="30" rx="4" fill="var(--color-bg-muted)" stroke="currentColor" stroke-width="2"/>
            <line x1="30" y1="20" x2="50" y2="20" stroke="currentColor" stroke-width="2"/>
            <line x1="30" y1="28" x2="45" y2="28" stroke="currentColor" stroke-width="2"/>
          </svg>
        </div>
        <h3>暂无项目</h3>
        <p v-if="activeTab === 'my'">您还没有参与任何项目</p>
        <p v-else>创建您的第一个项目开始吧</p>
        <button @click="showCreateDialog = true" class="btn btn-primary">
          创建第一个项目
        </button>
      </div>

      <!-- Projects Grid -->
      <div v-else class="projects-grid">
        <div
          v-for="(project, index) in projects"
          :key="project.id"
          class="project-card animate-in"
          :style="{ animationDelay: `${index * 0.05}s` }"
          @click="goToProject(project.id)"
        >
          <!-- Card Accent -->
          <div class="card-accent"></div>

          <!-- Card Header -->
          <div class="project-card-header">
            <div class="project-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                <line x1="3" y1="9" x2="21" y2="9"/>
                <line x1="9" y1="21" x2="9" y2="9"/>
              </svg>
            </div>
            <h3 class="project-name">{{ project.name }}</h3>
            <span v-if="project.is_archived" class="archived-badge">已归档</span>
          </div>

          <!-- Project Description -->
          <p class="project-description">{{ project.description || '暂无描述' }}</p>

          <!-- Task Progress -->
          <div v-if="project.task_counts" class="task-progress">
            <div class="progress-bar">
              <div
                class="progress-fill"
                :style="{ width: `${getProgressPercent(project.task_counts)}%` }"
              ></div>
            </div>
            <span class="progress-text">{{ getDoneTasks(project.task_counts) }} / {{ getTotalTasks(project.task_counts) }}</span>
          </div>

          <!-- Card Footer -->
          <div class="project-card-footer">
            <div class="task-count-badge">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="9 11 12 14 22 4"/>
                <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"/>
              </svg>
              {{ getTotalTasks(project.task_counts) }} 任务
            </div>
            <span class="created-at">
              {{ formatDate(project.created_at) }}
            </span>
          </div>
        </div>
      </div>
    </main>

    <!-- Create Project Dialog -->
    <CreateProjectDialog
      v-model:visible="showCreateDialog"
      @success="onProjectCreated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getProjects, getMyProjects } from '../../api/project'
import CreateProjectDialog from '../../components/project/CreateProjectDialog.vue'
import AppHeader from '../../components/AppHeader.vue'
import type { Project, ApiResponse } from '../../api/project'

const router = useRouter()
const route = useRoute()
const projects = ref<Project[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const showCreateDialog = ref(false)
const activeTab = ref<'my' | 'all'>('my')

onMounted(async () => {
  if (route.name === 'ProjectNew' || route.query.new === 'true') {
    showCreateDialog.value = true
  }
  await loadProjects()
})

async function loadProjects() {
  loading.value = true
  error.value = null
  try {
    const res = activeTab.value === 'my'
      ? await getMyProjects() as ApiResponse<Project[]>
      : await getProjects() as ApiResponse<Project[]>
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

function switchTab(tab: 'my' | 'all') {
  activeTab.value = tab
  loadProjects()
}

function getTotalTasks(taskCounts: any): number {
  if (!taskCounts) return 0
  return (taskCounts.todo || 0) + (taskCounts.in_progress || 0) +
         (taskCounts.review || 0) + (taskCounts.done || 0)
}

function getDoneTasks(taskCounts: any): number {
  if (!taskCounts) return 0
  return taskCounts.done || 0
}

function getProgressPercent(taskCounts: any): number {
  const total = getTotalTasks(taskCounts)
  if (total === 0) return 0
  return Math.round((getDoneTasks(taskCounts) / total) * 100)
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
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
  min-height: 100vh;
  background: var(--color-bg-base);
}

.page-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: var(--space-6);
}

/* Tab Bar */
.tab-bar {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  background: var(--color-bg-surface);
  padding: var(--space-2);
  border-radius: var(--radius-xl);
  border: 1px solid var(--color-border);
  margin-bottom: var(--space-6);
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-4);
  background: transparent;
  border: none;
  border-radius: var(--radius-lg);
  cursor: pointer;
  font-family: var(--font-body);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-secondary);
  transition: all var(--transition-base);
}

.tab-btn:hover {
  color: var(--color-text-primary);
  background: var(--color-bg-muted);
}

.tab-btn.active {
  background: var(--color-primary);
  color: white;
}

.tab-btn.active svg {
  opacity: 1;
}

.tab-btn svg {
  opacity: 0.7;
}

.tab-spacer {
  flex: 1;
}

/* Buttons */
.btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
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

.btn-primary:hover {
  background: var(--color-primary-hover);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(91, 95, 199, 0.4);
}

/* Loading State */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-16);
  color: var(--color-text-muted);
  gap: var(--space-4);
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

/* Error State */
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-16);
  color: var(--color-danger);
  gap: var(--space-4);
  text-align: center;
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

/* Empty State */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-16);
  text-align: center;
}

.empty-illustration {
  color: var(--color-text-muted);
  margin-bottom: var(--space-6);
}

.empty-state h3 {
  font-family: var(--font-display);
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  margin: 0 0 var(--space-2);
}

.empty-state p {
  color: var(--color-text-muted);
  margin: 0 0 var(--space-6);
}

/* Projects Grid */
.projects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: var(--space-5);
}

/* Project Card */
.project-card {
  position: relative;
  background: var(--color-bg-surface);
  border-radius: var(--radius-xl);
  padding: var(--space-6);
  box-shadow: var(--shadow-card);
  cursor: pointer;
  transition: all var(--transition-base);
  border: 1px solid var(--color-border);
  overflow: hidden;
}

.project-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-card-hover);
  border-color: var(--color-primary-light);
}

.project-card:hover .card-accent {
  opacity: 1;
}

/* Card Accent */
.card-accent {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--color-primary) 0%, var(--color-primary-light) 100%);
  opacity: 0;
  transition: opacity var(--transition-base);
}

/* Card Header */
.project-card-header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-3);
}

.project-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  background: var(--color-bg-muted);
  border-radius: var(--radius-lg);
  color: var(--color-primary);
}

.project-name {
  flex: 1;
  font-family: var(--font-display);
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  margin: 0;
}

.archived-badge {
  padding: var(--space-1) var(--space-3);
  background: var(--color-bg-muted);
  color: var(--color-text-muted);
  border-radius: var(--radius-full);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
}

/* Project Description */
.project-description {
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  line-height: var(--line-height-relaxed);
  margin: 0 0 var(--space-4);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Task Progress */
.task-progress {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-4);
}

.progress-bar {
  flex: 1;
  height: 6px;
  background: var(--color-bg-muted);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--color-primary) 0%, var(--color-success) 100%);
  border-radius: var(--radius-full);
  transition: width var(--transition-slow);
}

.progress-text {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
  font-weight: var(--font-weight-medium);
}

/* Card Footer */
.project-card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: var(--space-4);
  border-top: 1px solid var(--color-border);
}

.task-count-badge {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

.task-count-badge svg {
  color: var(--color-primary);
}

.created-at {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
}

/* Responsive */
@media (max-width: 768px) {
  .page-content {
    padding: var(--space-4);
  }

  .tab-bar {
    flex-wrap: wrap;
  }

  .tab-spacer {
    width: 100%;
    height: 0;
  }

  .btn {
    width: 100%;
    justify-content: center;
  }

  .projects-grid {
    grid-template-columns: 1fr;
  }
}
</style>

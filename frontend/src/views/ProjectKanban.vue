<template>
  <div class="project-kanban">
    <AppHeader :title="selectedProjectName">
      <template #center>
        <div class="project-selector" v-if="projects.length > 0">
          <el-select
            v-model="selectedProjectId"
            placeholder="选择项目"
            size="default"
            class="project-select"
            @change="onProjectChange"
          >
            <template #prefix>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                <line x1="3" y1="9" x2="21" y2="9"/>
                <line x1="9" y1="21" x2="9" y2="9"/>
              </svg>
            </template>
            <el-option
              v-for="project in projects"
              :key="project.id"
              :label="project.name"
              :value="project.id"
            />
          </el-select>
        </div>
      </template>
      <template #actions>
        <router-link
          v-if="selectedProjectId && isOwner"
          :to="`/projects/${selectedProjectId}/settings`"
          class="settings-btn"
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3"/>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
          </svg>
          <span>项目设置</span>
        </router-link>
      </template>
    </AppHeader>

    <main class="kanban-main">
      <div class="content-wrapper">
        <div class="main-content">
          <div v-if="projectStore.loading" class="loading-state animate-in">
            <div class="loading-spinner"></div>
            <span>加载中...</span>
          </div>
          <div v-else-if="projectStore.error" class="error-state animate-in">
            <div class="error-icon">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <line x1="12" y1="8" x2="12" y2="12"/>
                <line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
            </div>
            <span>{{ projectStore.error }}</span>
          </div>
          <template v-else>
            <KanbanBoard
              :columns="projectStore.columns"
              :tasks="projectStore.tasks"
              :current-user-id="currentUserId"
              :project-id="selectedProjectId"
              @task-click="onTaskClick"
              @task-moved="onTaskMoved"
              @tasks-updated="onTasksUpdated"
            />
          </template>
        </div>
        <aside class="sidebar">
          <TeamMembers ref="teamMembersRef" :current-user-id="currentUserId" :project-id="selectedProjectId" />
        </aside>
      </div>
    </main>

    <!-- Task Detail Dialog -->
    <TaskDetail
      v-model:visible="taskDetailVisible"
      :task="selectedTask"
      @task-updated="onTaskUpdated"
      @close="onTaskDetailClose"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '../stores/auth'
import { useProjectStore } from '../stores/project'
import { useSSEStoreUpdater } from '../composables/useSSEStoreUpdater'
import { getProjects, getProjectTask } from '../api/project'
import { updateTaskStatus, type UpdateTaskStatusRequest } from '../api/task'
import type { Project, Task } from '../api/project'
import KanbanBoard from '../components/kanban/KanbanBoard.vue'
import TeamMembers from '../components/TeamMembers.vue'
import TaskDetail from '../components/task/TaskDetail.vue'
import AppHeader from '../components/AppHeader.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const projectStore = useProjectStore()

// Initialize SSE connection
useSSEStoreUpdater()

const projects = ref<Project[]>([])
const selectedProjectId = ref<number>(0)
const taskDetailVisible = ref(false)
const selectedTask = ref<Task | null>(null)
const teamMembersRef = ref<InstanceType<typeof TeamMembers> | null>(null)

const currentUserId = computed(() => authStore.user?.id ?? null)

// Selected project name for header title
const selectedProjectName = computed(() => {
  if (!selectedProjectId.value) return ''
  const project = projects.value.find(p => p.id === selectedProjectId.value)
  return project?.name || projectStore.currentProject?.name || ''
})

// Check if current user is owner of selected project
const isOwner = computed(() => {
  if (!selectedProjectId.value || !currentUserId.value) return false
  return projectStore.currentProject?.owner_id === currentUserId.value
})

// Get project ID from route
const projectIdFromRoute = computed(() => Number(route.params.id))

onMounted(async () => {
  if (projectIdFromRoute.value) {
    selectedProjectId.value = projectIdFromRoute.value
    await projectStore.loadProject(selectedProjectId.value)
    teamMembersRef.value?.loadMembers(selectedProjectId.value)
    // Handle direct task link (share feature)
    const taskId = Number(route.params.taskId)
    if (taskId) {
      await loadTaskFromUrl(taskId)
    }
  } else {
    try {
      const res = await getProjects()
      if (res.code === 0 && res.data.projects && res.data.projects.length > 0) {
        router.replace(`/projects/${res.data.projects[0].id}`)
      }
    } catch (e) {
      console.error('Failed to load projects:', e)
    }
  }
})

// Watch for route changes
watch(() => route.params.id, async (newId) => {
  if (newId && Number(newId) !== selectedProjectId.value) {
    selectedProjectId.value = Number(newId)
    await projectStore.loadProject(selectedProjectId.value)
    teamMembersRef.value?.loadMembers(selectedProjectId.value)
  }
})

// Watch for taskId param (direct link access via share)
watch(() => route.params.taskId, async (newTaskId) => {
  if (newTaskId && selectedProjectId.value) {
    await loadTaskFromUrl(Number(newTaskId))
  }
})

// Load task from URL (for share links)
async function loadTaskFromUrl(taskId: number) {
  try {
    const res = await getProjectTask(selectedProjectId.value, taskId)
    if (res.code === 0) {
      selectedTask.value = res.data
      taskDetailVisible.value = true
    } else if (res.code === 403) {
      ElMessage.error('您不是该项目成员，无权访问此任务')
    } else if (res.code === 404) {
      ElMessage.error('任务不存在')
    }
  } catch (e: any) {
    if (e.response?.status === 403) {
      ElMessage.error('您不是该项目成员，无权访问此任务')
    } else if (e.response?.status === 404) {
      ElMessage.error('任务不存在')
    } else {
      ElMessage.error('加载任务失败')
    }
  }
}

async function onProjectChange(projectId: number) {
  router.push(`/projects/${projectId}`)
  teamMembersRef.value?.loadMembers(projectId)
}

// Task movement handler
async function onTaskMoved(taskId: number, newStatus: string, newPosition: number) {
  const task = projectStore.tasks.find(t => t.id === taskId)
  if (!task) return

  const oldStatus = task.status
  const oldPosition = task.position

  // Optimistic update
  projectStore.updateTask(taskId, { status: newStatus, position: newPosition })

  try {
    const data: UpdateTaskStatusRequest = {
      status: newStatus,
      position: newPosition
    }
    const res = await updateTaskStatus(taskId, data)
    if (res.code !== 0) {
      // Rollback on failure
      projectStore.updateTask(taskId, { status: oldStatus, position: oldPosition })
      ElMessage.error(res.message || '移动任务失败')
    }
  } catch (e: any) {
    // Rollback on error
    projectStore.updateTask(taskId, { status: oldStatus, position: oldPosition })
    ElMessage.error(e.message || '移动任务失败')
  }
}

// Task list update handler (after drag reorder)
function onTasksUpdated(updatedTasks: Task[]) {
  projectStore.setTasks(updatedTasks)
}

function onTaskClick(task: Task) {
  selectedTask.value = task
  taskDetailVisible.value = true
}

function onTaskUpdated(updatedTask: Task) {
  projectStore.updateTask(updatedTask.id, updatedTask)
}

// Task detail dialog closed - navigate back to kanban view
function onTaskDetailClose() {
  if (selectedProjectId.value) {
    router.push(`/projects/${selectedProjectId.value}`)
  }
}
</script>

<style scoped>
.project-kanban {
  min-height: 100vh;
  background: var(--color-bg-base);
  display: flex;
  flex-direction: column;
}

.kanban-main {
  flex: 1;
  padding: var(--space-6);
}

.content-wrapper {
  display: flex;
  gap: var(--space-6);
  align-items: flex-start;
  height: 100%;
}

.main-content {
  flex: 1;
  min-width: 0;
}

.sidebar {
  flex: 0 0 280px;
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

/* Project Selector */
.project-selector {
  min-width: 200px;
}

.project-select {
  width: 200px;
}

.project-select :deep(.el-input__wrapper) {
  border-radius: var(--radius-lg);
  box-shadow: 0 0 0 1px var(--color-border);
}

.project-select :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px var(--color-border-hover);
}

.project-select :deep(.el-input.is-focus .el-input__wrapper) {
  box-shadow: 0 0 0 2px var(--color-primary-light) !important;
}

/* Settings Button */
.settings-btn {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-3);
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  color: var(--color-text-secondary);
  text-decoration: none;
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  transition: all var(--transition-base);
}

.settings-btn:hover {
  background: var(--color-bg-muted);
  border-color: var(--color-border-hover);
  color: var(--color-text-primary);
}

.settings-btn svg {
  opacity: 0.7;
}

.settings-btn:hover svg {
  opacity: 1;
}

/* Responsive */
@media (max-width: 1024px) {
  .sidebar {
    display: none;
  }
}

@media (max-width: 768px) {
  .kanban-main {
    padding: var(--space-4);
  }

  .project-selector {
    display: none;
  }
}
</style>

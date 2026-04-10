<template>
  <div class="project-kanban">
    <div class="header">
      <h1>{{ projectStore.currentProject?.name || '加载中...' }}</h1>
      <div class="header-actions">
        <el-select v-model="selectedProjectId" placeholder="选择项目" @change="onProjectChange">
          <el-option
            v-for="project in projects"
            :key="project.id"
            :label="project.name"
            :value="project.id"
          />
        </el-select>
        <el-button @click="goToSettings">项目设置</el-button>
        <el-button @click="handleLogout">退出登录</el-button>
      </div>
    </div>

    <div class="content-wrapper">
      <div class="main-content">
        <div v-if="projectStore.loading" class="loading">加载中...</div>
        <div v-else-if="projectStore.error" class="error">{{ projectStore.error }}</div>
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

    <!-- 任务详情对话框 -->
    <TaskDetail
      v-model:visible="taskDetailVisible"
      :task="selectedTask"
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
import { getProjects } from '../api/project'
import { updateTaskStatus, type UpdateTaskStatusRequest } from '../api/task'
import type { Project, Task } from '../api/project'
import KanbanBoard from '../components/kanban/KanbanBoard.vue'
import TeamMembers from '../components/TeamMembers.vue'
import TaskDetail from '../components/task/TaskDetail.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const projectStore = useProjectStore()

// 初始化 SSE 连接
useSSEStoreUpdater()

const projects = ref<Project[]>([])
const selectedProjectId = ref<number>(0)
const taskDetailVisible = ref(false)
const selectedTask = ref<Task | null>(null)
const teamMembersRef = ref<InstanceType<typeof TeamMembers> | null>(null)

const currentUserId = computed(() => authStore.user?.id ?? null)

// 从路由获取项目ID
const projectIdFromRoute = computed(() => Number(route.params.id))

onMounted(async () => {
  if (projectIdFromRoute.value) {
    selectedProjectId.value = projectIdFromRoute.value
    await projectStore.loadProject(selectedProjectId.value)
    teamMembersRef.value?.loadMembers(selectedProjectId.value)
  } else {
    // 如果没有项目ID，先获取项目列表
    try {
      const res = await getProjects()
      if (res.code === 0 && res.data.length > 0) {
        router.replace(`/projects/${res.data[0].id}`)
      }
    } catch (e) {
      console.error('Failed to load projects:', e)
    }
  }
})

// 监听路由变化
watch(() => route.params.id, async (newId) => {
  if (newId && Number(newId) !== selectedProjectId.value) {
    selectedProjectId.value = Number(newId)
    await projectStore.loadProject(selectedProjectId.value)
    teamMembersRef.value?.loadMembers(selectedProjectId.value)
  }
})

async function onProjectChange(projectId: number) {
  router.push(`/projects/${projectId}`)
  teamMembersRef.value?.loadMembers(projectId)
}

function goToSettings() {
  router.push(`/projects/${selectedProjectId.value}/settings`)
}

// 任务移动处理
async function onTaskMoved(taskId: number, newStatus: string, newPosition: number) {
  // 乐观更新：立即更新本地状态
  const task = projectStore.tasks.find(t => t.id === taskId)
  if (!task) return

  const oldStatus = task.status
  const oldPosition = task.position

  // 立即更新 UI
  projectStore.updateTask(taskId, { status: newStatus, position: newPosition })

  // 调用 API
  try {
    const data: UpdateTaskStatusRequest = {
      status: newStatus,
      position: newPosition
    }
    const res = await updateTaskStatus(taskId, data)
    if (res.code !== 0) {
      // API 失败，回滚
      projectStore.updateTask(taskId, { status: oldStatus, position: oldPosition })
      ElMessage.error(res.message || '移动任务失败')
    }
  } catch (e: any) {
    // 网络错误，回滚
    projectStore.updateTask(taskId, { status: oldStatus, position: oldPosition })
    ElMessage.error(e.message || '移动任务失败')
  }
}

// 任务列表更新（拖拽后重新排序）
function onTasksUpdated(updatedTasks: Task[]) {
  projectStore.setTasks(updatedTasks)
}

function onTaskClick(task: Task) {
  selectedTask.value = task
  taskDetailVisible.value = true
}

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.project-kanban {
  width: 100%;
  min-height: 100vh;
  background: #f5f7fa;
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 0 10px;
}

.header h1 {
  font-size: 24px;
  color: #333;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.content-wrapper {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

.main-content {
  flex: 1;
  min-width: 0;
}

.sidebar {
  flex: 0 0 260px;
}

.loading,
.error {
  text-align: center;
  padding: 40px;
  color: #666;
  font-size: 16px;
}

.error {
  color: #e6a23c;
}
</style>

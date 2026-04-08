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
        <el-button @click="goBack">切换用户</el-button>
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
        <TeamMembers :current-user-id="currentUserId" />
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
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../stores/user'
import { useProjectStore } from '../stores/project'
import { getProjects } from '../api/project'
import { updateTaskStatus, type UpdateTaskStatusRequest } from '../api/task'
import type { Project, Task } from '../api/project'
import KanbanBoard from '../components/kanban/KanbanBoard.vue'
import TeamMembers from '../components/TeamMembers.vue'
import TaskDetail from '../components/task/TaskDetail.vue'

const router = useRouter()
const userStore = useUserStore()
const projectStore = useProjectStore()

const projects = ref<Project[]>([])
const selectedProjectId = ref<number>(1)
const taskDetailVisible = ref(false)
const selectedTask = ref<Task | null>(null)

const currentUserId = computed(() => userStore.currentUserId)

onMounted(async () => {
  try {
    const res = await getProjects()
    if (res.success) {
      projects.value = res.data
      if (projects.value.length > 0) {
        selectedProjectId.value = projects.value[0].id
        await projectStore.loadProject(selectedProjectId.value)
      }
    }
  } catch (e) {
    console.error('Failed to load projects:', e)
  }
})

async function onProjectChange(projectId: number) {
  await projectStore.loadProject(projectId)
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
    if (!res.success) {
      // API 失败，回滚
      projectStore.updateTask(taskId, { status: oldStatus, position: oldPosition })
      ElMessage.error(res.error?.message || '移动任务失败')
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

function goBack() {
  router.push('/')
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

<template>
  <div class="project-kanban">
    <div class="header">
      <h1>{{ project?.name || '加载中...' }}</h1>
      <el-button @click="goBack">切换用户</el-button>
    </div>
    <div class="kanban-board">
      <div v-for="column in columns" :key="column.status" class="kanban-column">
        <div class="column-header">
          <span class="column-title">{{ column.name }}</span>
          <span class="column-count">{{ getTasksByStatus(column.status).length }}</span>
        </div>
        <div class="column-content">
          <div
            v-for="task in getTasksByStatus(column.status)"
            :key="task.id"
            class="task-card"
            :class="{ 'my-task': task.assignee.id === currentUserId }"
          >
            <div class="task-title">{{ task.title }}</div>
            <div class="task-assignee">
              <span
                class="assignee-avatar"
                :style="{ backgroundColor: task.assignee.avatar }"
              >
                {{ task.assignee.name.charAt(0) }}
              </span>
              <span class="assignee-name">{{ task.assignee.name }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { getProject } from '../api/project'
import type { Project, Task } from '../api/project'

const router = useRouter()
const userStore = useUserStore()

const project = ref<Project | null>(null)
const tasks = ref<Task[]>([])
const loading = ref(true)

const currentUserId = computed(() => userStore.currentUserId)

const columns = [
  { status: 'todo', name: '待办' },
  { status: 'in_progress', name: '进行中' },
  { status: 'review', name: '审核中' },
  { status: 'done', name: '已完成' }
]

onMounted(async () => {
  // 默认加载第一个项目
  try {
    const res = await getProject(1)
    if (res.success) {
      project.value = {
        id: res.data.id,
        name: res.data.name,
        description: res.data.description,
        created_at: res.data.created_at
      }
      tasks.value = res.data.tasks || []
    }
  } catch (e) {
    console.error('Failed to load project:', e)
  } finally {
    loading.value = false
  }
})

function getTasksByStatus(status: string): Task[] {
  return tasks.value.filter(t => t.status === status)
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

.kanban-board {
  display: flex;
  gap: 16px;
  overflow-x: auto;
  padding-bottom: 20px;
}

.kanban-column {
  flex: 0 0 280px;
  background: #eef1f5;
  border-radius: 12px;
  padding: 12px;
}

.column-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding: 0 4px;
}

.column-title {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.column-count {
  background: #dfe4ea;
  color: #666;
  padding: 2px 8px;
  border-radius: 10px;
  font-size: 12px;
}

.column-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-height: 100px;
}

.task-card {
  background: white;
  border-radius: 8px;
  padding: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.2s;
}

.task-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.task-card.my-task {
  border: 2px solid #667eea;
}

.task-title {
  font-size: 14px;
  color: #333;
  margin-bottom: 8px;
}

.task-assignee {
  display: flex;
  align-items: center;
  gap: 6px;
}

.assignee-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 12px;
}

.assignee-name {
  font-size: 12px;
  color: #666;
}
</style>

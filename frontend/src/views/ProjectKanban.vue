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
            @task-click="onTaskClick"
          />
        </template>
      </div>
      <aside class="sidebar">
        <TeamMembers :current-user-id="currentUserId" />
      </aside>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { useProjectStore } from '../stores/project'
import { getProjects } from '../api/project'
import type { Project, Task } from '../api/project'
import KanbanBoard from '../components/kanban/KanbanBoard.vue'
import TeamMembers from '../components/TeamMembers.vue'

const router = useRouter()
const userStore = useUserStore()
const projectStore = useProjectStore()

const projects = ref<Project[]>([])
const selectedProjectId = ref<number>(1)

const currentUserId = computed(() => userStore.currentUserId)

onMounted(async () => {
  // 加载项目列表
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

function onTaskClick(task: Task) {
  console.log('Task clicked:', task)
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

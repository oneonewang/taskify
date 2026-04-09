<template>
  <div class="kanban-board">
    <KanbanColumn
      v-for="column in columns"
      :key="column.status"
      :name="column.name"
      :status="column.status"
      :tasks="getTasksByStatus(column.status)"
      :current-user-id="currentUserId"
      :project-id="projectId"
      @task-click="(task) => emit('taskClick', task)"
      @task-moved="(taskId, newStatus, newPosition) => emit('taskMoved', taskId, newStatus, newPosition)"
      @tasks-updated="(tasks) => onTasksUpdated(column.status, tasks)"
    />
  </div>
</template>

<script setup lang="ts">
import type { Task } from '../../api/project'
import KanbanColumn from './KanbanColumn.vue'

const props = defineProps<{
  columns: Array<{ status: string; name: string }>
  tasks: Task[]
  currentUserId: number | null
  projectId: number
}>()

const emit = defineEmits<{
  (e: 'taskClick', task: Task): void
  (e: 'taskMoved', taskId: number, newStatus: string, newPosition: number): void
  (e: 'tasksUpdated', tasks: Task[]): void
}>()

function getTasksByStatus(status: string): Task[] {
  return props.tasks.filter(t => t.status === status)
}

// 当某一列的任务列表更新时
function onTasksUpdated(status: string, updatedTasks: Task[]) {
  // 合并更新后的任务到总任务列表
  // 注意：任务可能同时出现在 otherTasks 和 updatedTasks 中（如果状态还未更新）
  // 需要去重，以 updatedTasks 中的任务状态为准
  const otherTasks = props.tasks.filter(t => t.status !== status)
  const merged = [...otherTasks, ...updatedTasks]
  // 按 ID 去重，后出现的优先（updatedTasks 中的任务状态更新）
  const seen = new Set<number>()
  const deduplicated = merged.filter(task => {
    if (seen.has(task.id)) return false
    seen.add(task.id)
    return true
  })
  emit('tasksUpdated', deduplicated)
}
</script>

<style scoped>
.kanban-board {
  display: flex;
  gap: 16px;
  overflow-x: auto;
  padding-bottom: 20px;
}
</style>

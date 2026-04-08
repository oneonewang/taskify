<template>
  <div class="kanban-board">
    <KanbanColumn
      v-for="column in columns"
      :key="column.status"
      :name="column.name"
      :status="column.status"
      :tasks="getTasksByStatus(column.status)"
      :current-user-id="currentUserId"
      @task-click="$emit('taskClick', $event)"
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
}>()

defineEmits<{
  taskClick: [task: Task]
}>()

function getTasksByStatus(status: string): Task[] {
  return props.tasks.filter(t => t.status === status)
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
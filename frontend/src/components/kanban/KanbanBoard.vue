<template>
  <div class="kanban-board">
    <KanbanColumn
      v-for="(column, index) in columns"
      :key="column.status"
      :name="column.name"
      :status="column.status"
      :tasks="getTasksByStatus(column.status)"
      :current-user-id="currentUserId"
      :project-id="projectId"
      :style="{ animationDelay: `${index * 0.1}s` }"
      class="column-animate-in"
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

// When a column's task list is updated
function onTasksUpdated(status: string, updatedTasks: Task[]) {
  // Merge updated tasks into the main task list
  // Tasks may appear in both otherTasks and updatedTasks if status hasn't been updated yet
  // Deduplicate by ID, with updatedTasks taking precedence
  const otherTasks = props.tasks.filter(t => t.status !== status)
  const merged = [...otherTasks, ...updatedTasks]
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
  gap: var(--space-5);
  overflow-x: auto;
  padding: var(--space-2) var(--space-1);
  padding-bottom: var(--space-6);
  min-height: 400px;
}

/* Column staggered animation */
.column-animate-in {
  animation: columnEnter 0.5s ease forwards;
  opacity: 0;
}

@keyframes columnEnter {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Custom scrollbar for board */
.kanban-board::-webkit-scrollbar {
  height: 10px;
}

.kanban-board::-webkit-scrollbar-track {
  background: var(--color-bg-muted);
  border-radius: var(--radius-full);
}

.kanban-board::-webkit-scrollbar-thumb {
  background: var(--color-border-hover);
  border-radius: var(--radius-full);
}

.kanban-board::-webkit-scrollbar-thumb:hover {
  background: var(--color-text-muted);
}

/* Smooth scrolling when dragging */
.kanban-board {
  scroll-behavior: smooth;
}
</style>

<template>
  <div
    class="kanban-card"
    :class="{
      'my-task': isMyTask,
      'has-assignee': task.assignee,
      'has-description': task.description
    }"
    @click="$emit('click', task)"
  >
    <!-- Task Title -->
    <div class="task-title">
      <span class="task-number">#{{ task.id }}</span>
      {{ task.title }}
    </div>

    <!-- Task Description -->
    <div v-if="task.description" class="task-description">
      {{ truncatedDescription }}
    </div>

    <!-- Task Meta -->
    <div class="task-meta">
      <!-- Spacer -->
      <div class="meta-spacer"></div>

      <!-- Assignee -->
      <div v-if="task.assignee" class="task-assignee">
        <span
          class="assignee-avatar"
          :style="{ backgroundColor: getAvatarColor(task.assignee) }"
        >
          {{ getInitials(task.assignee) }}
        </span>
      </div>
    </div>

    <!-- Hover Glow Effect -->
    <div class="card-glow"></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Task, Assignee } from '../../api/project'

const props = defineProps<{
  task: Task
  currentUserId: number | null
}>()

defineEmits<{
  click: [task: Task]
}>()

const isMyTask = computed(() => {
  return props.currentUserId !== null && props.task.assignee?.id === props.currentUserId
})

const truncatedDescription = computed(() => {
  if (!props.task.description) return ''
  return props.task.description.length > 60
    ? props.task.description.substring(0, 60) + '...'
    : props.task.description
})

const avatarColors = [
  '#5b5fc7', '#8b5cf6', '#ec4899', '#f97316',
  '#14b8a6', '#06b6d4', '#3b82f6', '#84cc16'
]

function getAvatarColor(assignee: Assignee | undefined): string {
  if (!assignee) return avatarColors[0]
  const name = assignee.name || assignee.display_name || '?'
  const charCode = name.charCodeAt(0)
  return avatarColors[charCode % avatarColors.length]
}

function getInitials(assignee: Assignee | undefined): string {
  if (!assignee) return '?'
  const name = assignee.name || assignee.display_name || '?'
  return name.charAt(0).toUpperCase()
}
</script>

<style scoped>
.kanban-card {
  position: relative;
  background: var(--color-bg-surface);
  border-radius: var(--radius-lg);
  padding: var(--space-4);
  box-shadow: var(--shadow-sm);
  cursor: pointer;
  transition: all var(--transition-base);
  border: 1px solid transparent;
  overflow: hidden;
}

.kanban-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-lg);
  border-color: var(--color-border);
}

.kanban-card:hover .card-glow {
  opacity: 1;
}

.kanban-card.my-task {
  border: 1px solid var(--color-primary-light);
  background: linear-gradient(135deg, rgba(91, 95, 199, 0.03) 0%, var(--color-bg-surface) 100%);
}

/* Card Glow */
.card-glow {
  position: absolute;
  inset: 0;
  background: radial-gradient(
    circle at 50% 0%,
    rgba(91, 95, 199, 0.08) 0%,
    transparent 70%
  );
  opacity: 0;
  transition: opacity var(--transition-base);
  pointer-events: none;
}

/* Task Title */
.task-title {
  font-family: var(--font-display);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  margin-bottom: var(--space-2);
  line-height: var(--line-height-tight);
  padding-right: var(--space-4);
  display: flex;
  align-items: baseline;
  gap: var(--space-1);
}

.task-number {
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-muted);
  flex-shrink: 0;
}

/* Task Description */
.task-description {
  font-size: var(--font-size-xs);
  color: var(--color-text-muted);
  line-height: var(--line-height-normal);
  margin-bottom: var(--space-3);
}

/* Task Meta */
.task-meta {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  margin-top: auto;
}

.meta-spacer {
  flex: 1;
}

/* Assignee */
.task-assignee {
  display: flex;
  align-items: center;
}

.assignee-avatar {
  width: 26px;
  height: 26px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-semibold);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.15);
  border: 2px solid var(--color-bg-surface);
}

/* Animations */
@keyframes cardEnter {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.kanban-card {
  animation: cardEnter 0.3s ease forwards;
}
</style>

<template>
  <div
    class="kanban-card"
    :class="{ 'my-task': isMyTask }"
    @click="$emit('click', task)"
  >
    <div class="task-title">{{ task.title }}</div>
    <div v-if="task.description" class="task-description">
      {{ task.description.substring(0, 50) }}{{ task.description.length > 50 ? '...' : '' }}
    </div>
    <div v-if="task.assignee" class="task-assignee">
      <span
        class="assignee-avatar"
        :style="{ backgroundColor: task.assignee.avatar || task.assignee.avatar_url || '#667eea' }"
      >
        {{ (task.assignee.name || task.assignee.display_name || '?').charAt(0) }}
      </span>
      <span class="assignee-name">{{ task.assignee.name || task.assignee.display_name }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Task } from '../../api/project'

const props = defineProps<{
  task: Task
  currentUserId: number | null
}>()

defineEmits<{
  click: [task: Task]
}>()

const isMyTask = computed(() => {
  return props.currentUserId !== null && props.task.assignee.id === props.currentUserId
})

import { computed } from 'vue'
</script>

<style scoped>
.kanban-card {
  background: white;
  border-radius: 8px;
  padding: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.2s;
}

.kanban-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.kanban-card.my-task {
  border: 2px solid #667eea;
  background: #f8f9ff;
}

.task-title {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 6px;
  line-height: 1.4;
}

.task-description {
  font-size: 12px;
  color: #999;
  margin-bottom: 8px;
  line-height: 1.4;
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
  font-weight: 600;
}

.assignee-name {
  font-size: 12px;
  color: #666;
}
</style>
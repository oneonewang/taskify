<template>
  <div class="kanban-column">
    <div class="column-header">
      <span class="column-title">{{ name }}</span>
      <span class="column-count">{{ tasks.length }}</span>
    </div>
    <div class="column-content">
      <KanbanCard
        v-for="task in tasks"
        :key="task.id"
        :task="task"
        :current-user-id="currentUserId"
        @click="$emit('taskClick', task)"
      />
      <div v-if="tasks.length === 0" class="empty-column">
        暂无任务
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Task } from '../../api/project'
import KanbanCard from './KanbanCard.vue'

defineProps<{
  name: string
  status: string
  tasks: Task[]
  currentUserId: number | null
}>()

defineEmits<{
  taskClick: [task: Task]
}>()
</script>

<style scoped>
.kanban-column {
  flex: 0 0 280px;
  background: #eef1f5;
  border-radius: 12px;
  padding: 12px;
  display: flex;
  flex-direction: column;
  max-height: calc(100vh - 140px);
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
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-height: 100px;
}

.empty-column {
  text-align: center;
  color: #999;
  font-size: 13px;
  padding: 20px 0;
}
</style>
<template>
  <div class="kanban-column" :class="`status-${status}`">
    <!-- Column Header -->
    <div class="column-header">
      <div class="header-left">
        <div class="status-dot"></div>
        <span class="column-title">{{ name }}</span>
      </div>
      <span class="column-count">{{ tasks.length }}</span>
    </div>

    <!-- Column Content -->
    <div class="column-content">
      <VueDraggable
        v-model="localTasks"
        class="drag-area"
        group="kanban"
        :animation="200"
        ghost-class="ghost-card"
        drag-class="dragging-card"
        @end="onDragEnd"
        @add="onAdd"
      >
        <KanbanCard
          v-for="(task, index) in localTasks"
          :key="task.id"
          :task="task"
          :current-user-id="currentUserId"
          :style="{ animationDelay: `${index * 0.05}s` }"
          @click="$emit('taskClick', task)"
        />
      </VueDraggable>

      <div v-if="localTasks.length === 0" class="empty-column">
        <div class="empty-icon">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
            <line x1="12" y1="8" x2="12" y2="16"/>
            <line x1="8" y1="12" x2="16" y2="12"/>
          </svg>
        </div>
        <span>暂无任务</span>
      </div>
    </div>

    <!-- Column Footer -->
    <div v-if="canEdit && status === 'todo'" class="column-footer">
      <el-button
        type="primary"
        text
        class="add-task-btn"
        @click="showTaskForm = true"
      >
        <el-icon><Plus /></el-icon>
        添加任务
      </el-button>
    </div>

    <!-- Task Form Dialog -->
    <TaskForm
      v-model:visible="showTaskForm"
      :project-id="projectId"
      @success="onTaskCreated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'
import { Plus } from '@element-plus/icons-vue'
import type { Task } from '../../api/project'
import KanbanCard from './KanbanCard.vue'
import TaskForm from '../task/TaskForm.vue'
import { usePermission } from '../../composables/usePermission'

const props = defineProps<{
  name: string
  status: string
  tasks: Task[]
  currentUserId: number | null
  projectId: number
}>()

const emit = defineEmits<{
  taskClick: [task: Task]
  taskMoved: [taskId: number, newStatus: string, newPosition: number]
  tasksUpdated: [tasks: Task[]]
}>()

const permission = usePermission()

// Local task list for drag-and-drop
const localTasks = ref<Task[]>([...props.tasks])
const showTaskForm = ref(false)
const canEdit = ref(false)

async function checkEditPermission() {
  if (!props.projectId || props.projectId === 0) {
    canEdit.value = false
    return
  }
  canEdit.value = await permission.canEditTask(props.projectId)
}

onMounted(() => {
  checkEditPermission()
})

watch(() => props.projectId, () => {
  checkEditPermission()
})

// Watch for external task changes
watch(() => props.tasks, (newTasks) => {
  localTasks.value = newTasks.map(task => {
    const existing = localTasks.value.find(t => t.id === task.id)
    if (existing) {
      return { ...task }
    }
    return task
  })
}, { deep: true })

// Drag end event
function onDragEnd() {
  emit('tasksUpdated', localTasks.value)
}

// Task added to new column (cross-column drag)
function onAdd(evt: any) {
  const task = localTasks.value[evt.newIndex]
  if (task && task.status !== props.status) {
    task.status = props.status
    emit('taskMoved', task.id, props.status, evt.newIndex)
  }
  emit('tasksUpdated', localTasks.value)
}

// Task created successfully
function onTaskCreated(task: Task) {
  emit('tasksUpdated', [...localTasks.value, task])
}
</script>

<style scoped>
.kanban-column {
  flex: 0 0 300px;
  background: var(--color-bg-muted);
  border-radius: var(--radius-xl);
  padding: var(--space-4);
  display: flex;
  flex-direction: column;
  max-height: calc(100vh - 160px);
  border: 1px solid var(--color-border);
  transition: all var(--transition-base);
}

.kanban-column:hover {
  border-color: var(--color-border-hover);
}

/* Status-specific accent colors */
.kanban-column.status-todo {
  --column-accent: var(--color-info);
}

.kanban-column.status-in_progress {
  --column-accent: var(--color-primary);
}

.kanban-column.status-done {
  --column-accent: var(--color-success);
}

/* Column Header */
.column-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-4);
  padding: 0 var(--space-1);
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--column-accent);
  box-shadow: 0 0 8px var(--column-accent);
}

.column-title {
  font-family: var(--font-display);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
}

.column-count {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 24px;
  height: 24px;
  padding: 0 var(--space-2);
  background: var(--color-bg-surface);
  color: var(--color-text-secondary);
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
  border-radius: var(--radius-full);
  box-shadow: var(--shadow-sm);
}

/* Column Content */
.column-content {
  flex: 1;
  overflow-y: auto;
  min-height: 120px;
  padding-right: var(--space-1);
}

.drag-area {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  min-height: 80px;
}

/* Empty State */
.empty-column {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-8) var(--space-4);
  color: var(--color-text-muted);
  text-align: center;
}

.empty-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background: var(--color-bg-surface);
  border-radius: var(--radius-lg);
  margin-bottom: var(--space-3);
  color: var(--color-text-muted);
  box-shadow: var(--shadow-sm);
}

.empty-column span {
  font-size: var(--font-size-sm);
}

/* Column Footer */
.column-footer {
  margin-top: var(--space-3);
  padding-top: var(--space-3);
  border-top: 1px dashed var(--color-border);
}

.add-task-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  font-weight: var(--font-weight-medium);
  color: var(--color-primary);
  background: var(--color-bg-surface);
  border-radius: var(--radius-md);
  transition: all var(--transition-base);
}

.add-task-btn:hover {
  background: var(--color-primary);
  color: white;
}

/* Drag Styles */
.ghost-card {
  opacity: 0.6;
  background: var(--color-bg-muted);
  border: 2px dashed var(--color-primary-light);
  border-radius: var(--radius-lg);
}

.dragging-card {
  transform: rotate(3deg) scale(1.02);
  box-shadow: var(--shadow-xl);
  cursor: grabbing;
}

/* Responsive */
@media (max-width: 768px) {
  .kanban-column {
    flex: 0 0 260px;
  }
}
</style>

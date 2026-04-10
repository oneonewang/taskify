<template>
  <div class="kanban-column">
    <div class="column-header">
      <span class="column-title">{{ name }}</span>
      <span class="column-count">{{ tasks.length }}</span>
    </div>

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
          v-for="task in localTasks"
          :key="task.id"
          :task="task"
          :current-user-id="currentUserId"
          @click="$emit('taskClick', task)"
        />
      </VueDraggable>

      <div v-if="localTasks.length === 0" class="empty-column">
        暂无任务
      </div>
    </div>

    <!-- 创建任务按钮 -->
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

    <!-- 任务创建表单 -->
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

// 本地任务列表（用于拖拽）
const localTasks = ref<Task[]>([...props.tasks])
const showTaskForm = ref(false)
const canEdit = ref(false)

async function checkEditPermission() {
  // 只有有效的项目ID才检查权限
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

// 监听 tasks 变化
watch(() => props.tasks, (newTasks) => {
  // 保留拖拽过程中的位置信息
  localTasks.value = newTasks.map(task => {
    const existing = localTasks.value.find(t => t.id === task.id)
    if (existing) {
      return { ...task }
    }
    return task
  })
}, { deep: true })

// 拖拽结束事件
function onDragEnd() {
  // 更新父组件的任务列表
  emit('tasksUpdated', localTasks.value)
}

// 任务添加到新列时触发（跨列拖拽）
function onAdd(evt: any) {
  const task = localTasks.value[evt.newIndex]
  if (task && task.status !== props.status) {
    task.status = props.status
    emit('taskMoved', task.id, props.status, evt.newIndex)
  }
  // 更新父组件的任务列表
  emit('tasksUpdated', localTasks.value)
}

// 任务创建成功
function onTaskCreated(task: Task) {
  emit('tasksUpdated', [...localTasks.value, task])
}
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
  min-height: 100px;
}

.drag-area {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-height: 60px;
}

.empty-column {
  text-align: center;
  color: #999;
  font-size: 13px;
  padding: 20px 0;
}

.column-footer {
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid #dfe4ea;
}

.add-task-btn {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

/* 拖拽样式 */
.ghost-card {
  opacity: 0.5;
  background: #c8e6ff;
  border: 2px dashed #667eea;
}

.dragging-card {
  transform: rotate(3deg);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
}
</style>

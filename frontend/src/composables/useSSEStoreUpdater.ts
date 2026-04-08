import { useProjectStore } from '../stores/project'
import { useSSE } from './useSSE'
import type { SSETaskCreatedEvent, SSETaskUpdatedEvent, SSETaskDeletedEvent } from './useSSE'

export function useSSEStoreUpdater() {
  const projectStore = useProjectStore()

  // SSE 事件处理器
  const sseHandlers = {
    onTaskCreated: (task: SSETaskCreatedEvent) => {
      // 新任务添加到当前项目
      if (projectStore.currentProject && task.project_id === projectStore.currentProject.id) {
        projectStore.addTask(task)
      }
    },

    onTaskUpdated: (task: SSETaskUpdatedEvent) => {
      // 更新任务
      if (projectStore.currentProject && task.project_id === projectStore.currentProject.id) {
        projectStore.updateTaskBySSE(task)
      }
    },

    onTaskDeleted: (event: SSETaskDeletedEvent) => {
      // 删除任务
      projectStore.removeTask(event.id)
    },

    onConnected: () => {
      console.log('SSE connected')
    },

    onDisconnected: () => {
      console.log('SSE disconnected')
    },

    onError: (error: Event) => {
      console.error('SSE error:', error)
    }
  }

  // 连接 SSE
  useSSE(sseHandlers)
}

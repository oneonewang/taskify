import { ref, onMounted, onUnmounted } from 'vue'
import type { Task } from '../api/project'

export interface SSETaskMovedEvent {
  id: number
  from: string
  to: string
  position: number
}

export interface SSETaskCreatedEvent extends Task {}
export interface SSETaskUpdatedEvent extends Task {}
export interface SSETaskDeletedEvent {
  id: number
}
export interface SSECommentAddedEvent {
  id: number
  content: string
  user: {
    id: number
    name: string
    avatar: string
  }
  task_id: number
  created_at: string
  updated_at: string
}

export type SSEEventHandler = {
  onTaskCreated?: (task: SSETaskCreatedEvent) => void
  onTaskUpdated?: (task: SSETaskUpdatedEvent) => void
  onTaskDeleted?: (event: SSETaskDeletedEvent) => void
  onTaskMoved?: (event: SSETaskMovedEvent) => void
  onCommentAdded?: (event: SSECommentAddedEvent) => void
  onConnected?: () => void
  onDisconnected?: () => void
  onError?: (error: Event) => void
}

export function useSSE(handlers: SSEEventHandler) {
  const connected = ref(false)
  let eventSource: EventSource | null = null

  function connect() {
    if (eventSource) {
      eventSource.close()
    }

    eventSource = new EventSource('/api/events')

    eventSource.onopen = () => {
      connected.value = true
      handlers.onConnected?.()
    }

    eventSource.onerror = (error) => {
      connected.value = false
      handlers.onError?.(error)
      // EventSource 会自动重连，这里可以做一些清理工作
      if (eventSource) {
        eventSource.close()
        eventSource = null
      }
    }

    // 监听 task_created 事件
    eventSource.addEventListener('task_created', (event) => {
      try {
        const data = JSON.parse(event.data)
        handlers.onTaskCreated?.(data)
      } catch (e) {
        console.error('Failed to parse task_created event:', e)
      }
    })

    // 监听 task_updated 事件
    eventSource.addEventListener('task_updated', (event) => {
      try {
        const data = JSON.parse(event.data)
        handlers.onTaskUpdated?.(data)
      } catch (e) {
        console.error('Failed to parse task_updated event:', e)
      }
    })

    // 监听 task_deleted 事件
    eventSource.addEventListener('task_deleted', (event) => {
      try {
        const data = JSON.parse(event.data)
        handlers.onTaskDeleted?.(data)
      } catch (e) {
        console.error('Failed to parse task_deleted event:', e)
      }
    })

    // 监听 task_moved 事件
    eventSource.addEventListener('task_moved', (event) => {
      try {
        const data = JSON.parse(event.data)
        handlers.onTaskMoved?.(data)
      } catch (e) {
        console.error('Failed to parse task_moved event:', e)
      }
    })

    // 监听 comment_added 事件
    eventSource.addEventListener('comment_added', (event) => {
      try {
        const data = JSON.parse(event.data)
        handlers.onCommentAdded?.(data)
      } catch (e) {
        console.error('Failed to parse comment_added event:', e)
      }
    })
  }

  function disconnect() {
    if (eventSource) {
      eventSource.close()
      eventSource = null
      connected.value = false
      handlers.onDisconnected?.()
    }
  }

  onMounted(() => {
    connect()
  })

  onUnmounted(() => {
    disconnect()
  })

  return {
    connected,
    connect,
    disconnect
  }
}

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getProject, getProjectTasks, type Task, type TaskCount } from '../api/project'

export interface ProjectState {
  id: number
  name: string
  description: string
  created_at: string
}

export const useProjectStore = defineStore('project', () => {
  const currentProject = ref<ProjectState | null>(null)
  const tasks = ref<Task[]>([])
  const taskCounts = ref<TaskCount>({ todo: 0, in_progress: 0, review: 0, done: 0 })
  const loading = ref(false)
  const error = ref('')

  const columns = [
    { status: 'todo', name: '待办' },
    { status: 'in_progress', name: '进行中' },
    { status: 'review', name: '审核中' },
    { status: 'done', name: '已完成' }
  ]

  function getTasksByStatus(status: string): Task[] {
    return tasks.value.filter(t => t.status === status)
  }

  async function loadProject(projectId: number) {
    loading.value = true
    error.value = ''
    try {
      // 并行加载项目详情和任务列表
      const [projectRes, tasksRes] = await Promise.all([
        getProject(projectId),
        getProjectTasks(projectId)
      ])

      if (projectRes.code === 0) {
        currentProject.value = {
          id: projectRes.data.id,
          name: projectRes.data.name,
          description: projectRes.data.description,
          created_at: projectRes.data.created_at
        }
        taskCounts.value = projectRes.data.task_counts || { todo: 0, in_progress: 0, review: 0, done: 0 }
      } else {
        error.value = projectRes.message || '加载项目失败'
        return
      }

      if (tasksRes.code === 0) {
        tasks.value = tasksRes.data
        // 更新任务数量
        taskCounts.value = {
          todo: tasksRes.data.filter(t => t.status === 'todo').length,
          in_progress: tasksRes.data.filter(t => t.status === 'in_progress').length,
          review: tasksRes.data.filter(t => t.status === 'review').length,
          done: tasksRes.data.filter(t => t.status === 'done').length
        }
      }
    } catch (e: any) {
      error.value = e.message || '加载项目失败'
    } finally {
      loading.value = false
    }
  }

  function setTasks(newTasks: Task[]) {
    tasks.value = newTasks
    // 重新计算任务数量
    taskCounts.value = {
      todo: newTasks.filter(t => t.status === 'todo').length,
      in_progress: newTasks.filter(t => t.status === 'in_progress').length,
      review: newTasks.filter(t => t.status === 'review').length,
      done: newTasks.filter(t => t.status === 'done').length
    }
  }

  function addTask(task: Task) {
    // 检查任务是否已存在（防止重复添加）
    const exists = tasks.value.some(t => t.id === task.id)
    if (!exists) {
      tasks.value.push(task)
      // 更新计数
      if (task.status === 'todo') taskCounts.value.todo++
      else if (task.status === 'in_progress') taskCounts.value.in_progress++
      else if (task.status === 'review') taskCounts.value.review++
      else if (task.status === 'done') taskCounts.value.done++
    }
  }

  function updateTask(taskId: number, updates: Partial<Task>) {
    const index = tasks.value.findIndex(t => t.id === taskId)
    if (index !== -1) {
      const oldTask = tasks.value[index]
      tasks.value[index] = { ...oldTask, ...updates }
    }
  }

  function updateTaskBySSE(task: Task) {
    const index = tasks.value.findIndex(t => t.id === task.id)
    if (index !== -1) {
      tasks.value[index] = task
    }
  }

  function removeTask(taskId: number) {
    const task = tasks.value.find(t => t.id === taskId)
    if (task) {
      // 更新计数
      if (task.status === 'todo') taskCounts.value.todo--
      else if (task.status === 'in_progress') taskCounts.value.in_progress--
      else if (task.status === 'review') taskCounts.value.review--
      else if (task.status === 'done') taskCounts.value.done--
      tasks.value = tasks.value.filter(t => t.id !== taskId)
    }
  }

  return {
    currentProject,
    tasks,
    taskCounts,
    loading,
    error,
    columns,
    getTasksByStatus,
    loadProject,
    setTasks,
    addTask,
    updateTask,
    updateTaskBySSE,
    removeTask
  }
})
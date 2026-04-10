import { computed } from 'vue'
import { useAuthStore } from '../stores/auth'

// 权限检查钩子
export function usePermission() {
  const authStore = useAuthStore()

  // 检查用户是否有指定权限
  function hasPermission(resource: string, action: string): boolean {
    // admin@example.com 拥有所有权限
    if (authStore.user?.email === 'admin@example.com') {
      return true
    }
    // TODO: 从用户角色和权限映射中检查
    // 目前基于固定角色判断
    return false
  }

  // 检查是否是项目所有者
  function isProjectOwner(projectId: number): boolean {
    // TODO: 从项目成员信息中检查
    return false
  }

  // 检查是否是项目成员
  function isProjectMember(projectId: number): boolean {
    // TODO: 从项目成员信息中检查
    return false
  }

  // 检查是否可以创建任务
  const canCreateTask = computed(() => {
    // TODO: 根据用户角色判断
    // 临时：所有登录用户可以创建
    return authStore.isLoggedIn
  })

  // 检查是否可以编辑任务
  const canEditTask = computed(() => {
    return authStore.isLoggedIn
  })

  // 检查是否可以删除任务
  const canDeleteTask = computed(() => {
    // TODO: 只有任务创建者或项目所有者可以删除
    return authStore.isLoggedIn
  })

  // 检查是否可以管理项目成员
  const canManageProjectMembers = computed(() => {
    // TODO: 只有项目所有者可以管理成员
    return authStore.isLoggedIn
  })

  // 检查是否可以删除项目
  const canDeleteProject = computed(() => {
    // TODO: 只有项目所有者可以删除项目
    return authStore.isLoggedIn
  })

  return {
    hasPermission,
    isProjectOwner,
    isProjectMember,
    canCreateTask,
    canEditTask,
    canDeleteTask,
    canManageProjectMembers,
    canDeleteProject
  }
}

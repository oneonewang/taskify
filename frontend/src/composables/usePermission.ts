import { useAuthStore } from '../stores/auth'
import { getMyProjectMembership } from '../api/membership'

// 缓存用户成员资格
const membershipCache = new Map<number, { role_name: string; role_display_name: string } | null>()
const cacheExpiry = new Map<number, number>()
const CACHE_TTL = 30000 // 30秒缓存

// 权限检查钩子
export function usePermission() {
  const authStore = useAuthStore()

  // 获取用户在项目中的成员资格
  async function fetchMyMembership(projectId: number): Promise<{ role_name: string; role_display_name: string } | null> {
    // 无效的项目ID直接返回null
    if (!projectId || projectId === 0) {
      return null
    }

    const now = Date.now()
    const cached = membershipCache.get(projectId)

    if (cached !== undefined && cacheExpiry.has(projectId) && now - cacheExpiry.get(projectId)! < CACHE_TTL) {
      return cached
    }

    try {
      const res = await getMyProjectMembership(projectId)
      if (res.code === 0 && res.data) {
        membershipCache.set(projectId, res.data)
        cacheExpiry.set(projectId, now)
        return res.data
      }
    } catch (e) {
      console.error('获取成员资格失败', e)
    }
    membershipCache.set(projectId, null)
    cacheExpiry.set(projectId, now)
    return null
  }

  // 清除成员资格缓存
  function clearMembershipCache(projectId?: number) {
    if (projectId) {
      membershipCache.delete(projectId)
      cacheExpiry.delete(projectId)
    } else {
      membershipCache.clear()
      cacheExpiry.clear()
    }
  }

  // 检查用户是否有指定权限
  function hasPermission(_resource: string, _action: string): boolean {
    // admin@example.com 拥有所有权限
    if (authStore.user?.email === 'admin@example.com') {
      return true
    }
    // TODO: 从用户角色和权限映射中检查
    return false
  }

  // 检查是否是项目所有者
  async function isProjectOwner(projectId: number): Promise<boolean> {
    if (authStore.user?.email === 'admin@example.com') {
      return true
    }
    const membership = await fetchMyMembership(projectId)
    return membership?.role_name === 'owner'
  }

  // 检查是否是项目成员
  async function isProjectMember(projectId: number): Promise<boolean> {
    if (authStore.user?.email === 'admin@example.com') {
      return true
    }
    const membership = await fetchMyMembership(projectId)
    return membership !== null
  }

  // 检查是否可以管理项目成员（只有所有者可以）
  async function canManageProjectMembers(projectId: number): Promise<boolean> {
    if (authStore.user?.email === 'admin@example.com') {
      return true
    }
    const membership = await fetchMyMembership(projectId)
    return membership?.role_name === 'owner'
  }

  // 检查是否可以删除项目（只有所有者可以）
  async function canDeleteProject(projectId: number): Promise<boolean> {
    if (authStore.user?.email === 'admin@example.com') {
      return true
    }
    const membership = await fetchMyMembership(projectId)
    return membership?.role_name === 'owner'
  }

  // 检查是否可以编辑任务
  async function canEditTask(projectId: number): Promise<boolean> {
    if (authStore.user?.email === 'admin@example.com') {
      return true
    }
    const membership = await fetchMyMembership(projectId)
    return membership !== null // 成员都可以编辑任务
  }

  // 检查是否可以删除任务（任务创建者或项目所有者）
  async function canDeleteTask(projectId: number, taskCreatorId: number): Promise<boolean> {
    if (authStore.user?.email === 'admin@example.com') {
      return true
    }
    // 任务创建者可以删除自己的任务
    if (authStore.user?.id === taskCreatorId) {
      return true
    }
    // 所有者可以删除任何任务
    const membership = await fetchMyMembership(projectId)
    return membership?.role_name === 'owner'
  }

  return {
    hasPermission,
    isProjectOwner,
    isProjectMember,
    canManageProjectMembers,
    canDeleteProject,
    canEditTask,
    canDeleteTask,
    clearMembershipCache
  }
}
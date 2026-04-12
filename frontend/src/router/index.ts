import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import ProjectKanban from '../views/ProjectKanban.vue'
import ProjectList from '../views/projects/ProjectList.vue'
import ProjectSettings from '../views/projects/ProjectSettings.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import Profile from '../views/Profile.vue'
import TokenSettings from '../views/TokenSettings.vue'
import AdminLayout from '../views/admin/Layout.vue'
import AdminUsers from '../views/admin/Users.vue'
import AdminRoles from '../views/admin/Roles.vue'
import AdminAuditLogs from '../views/admin/AuditLogs.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: () => {
        const authStore = useAuthStore()
        return authStore.isLoggedIn ? '/projects' : '/login'
      }
    },
    {
      path: '/login',
      name: 'Login',
      component: Login,
      meta: { guestOnly: true }
    },
    {
      path: '/register',
      name: 'Register',
      component: Register,
      meta: { guestOnly: true }
    },
    {
      path: '/profile',
      name: 'Profile',
      component: Profile,
      meta: { requiresAuth: true }
    },
    {
      path: '/profile/tokens',
      name: 'TokenSettings',
      component: TokenSettings,
      meta: { requiresAuth: true }
    },
    {
      path: '/projects',
      name: 'ProjectList',
      component: ProjectList,
      meta: { requiresAuth: true }
    },
    {
      path: '/projects/new',
      name: 'ProjectNew',
      component: ProjectList,
      meta: { requiresAuth: true },
      props: { openCreateDialog: true }
    },
    {
      path: '/projects/:id',
      name: 'ProjectKanban',
      component: ProjectKanban,
      meta: { requiresAuth: true }
    },
    {
      path: '/projects/:id/tasks/:taskId',
      name: 'TaskDetailPage',
      component: ProjectKanban,
      meta: { requiresAuth: true }
    },
    {
      path: '/projects/:id/settings',
      name: 'ProjectSettings',
      component: ProjectSettings,
      meta: { requiresAuth: true }
    },
    {
      path: '/admin',
      component: AdminLayout,
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: 'users',
          name: 'AdminUsers',
          component: AdminUsers
        },
        {
          path: 'roles',
          name: 'AdminRoles',
          component: AdminRoles
        },
        {
          path: 'audit-logs',
          name: 'AdminAuditLogs',
          component: AdminAuditLogs
        }
      ]
    }
  ]
})

// 路由守卫
router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()

  // 如果需要认证且用户未登录
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    // 尝试初始化用户状态
    await authStore.init()

    if (!authStore.isLoggedIn) {
      next({ name: 'Login', query: { redirect: to.fullPath } })
      return
    }
  }

  // 如果需要管理员权限
  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    next({ name: 'ProjectKanban' })
    return
  }

  // 如果是访客专用页面（如登录、注册）且已登录
  if (to.meta.guestOnly && authStore.isLoggedIn) {
    next({ name: 'ProjectKanban' })
    return
  }

  next()
})

export default router
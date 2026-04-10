import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import UserSelect from '../views/UserSelect.vue'
import ProjectKanban from '../views/ProjectKanban.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import AdminLayout from '../views/admin/Layout.vue'
import AdminUsers from '../views/admin/Users.vue'
import AdminRoles from '../views/admin/Roles.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'UserSelect',
      component: UserSelect
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
      path: '/projects',
      name: 'ProjectKanban',
      component: ProjectKanban,
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
        }
      ]
    }
  ]
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
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
    // 临时：如果是 admin@example.com 则允许
    if (authStore.user?.email !== 'admin@example.com') {
      next({ name: 'ProjectKanban' })
      return
    }
  }

  // 如果是访客专用页面（如登录、注册）且已登录
  if (to.meta.guestOnly && authStore.isLoggedIn) {
    next({ name: 'ProjectKanban' })
    return
  }

  next()
})

export default router
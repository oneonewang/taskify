import { createRouter, createWebHistory } from 'vue-router'
import UserSelect from '../views/UserSelect.vue'
import ProjectKanban from '../views/ProjectKanban.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'UserSelect',
      component: UserSelect
    },
    {
      path: '/projects',
      name: 'ProjectKanban',
      component: ProjectKanban
    }
  ]
})

export default router

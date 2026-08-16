import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('../views/Home.vue')
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/trainer',
    name: 'trainer',
    component: () => import('../views/Trainer.vue')
  },
  {
    path: '/daily',
    name: 'daily',
    component: () => import('../views/Daily.vue')
  },
  {
    path: '/notebooks',
    name: 'notebooks',
    component: () => import('../views/Notebooks.vue')
  },
  {
    path: '/notebook/:id',
    name: 'notebook-read',
    component: () => import('../views/NotebookRead.vue')
  },
  {
    path: '/notebook/:id/edit',
    name: 'notebook-edit',
    component: () => import('../views/NotebookWrite.vue')
  },
  {
    path: '/community',
    name: 'community',
    component: () => import('../views/Community.vue')
  },
  {
    path: '/profile',
    name: 'profile',
    component: () => import('../views/Profile.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
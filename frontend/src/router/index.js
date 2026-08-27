import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/login', name: 'login', component: () => import('../views/Login.vue') },
  { path: '/register', name: 'register', component: () => import('../views/Register.vue') },
  { path: '/', name: 'home', component: () => import('../views/Home.vue'), meta: { requiresAuth: false } },
  { path: '/trainer', name: 'trainer', component: () => import('../views/Trainer.vue'), meta: { requiresAuth: true } },
  { path: '/daily', name: 'daily', component: () => import('../views/Daily.vue'), meta: { requiresAuth: true } },
  { path: '/notebooks', name: 'notebooks', component: () => import('../views/Notebooks.vue'), meta: { requiresAuth: true } },
  { path: '/notebook/:id', name: 'notebook-read', component: () => import('../views/NotebookRead.vue'), meta: { requiresAuth: false } },
  { path: '/notebook/:id/edit', name: 'notebook-edit', component: () => import('../views/NotebookWrite.vue'), meta: { requiresAuth: true } },
  { path: '/community', name: 'community', component: () => import('../views/Community.vue'), meta: { requiresAuth: false } },
  { path: '/profile', name: 'profile', component: () => import('../views/Profile.vue'), meta: { requiresAuth: true } },
  { path: '/admin', name: 'admin', component: () => import('../views/Admin.vue'), meta: { requiresAuth: true } },
  { path: '/rubium_tech', name: 'rubium_tech', component: () => import('../views/Rubium_tech.vue'), meta: { requiresAuth: false } },
  { path: '/courses', name: 'courses', component: () => import('../views/Courses.vue'), meta: { requiresAuth: false } },
  { path: '/verify-email', name: 'VerifyEmail', component: () => import('../views/VerifyEmail.vue') }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to) => {
  const token = localStorage.getItem('supabase_token')
  
  if (to.meta.requiresAuth && !token) {
    return '/login'
  }
})

export default router
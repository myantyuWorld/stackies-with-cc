import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { authGuard } from './authGuard'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { requiresAuth: true }
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/screens/auth/LoginScreen/LoginScreen.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/auth/callback',
      name: 'auth-callback',
      component: () => import('@/screens/auth/CallbackScreen/CallbackScreen.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/sample',
      name: 'sample',
      component: () => import('../components/SamplePage.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/AboutView.vue'),
      meta: { requiresAuth: false }
    },
  ],
})

// 認証ガードを追加
router.beforeEach(authGuard)

export default router

import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router'
import { useAuthStore } from '@/features/auth/model/authStore'

export const authGuard = async (
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext
) => {
  const authStore = useAuthStore()
  
  // 認証状態を初期化（必要に応じて）
  if (!authStore.accessToken) {
    await authStore.initializeAuth()
  }

  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const isAuthenticated = authStore.isAuthenticated

  if (requiresAuth && !isAuthenticated) {
    // 認証が必要なページで未認証の場合、ログイン画面にリダイレクト
    next({ name: 'login', query: { redirect: to.fullPath } })
  } else if (to.name === 'login' && isAuthenticated) {
    // 既に認証済みでログイン画面にアクセスした場合、ホーム画面にリダイレクト
    const redirectPath = (to.query.redirect as string) || '/'
    next(redirectPath)
  } else {
    next()
  }
}
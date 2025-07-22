import { computed } from 'vue'
import { useAuthStore } from '../model/authStore'
import type { LoginRequest } from '@/shared/auth/types'

export function useAuth() {
  const authStore = useAuthStore()

  const user = computed(() => authStore.user)
  const accessToken = computed(() => authStore.accessToken)
  const refreshToken = computed(() => authStore.refreshToken)
  const isLoading = computed(() => authStore.isLoading)
  const error = computed(() => authStore.error)
  const isAuthenticated = computed(() => authStore.isAuthenticated)

  const loginWithGoogle = (request: LoginRequest) => {
    authStore.login(request)
  }

  const logout = async () => {
    await authStore.logout()
  }

  const initialize = async () => {
    await authStore.initializeAuth()
  }

  const refreshAuth = async () => {
    await authStore.refreshAuthToken()
  }

  return {
    user,
    accessToken,
    refreshToken,
    isLoading,
    error,
    isAuthenticated,
    loginWithGoogle,
    logout,
    initialize,
    refreshAuth
  }
}
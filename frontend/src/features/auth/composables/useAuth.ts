import { computed } from 'vue'
import { useAuthStore } from '../model/authStore'
import { AUTH_CONSTANTS } from '@/shared/constants/auth'

export function useAuth() {
  const authStore = useAuthStore()

  const user = computed(() => authStore.user)
  const accessToken = computed(() => authStore.accessToken)
  const refreshToken = computed(() => authStore.refreshToken)
  const isLoading = computed(() => authStore.isLoading)
  const error = computed(() => authStore.error)
  const isAuthenticated = computed(() => authStore.isAuthenticated)

  const loginWithGoogle = () => {
    if (!(window as any).google) {
      console.error('Google Identity Services library not loaded')
      return
    }

    const client = (window as any).google.accounts.oauth2.initCodeClient({
      client_id: import.meta.env.VUE_APP_GOOGLE_CLIENT_ID,
      scope: AUTH_CONSTANTS.GOOGLE_OAUTH.SCOPE,
      callback: (response: any) => {
        if (response.error) {
          console.error('Google login error:', response)
          return
        }

        authStore.login({
          code: response.code,
          redirectUri: window.location.origin + '/auth/callback'
        })
      },
      ux_mode: AUTH_CONSTANTS.GOOGLE_OAUTH.UX_MODE,
    })

    client.requestCode()
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
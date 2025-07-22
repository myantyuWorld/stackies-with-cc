import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authClient } from '@/shared/api/authClient'
import { AUTH_CONSTANTS } from '@/shared/constants/auth'
import type { User, LoginRequest } from '@/shared/auth/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!accessToken.value && !!user.value)

  const login = async (request: LoginRequest) => {
    try {
      isLoading.value = true
      error.value = null

      const response = await authClient.login(request)
      
      user.value = response.user
      accessToken.value = response.accessToken
      refreshToken.value = response.refreshToken

      // Save tokens to localStorage
      localStorage.setItem(AUTH_CONSTANTS.STORAGE_KEYS.ACCESS_TOKEN, response.accessToken)
      localStorage.setItem(AUTH_CONSTANTS.STORAGE_KEYS.REFRESH_TOKEN, response.refreshToken)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Login failed'
      user.value = null
      accessToken.value = null
      refreshToken.value = null
    } finally {
      isLoading.value = false
    }
  }

  const logout = async () => {
    try {
      if (accessToken.value) {
        await authClient.logout(accessToken.value)
      }
    } catch (err) {
      console.error('Logout error:', err)
    } finally {
      // Clear state regardless of API call success
      user.value = null
      accessToken.value = null
      refreshToken.value = null
      error.value = null

      // Clear localStorage
      localStorage.removeItem(AUTH_CONSTANTS.STORAGE_KEYS.ACCESS_TOKEN)
      localStorage.removeItem(AUTH_CONSTANTS.STORAGE_KEYS.REFRESH_TOKEN)
    }
  }

  const initializeAuth = async () => {
    try {
      isLoading.value = true
      
      const storedAccessToken = localStorage.getItem(AUTH_CONSTANTS.STORAGE_KEYS.ACCESS_TOKEN)
      const storedRefreshToken = localStorage.getItem(AUTH_CONSTANTS.STORAGE_KEYS.REFRESH_TOKEN)

      if (storedAccessToken) {
        try {
          const userData = await authClient.getUser(storedAccessToken)
          
          user.value = userData
          accessToken.value = storedAccessToken
          refreshToken.value = storedRefreshToken
        } catch (err) {
          console.error('Error initializing auth:', err)
          // Invalid token, clear localStorage
          localStorage.removeItem(AUTH_CONSTANTS.STORAGE_KEYS.ACCESS_TOKEN)
          localStorage.removeItem(AUTH_CONSTANTS.STORAGE_KEYS.REFRESH_TOKEN)
          user.value = null
          accessToken.value = null
          refreshToken.value = null
        }
      }
    } finally {
      isLoading.value = false
    }
  }

  const refreshAuthToken = async () => {
    try {
      if (!refreshToken.value) {
        throw new Error('No refresh token available')
      }

      const response = await authClient.refreshToken({ refreshToken: refreshToken.value })
      
      accessToken.value = response.accessToken
      refreshToken.value = response.refreshToken

      localStorage.setItem(AUTH_CONSTANTS.STORAGE_KEYS.ACCESS_TOKEN, response.accessToken)
      localStorage.setItem(AUTH_CONSTANTS.STORAGE_KEYS.REFRESH_TOKEN, response.refreshToken)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Token refresh failed'
      await logout()
    }
  }

  return {
    user,
    accessToken,
    refreshToken,
    isLoading,
    error,
    isAuthenticated,
    login,
    logout,
    initializeAuth,
    refreshAuthToken
  }
})
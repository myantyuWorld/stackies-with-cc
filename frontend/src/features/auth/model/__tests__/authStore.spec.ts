import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '../authStore'
import type { User, LoginRequest } from '@/shared/auth/types'

// Mock the auth client
vi.mock('@/shared/api/authClient', () => ({
  authClient: {
    login: vi.fn(),
    refreshToken: vi.fn(),
    getUser: vi.fn(),
    logout: vi.fn()
  }
}))

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn()
}
Object.defineProperty(window, 'localStorage', {
  value: localStorageMock
})

describe('Auth Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  describe('Initial State', () => {
    it('should initialize with default state', () => {
      const store = useAuthStore()
      
      expect(store.user).toBeNull()
      expect(store.accessToken).toBeNull()
      expect(store.refreshToken).toBeNull()
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeNull()
    })

    it('should have computed isAuthenticated as false initially', () => {
      const store = useAuthStore()
      
      expect(store.isAuthenticated).toBe(false)
    })
  })

  describe('Actions', () => {
    describe('login', () => {
      it('should set loading state during login', async () => {
        const store = useAuthStore()
        const { authClient } = await import('@/shared/api/authClient')
        
        // Mock delayed response
        vi.mocked(authClient.login).mockImplementation(() => 
          new Promise(resolve => setTimeout(() => resolve({
            user: {
              id: 'test-id',
              email: 'test@example.com',
              name: 'Test User',
              picture: 'https://example.com/picture.jpg'
            },
            accessToken: 'access-token',
            refreshToken: 'refresh-token'
          }), 100))
        )

        const loginRequest: LoginRequest = {
          code: 'google-auth-code',
          state: 'google-auth-state'
        }

        const loginPromise = store.login(loginRequest)
        expect(store.isLoading).toBe(true)

        await loginPromise
        expect(store.isLoading).toBe(false)
      })

      it('should update state on successful login', async () => {
        const store = useAuthStore()
        const { authClient } = await import('@/shared/api/authClient')
        
        const mockUser: User = {
          id: 'test-id',
          email: 'test@example.com',
          name: 'Test User',
          picture: 'https://example.com/picture.jpg'
        }

        vi.mocked(authClient.login).mockResolvedValue({
          user: mockUser,
          accessToken: 'access-token',
          refreshToken: 'refresh-token'
        })

        const loginRequest: LoginRequest = {
          code: 'google-auth-code',
          state: 'google-auth-state'
        }

        await store.login(loginRequest)

        expect(store.user).toEqual(mockUser)
        expect(store.accessToken).toBe('access-token')
        expect(store.refreshToken).toBe('refresh-token')
        expect(store.error).toBeNull()
        expect(store.isAuthenticated).toBe(true)
      })

      it('should save tokens to localStorage on successful login', async () => {
        const store = useAuthStore()
        const { authClient } = await import('@/shared/api/authClient')
        
        vi.mocked(authClient.login).mockResolvedValue({
          user: {
            id: 'test-id',
            email: 'test@example.com',
            name: 'Test User',
            picture: 'https://example.com/picture.jpg'
          },
          accessToken: 'access-token',
          refreshToken: 'refresh-token'
        })

        const loginRequest: LoginRequest = {
          code: 'google-auth-code',
          state: 'google-auth-state'
        }

        await store.login(loginRequest)

        expect(localStorageMock.setItem).toHaveBeenCalledWith('accessToken', 'access-token')
        expect(localStorageMock.setItem).toHaveBeenCalledWith('refreshToken', 'refresh-token')
      })

      it('should handle login error', async () => {
        const store = useAuthStore()
        const { authClient } = await import('@/shared/api/authClient')
        
        const errorMessage = 'Login failed'
        vi.mocked(authClient.login).mockRejectedValue(new Error(errorMessage))

        const loginRequest: LoginRequest = {
          code: 'invalid-code',
          state: 'google-auth-state'
        }

        await store.login(loginRequest)

        expect(store.error).toBe(errorMessage)
        expect(store.user).toBeNull()
        expect(store.accessToken).toBeNull()
        expect(store.refreshToken).toBeNull()
        expect(store.isLoading).toBe(false)
      })
    })

    describe('logout', () => {
      it('should clear state on logout', async () => {
        const store = useAuthStore()
        const { authClient } = await import('@/shared/api/authClient')
        
        // Set initial authenticated state
        store.user = {
          id: 'test-id',
          email: 'test@example.com',
          name: 'Test User',
          picture: 'https://example.com/picture.jpg'
        }
        store.accessToken = 'access-token'
        store.refreshToken = 'refresh-token'

        vi.mocked(authClient.logout).mockResolvedValue()

        await store.logout()

        expect(store.user).toBeNull()
        expect(store.accessToken).toBeNull()
        expect(store.refreshToken).toBeNull()
        expect(store.error).toBeNull()
        expect(store.isAuthenticated).toBe(false)
      })

      it('should clear localStorage on logout', async () => {
        const store = useAuthStore()
        const { authClient } = await import('@/shared/api/authClient')
        
        store.accessToken = 'access-token'
        vi.mocked(authClient.logout).mockResolvedValue()

        await store.logout()

        expect(localStorageMock.removeItem).toHaveBeenCalledWith('accessToken')
        expect(localStorageMock.removeItem).toHaveBeenCalledWith('refreshToken')
      })
    })

    describe('initializeAuth', () => {
      it('should load tokens from localStorage', async () => {
        const store = useAuthStore()
        const { authClient } = await import('@/shared/api/authClient')
        
        localStorageMock.getItem.mockImplementation((key) => {
          if (key === 'accessToken') return 'stored-access-token'
          if (key === 'refreshToken') return 'stored-refresh-token'
          return null
        })

        vi.mocked(authClient.getUser).mockResolvedValue({
          id: 'test-id',
          email: 'test@example.com',
          name: 'Test User',
          picture: 'https://example.com/picture.jpg'
        })

        await store.initializeAuth()

        expect(store.accessToken).toBe('stored-access-token')
        expect(store.refreshToken).toBe('stored-refresh-token')
        expect(store.user).toBeDefined()
      })

      it('should clear invalid tokens from localStorage', async () => {
        const store = useAuthStore()
        const { authClient } = await import('@/shared/api/authClient')
        
        localStorageMock.getItem.mockImplementation((key) => {
          if (key === 'accessToken') return 'invalid-token'
          return null
        })

        vi.mocked(authClient.getUser).mockRejectedValue(new Error('Invalid token'))

        await store.initializeAuth()

        expect(localStorageMock.removeItem).toHaveBeenCalledWith('accessToken')
        expect(localStorageMock.removeItem).toHaveBeenCalledWith('refreshToken')
        expect(store.accessToken).toBeNull()
        expect(store.refreshToken).toBeNull()
      })
    })
  })
})
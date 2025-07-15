import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

// Mock Google Auth library
const mockGoogleAuth = {
  accounts: {
    oauth2: {
      initCodeClient: vi.fn(),
      hasGrantedAllScopes: vi.fn()
    }
  }
}

Object.defineProperty(window, 'google', {
  value: mockGoogleAuth,
  writable: true
})

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

// Mock import.meta.env
Object.defineProperty(import.meta, 'env', {
  value: {
    VUE_APP_GOOGLE_CLIENT_ID: 'test-client-id'
  }
})

describe('useAuth Composable', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('should be defined', async () => {
    const { useAuth } = await import('../useAuth')
    expect(useAuth).toBeDefined()
  })

  it('should return reactive values', async () => {
    const { useAuth } = await import('../useAuth')
    const auth = useAuth()
    
    expect(auth.user).toBeDefined()
    expect(auth.accessToken).toBeDefined()
    expect(auth.isLoading).toBeDefined()
    expect(auth.error).toBeDefined()
    expect(auth.isAuthenticated).toBeDefined()
  })

  it('should have login function', async () => {
    const { useAuth } = await import('../useAuth')
    const auth = useAuth()
    
    expect(typeof auth.loginWithGoogle).toBe('function')
  })

  it('should have logout function', async () => {
    const { useAuth } = await import('../useAuth')
    const auth = useAuth()
    
    expect(typeof auth.logout).toBe('function')
  })

  it('should have initialize function', async () => {
    const { useAuth } = await import('../useAuth')
    const auth = useAuth()
    
    expect(typeof auth.initialize).toBe('function')
  })
})
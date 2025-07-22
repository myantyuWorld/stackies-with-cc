import { describe, it, expect, vi, beforeEach } from 'vitest'
import { authClient } from '../authClient'
import type { LoginRequest, LoginResponse, RefreshTokenRequest } from '../../auth/types'

// Mock fetch
global.fetch = vi.fn()

describe('Auth API Client', () => {
  beforeEach(() => {
    vi.resetAllMocks()
  })

  describe('login', () => {
    it('should make POST request to login endpoint', async () => {
      const mockResponse: LoginResponse = {
        user: {
          id: 'test-id',
          email: 'test@example.com',
          name: 'Test User',
          picture: 'https://example.com/picture.jpg'
        },
        accessToken: 'access-token',
        refreshToken: 'refresh-token'
      }

      const mockFetch = vi.mocked(fetch)
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockResponse)
      } as Response)

      const loginRequest: LoginRequest = {
        code: 'google-auth-code',
        state: 'google-auth-state'
      }

      const result = await authClient.login(loginRequest)

      expect(mockFetch).toHaveBeenCalledWith(import.meta.env.VITE_APP_API_BASE_URL + '/auth/google/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(loginRequest)
      })
      expect(result).toEqual(mockResponse)
    })

    it('should throw error on failed login', async () => {
      const mockFetch = vi.mocked(fetch)
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
        statusText: 'Unauthorized'
      } as Response)

      const loginRequest: LoginRequest = {
        code: 'invalid-code',
        state: 'google-auth-state'
      }

      await expect(authClient.login(loginRequest)).rejects.toThrow('Login failed: 401 Unauthorized')
    })
  })

  describe('refreshToken', () => {
    it('should make POST request to refresh endpoint', async () => {
      const mockResponse = {
        accessToken: 'new-access-token',
        refreshToken: 'new-refresh-token'
      }

      const mockFetch = vi.mocked(fetch)
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockResponse)
      } as Response)

      const refreshRequest: RefreshTokenRequest = {
        refreshToken: 'refresh-token'
      }

      const result = await authClient.refreshToken(refreshRequest)

      expect(mockFetch).toHaveBeenCalledWith(import.meta.env.VITE_APP_API_BASE_URL + '/auth/refresh', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(refreshRequest)
      })
      expect(result).toEqual(mockResponse)
    })

    it('should throw error on failed refresh', async () => {
      const mockFetch = vi.mocked(fetch)
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
        statusText: 'Unauthorized'
      } as Response)

      const refreshRequest: RefreshTokenRequest = {
        refreshToken: 'invalid-refresh-token'
      }

      await expect(authClient.refreshToken(refreshRequest)).rejects.toThrow('Token refresh failed: 401 Unauthorized')
    })
  })

  describe('getUser', () => {
    it('should make GET request to user endpoint with auth header', async () => {
      const mockUser = {
        id: 'test-id',
        email: 'test@example.com',
        name: 'Test User',
        picture: 'https://example.com/picture.jpg'
      }

      const mockFetch = vi.mocked(fetch)
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockUser)
      } as Response)

      const result = await authClient.getUser('access-token')

      expect(mockFetch).toHaveBeenCalledWith(import.meta.env.VITE_APP_API_BASE_URL + '/auth/me', {
        method: 'GET',
        headers: {
          'Authorization': 'Bearer access-token',
          'Content-Type': 'application/json'
        }
      })
      expect(result).toEqual(mockUser)
    })

    it('should throw error on failed user fetch', async () => {
      const mockFetch = vi.mocked(fetch)
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
        statusText: 'Unauthorized'
      } as Response)

      await expect(authClient.getUser('invalid-token')).rejects.toThrow('Get user failed: 401 Unauthorized')
    })
  })

  describe('logout', () => {
    it('should make POST request to logout endpoint with auth header', async () => {
      const mockFetch = vi.mocked(fetch)
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ success: true })
      } as Response)

      await authClient.logout('access-token')

      expect(mockFetch).toHaveBeenCalledWith(import.meta.env.VITE_APP_API_BASE_URL + '/auth/logout', {
        method: 'POST',
        headers: {
          'Authorization': 'Bearer access-token',
          'Content-Type': 'application/json'
        }
      })
    })

    it('should throw error on failed logout', async () => {
      const mockFetch = vi.mocked(fetch)
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error'
      } as Response)

      await expect(authClient.logout('access-token')).rejects.toThrow('Logout failed: 500 Internal Server Error')
    })
  })
})
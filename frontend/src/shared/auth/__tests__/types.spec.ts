import { describe, it, expect } from 'vitest'
import type { User, AuthState, LoginRequest, LoginResponse, RefreshTokenRequest } from '../types'

describe('Auth Types', () => {
  describe('User Type', () => {
    it('should have correct user properties', () => {
      const user: User = {
        id: 'test-id',
        email: 'test@example.com',
        name: 'Test User',
        picture: 'https://example.com/picture.jpg'
      }

      expect(user.id).toBe('test-id')
      expect(user.email).toBe('test@example.com')
      expect(user.name).toBe('Test User')
      expect(user.picture).toBe('https://example.com/picture.jpg')
    })
  })

  describe('AuthState Type', () => {
    it('should have correct auth state properties', () => {
      const authState: AuthState = {
        user: null,
        accessToken: null,
        refreshToken: null,
        isLoading: false,
        error: null
      }

      expect(authState.user).toBeNull()
      expect(authState.accessToken).toBeNull()
      expect(authState.refreshToken).toBeNull()
      expect(authState.isLoading).toBe(false)
      expect(authState.error).toBeNull()
    })

    it('should allow authenticated state', () => {
      const user: User = {
        id: 'test-id',
        email: 'test@example.com',
        name: 'Test User',
        picture: 'https://example.com/picture.jpg'
      }

      const authState: AuthState = {
        user,
        accessToken: 'access-token',
        refreshToken: 'refresh-token',
        isLoading: false,
        error: null
      }

      expect(authState.user).toBe(user)
      expect(authState.accessToken).toBe('access-token')
      expect(authState.refreshToken).toBe('refresh-token')
    })
  })

  describe('LoginRequest Type', () => {
    it('should have correct login request properties', () => {
      const loginRequest: LoginRequest = {
        code: 'google-auth-code',
        redirectUri: 'http://localhost:3000/auth/callback'
      }

      expect(loginRequest.code).toBe('google-auth-code')
      expect(loginRequest.redirectUri).toBe('http://localhost:3000/auth/callback')
    })
  })

  describe('LoginResponse Type', () => {
    it('should have correct login response properties', () => {
      const user: User = {
        id: 'test-id',
        email: 'test@example.com',
        name: 'Test User',
        picture: 'https://example.com/picture.jpg'
      }

      const loginResponse: LoginResponse = {
        user,
        accessToken: 'access-token',
        refreshToken: 'refresh-token'
      }

      expect(loginResponse.user).toBe(user)
      expect(loginResponse.accessToken).toBe('access-token')
      expect(loginResponse.refreshToken).toBe('refresh-token')
    })
  })

  describe('RefreshTokenRequest Type', () => {
    it('should have correct refresh token request properties', () => {
      const refreshRequest: RefreshTokenRequest = {
        refreshToken: 'refresh-token'
      }

      expect(refreshRequest.refreshToken).toBe('refresh-token')
    })
  })
})
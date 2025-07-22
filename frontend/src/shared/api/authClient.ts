import type { LoginRequest, LoginResponse, RefreshTokenRequest, RefreshTokenResponse, User } from '../auth/types'

class AuthClient {
  async login(request: LoginRequest): Promise<LoginResponse> {
    const response = await fetch(import.meta.env.VITE_APP_API_BASE_URL + '/auth/google/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(request)
    })

    console.log(response)

    if (!response.ok) {
      throw new Error(`Login failed: ${response.status} ${response.statusText}`)
    }

    return response.json()
  }

  async refreshToken(request: RefreshTokenRequest): Promise<RefreshTokenResponse> {
    const response = await fetch(import.meta.env.VITE_APP_API_BASE_URL + '/auth/refresh', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(request)
    })

    if (!response.ok) {
      throw new Error(`Token refresh failed: ${response.status} ${response.statusText}`)
    }

    return response.json()
  }

  async getUser(accessToken: string): Promise<User> {
    const response = await fetch(import.meta.env.VITE_APP_API_BASE_URL + '/auth/me', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      throw new Error(`Get user failed: ${response.status} ${response.statusText}`)
    }

    return response.json()
  }

  async logout(accessToken: string): Promise<void> {
    const response = await fetch(import.meta.env.VITE_APP_API_BASE_URL + '/auth/logout', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      throw new Error(`Logout failed: ${response.status} ${response.statusText}`)
    }
  }
}

export const authClient = new AuthClient()
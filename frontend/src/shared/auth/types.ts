export interface User {
  id: string
  email: string
  name: string
  picture: string
}

export interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  isLoading: boolean
  error: string | null
}

export interface LoginRequest {
  code: string
  redirectUri: string
}

export interface LoginResponse {
  user: User
  accessToken: string
  refreshToken: string
}

export interface RefreshTokenRequest {
  refreshToken: string
}

export interface RefreshTokenResponse {
  accessToken: string
  refreshToken: string
}

// Google OAuth関連の型定義
export interface GoogleOAuthResponse {
  code?: string
  error?: string
  error_description?: string
  state?: string
}

export interface GoogleOAuthError {
  error: string
  error_description?: string
  state?: string
}

export interface GoogleOAuthSuccess {
  code: string
  state?: string
}

// Google Identity Services関連の型定義
export interface GoogleIdentityServices {
  accounts: {
    oauth2: {
      initCodeClient: (config: GoogleOAuthClientConfig) => GoogleOAuthClient
      hasGrantedAllScopes: (token: string, scopes: string) => boolean
    }
  }
}

export interface GoogleOAuthClientConfig {
  client_id: string
  scope: string
  callback: (response: GoogleOAuthResponse) => void
  ux_mode: string
}

export interface GoogleOAuthClient {
  requestCode: () => void
}

// Window拡張の型定義
declare global {
  interface Window {
    google: GoogleIdentityServices
  }
}
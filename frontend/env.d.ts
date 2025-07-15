/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VUE_APP_GOOGLE_CLIENT_ID: string
  readonly VUE_APP_GOOGLE_REDIRECT_URI: string
  readonly VUE_APP_API_BASE_URL: string
  readonly VUE_APP_APP_NAME: string
  readonly VUE_APP_APP_VERSION: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

// Google OAuth関連の型定義
interface GoogleOAuthResponse {
  code?: string
  error?: string
  error_description?: string
  state?: string
}

interface GoogleIdentityServices {
  accounts: {
    oauth2: {
      initCodeClient: (config: {
        client_id: string
        scope: string
        callback: (response: GoogleOAuthResponse) => void
        ux_mode: string
      }) => {
        requestCode: () => void
      }
      hasGrantedAllScopes: (token: string, scopes: string) => boolean
    }
  }
}

declare global {
  interface Window {
    google: GoogleIdentityServices
  }
}

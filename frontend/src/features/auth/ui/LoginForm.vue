<template>
  <div class="login-form">
    <div v-if="error" data-testid="error-message" class="error-message">
      {{ error }}
    </div>
    
    <div v-if="isLoading" data-testid="loading-spinner" class="loading-spinner">
      ログイン中...
    </div>
    
    <button
      v-else
      data-testid="google-login-button"
      @click="handleGoogleLogin"
      class="google-login-button"
      :disabled="isLoading"
    >
      <svg class="google-icon" viewBox="0 0 24 24" width="20" height="20">
        <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
        <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
        <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
        <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
      </svg>
      Googleでログイン
    </button>
  </div>
</template>

<script setup lang="ts">
import { useAuth } from '@/features/auth/composables/useAuth'

const { isLoading, error } = useAuth()

const handleGoogleLogin = async () => {
  // まずバックエンドからGoogle認証URLを取得
  console.log(import.meta.env.VITE_APP_API_BASE_URL)

  const res = await fetch(import.meta.env.VITE_APP_API_BASE_URL + '/auth/google/url');
  console.log(res)
  
  const { auth_url } = await res.json();
  // 取得したURLにリダイレクト
  window.location.href = auth_url;
}
</script>

<style scoped>
.login-form {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 2rem;
}

.google-login-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  border: 1px solid #dadce0;
  border-radius: 4px;
  background-color: #fff;
  color: #3c4043;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.google-login-button:hover {
  box-shadow: 0 1px 2px 0 rgba(60, 64, 67, 0.3), 0 1px 3px 1px rgba(60, 64, 67, 0.15);
}

.google-login-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.google-icon {
  flex-shrink: 0;
}

.error-message {
  color: #d93025;
  background-color: #fce8e6;
  padding: 0.75rem;
  border-radius: 4px;
  border: 1px solid #f9ab00;
  font-size: 14px;
}

.loading-spinner {
  color: #1976d2;
  font-size: 14px;
  font-weight: 500;
}
</style>
<template>
  <div class="callback-screen">
    <div class="callback-container">
      <div data-testid="processing-message" class="processing-content">
        <div class="spinner">
          <div class="spinner-inner"></div>
        </div>
        <h2 class="processing-title">認証処理中...</h2>
        <p class="processing-description">
          Googleアカウントでの認証を処理しています。<br>
          しばらくお待ちください。
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '@/features/auth/composables/useAuth'

const router = useRouter()
const route = useRoute()
const { isAuthenticated, loginWithGoogle } = useAuth()

onMounted(async () => {
  // URLからGoogle OAuth認証コードを取得
  const code = route.query.code as string
  const error = route.query.error as string
  const state = route.query.state as string

  if (error) {
    // 認証エラーの場合はログイン画面にリダイレクト
    console.error('Google OAuth error:', error)
    router.push('/login')
    return
  }

  await loginWithGoogle({
    code: code,
    state: state
  })

  if (code) {
    try {
      // 認証処理は自動的に行われるため、少し待ってから状態をチェック
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      if (isAuthenticated.value) {
        // 認証成功時はホーム画面にリダイレクト
        router.push('/')
      } else {
        // 認証失敗時はログイン画面にリダイレクト
        router.push('/login')
      }
    } catch (error) {
      console.error('Authentication error:', error)
      router.push('/login')
    }
  } else {
    // 認証コードがない場合はログイン画面にリダイレクト
    router.push('/login')
  }
})
</script>

<style scoped>
.callback-screen {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 1rem;
}

.callback-container {
  width: 100%;
  max-width: 400px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  padding: 3rem 2rem;
}

.processing-content {
  text-align: center;
}

.spinner {
  width: 60px;
  height: 60px;
  margin: 0 auto 2rem;
  border: 4px solid #f3f4f6;
  border-top: 4px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.spinner-inner {
  width: 100%;
  height: 100%;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.processing-title {
  margin: 0 0 1rem;
  font-size: 1.5rem;
  font-weight: 600;
  color: #2d3748;
}

.processing-description {
  margin: 0;
  font-size: 1rem;
  color: #718096;
  line-height: 1.6;
}

@media (max-width: 640px) {
  .callback-container {
    margin: 0;
    border-radius: 0;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }
  
  .callback-screen {
    padding: 0;
  }
}
</style>
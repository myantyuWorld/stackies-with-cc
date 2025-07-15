<template>
  <div class="login-screen">
    <div class="login-container">
      <div class="login-header">
        <h1 class="app-title">Stackies</h1>
        <p class="login-subtitle">サインインしてください</p>
      </div>
      
      <div data-testid="login-form" class="login-form-container">
        <LoginForm />
      </div>
      
      <div class="login-footer">
        <p class="footer-text">
          サインインすることで、
          <a href="#" class="footer-link">利用規約</a>と
          <a href="#" class="footer-link">プライバシーポリシー</a>に同意したものとみなされます。
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import LoginForm from '@/features/auth/ui/LoginForm.vue'
</script>

<style scoped>
.login-screen {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 1rem;
}

.login-container {
  width: 100%;
  max-width: 400px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.login-header {
  padding: 2rem 2rem 1rem;
  text-align: center;
  background: white;
}

.app-title {
  margin: 0 0 0.5rem;
  font-size: 2rem;
  font-weight: 700;
  color: #2d3748;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.login-subtitle {
  margin: 0;
  font-size: 1rem;
  color: #718096;
  font-weight: 400;
}

.login-form-container {
  padding: 1rem 2rem;
}

.login-footer {
  padding: 1rem 2rem 2rem;
  text-align: center;
}

.footer-text {
  margin: 0;
  font-size: 0.875rem;
  color: #718096;
  line-height: 1.5;
}

.footer-link {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}

.footer-link:hover {
  text-decoration: underline;
}

@media (max-width: 640px) {
  .login-container {
    margin: 0;
    border-radius: 0;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }
  
  .login-screen {
    padding: 0;
  }
}
</style>
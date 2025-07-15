import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { useAuth } from '@/features/auth/composables/useAuth'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// 認証状態の初期化
const { initialize } = useAuth()
initialize()

app.mount('#app')

<template>
  <div v-if="user" data-testid="user-profile" class="user-profile">
    <div class="user-info">
      <img 
        :src="user.picture" 
        :alt="`${user.name}のアバター`"
        data-testid="user-picture"
        class="user-picture"
      />
      <div class="user-details">
        <h3 data-testid="user-name" class="user-name">{{ user.name }}</h3>
        <p data-testid="user-email" class="user-email">{{ user.email }}</p>
      </div>
    </div>
    
    <button
      data-testid="logout-button"
      @click="handleLogout"
      class="logout-button"
      :disabled="isLoading"
    >
      ログアウト
    </button>
  </div>
</template>

<script setup lang="ts">
import { useAuth } from '@/features/auth/composables/useAuth'

const { user, logout, isLoading } = useAuth()

const handleLogout = async () => {
  await logout()
}
</script>

<style scoped>
.user-profile {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  background-color: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-picture {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
}

.user-details {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.user-name {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #202124;
}

.user-email {
  margin: 0;
  font-size: 14px;
  color: #5f6368;
}

.logout-button {
  padding: 0.5rem 1rem;
  border: 1px solid #dadce0;
  border-radius: 4px;
  background-color: #fff;
  color: #3c4043;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.logout-button:hover {
  background-color: #f8f9fa;
  box-shadow: 0 1px 2px 0 rgba(60, 64, 67, 0.3);
}

.logout-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
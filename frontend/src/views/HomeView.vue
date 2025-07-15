<script setup lang="ts">
import TheWelcome from '../components/TheWelcome.vue'
import UserProfile from '@/features/auth/ui/UserProfile.vue'
import { useAuth } from '@/features/auth/composables/useAuth'

const { user } = useAuth()
</script>

<template>
  <main>
    <div class="home-header">
      <h1>Stackiesへようこそ</h1>
      <UserProfile v-if="user" />
    </div>
    <TheWelcome />
  </main>
</template>

<style scoped>
.home-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  border-bottom: 1px solid #e2e8f0;
  margin-bottom: 2rem;
}

.home-header h1 {
  margin: 0;
  color: #2d3748;
  font-size: 1.75rem;
  font-weight: 600;
}

@media (max-width: 768px) {
  .home-header {
    flex-direction: column;
    gap: 1rem;
    align-items: flex-start;
  }
}
</style>

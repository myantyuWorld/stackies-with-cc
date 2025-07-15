<script setup lang="ts">
import { ref } from 'vue'

const title = ref('サンプル画面')
const message = ref('こんにちは、Vue3 + TypeScriptです！')
const counter = ref(0)

const increment = () => {
  counter.value++
}

const reset = () => {
  counter.value = 0
}
</script>

<template>
  <div class="sample-page">
    <h1>{{ title }}</h1>
    <p class="message">{{ message }}</p>
    
    <div class="counter-section">
      <h2>カウンター</h2>
      <p class="counter-value">{{ counter }}</p>
      <div class="button-group">
        <button @click="increment" class="btn btn-primary">+1</button>
        <button @click="reset" class="btn btn-secondary">リセット</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.sample-page {
  padding: 2rem;
  max-width: 600px;
  margin: 0 auto;
  text-align: center;
}

h1 {
  color: #42b883;
  font-size: 2.5rem;
  margin-bottom: 1rem;
}

.message {
  font-size: 1.2rem;
  color: #666;
  margin-bottom: 2rem;
}

.counter-section {
  background: #f8f9fa;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

h2 {
  color: #333;
  margin-bottom: 1rem;
}

.counter-value {
  font-size: 3rem;
  font-weight: bold;
  color: #42b883;
  margin: 1rem 0;
}

.button-group {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.btn {
  padding: 0.8rem 1.5rem;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

.btn-primary {
  background-color: #42b883;
  color: white;
}

.btn-primary:hover {
  background-color: #369870;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background-color: #5a6268;
}
</style>
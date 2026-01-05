<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const localError = ref<string | null>(null)

async function handleSignUp() {
  localError.value = null

  if (!username.value || !password.value || !confirmPassword.value) {
    localError.value = 'Заполните все поля'
    return
  }

  if (password.value !== confirmPassword.value) {
    localError.value = 'Пароли не совпадают'
    return
  }

  if (password.value.length < 6) {
    localError.value = 'Пароль должен содержать минимум 6 символов'
    return
  }

  const result = await authStore.signUp(username.value, password.value)
  if (result.success) {
    router.push('/')
  } else {
    localError.value = result.error || 'Ошибка регистрации'
  }
}
</script>

<template>
  <div class="signup-container">
    <div class="signup-card">
      <div class="logo">
        <svg width="48" height="48" viewBox="0 0 48 48" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M12 12H36C38.2091 12 40 13.7909 40 16V32C40 34.2091 38.2091 36 36 36H12C9.79086 36 8 34.2091 8 32V16C8 13.7909 9.79086 12 12 12Z"/>
          <path d="M16 20H32M16 26H32M16 32H24"/>
        </svg>
      </div>
      <h1>Регистрация</h1>
      <p class="subtitle">Создайте новый аккаунт для доступа к файлам</p>
      <form @submit.prevent="handleSignUp">
        <div class="form-group">
          <label for="username">Имя пользователя</label>
          <input
            id="username"
            v-model="username"
            type="text"
            required
            autocomplete="username"
            placeholder="Введите имя пользователя"
          />
        </div>
        <div class="form-group">
          <label for="password">Пароль</label>
          <input
            id="password"
            v-model="password"
            type="password"
            required
            autocomplete="new-password"
            placeholder="Введите пароль"
          />
        </div>
        <div class="form-group">
          <label for="confirmPassword">Подтвердите пароль</label>
          <input
            id="confirmPassword"
            v-model="confirmPassword"
            type="password"
            required
            autocomplete="new-password"
            placeholder="Повторите пароль"
          />
        </div>
        <div v-if="localError" class="error-message">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="8" cy="8" r="7"/>
            <path d="M8 4V8M8 12H8.01"/>
          </svg>
          {{ localError }}
        </div>
        <button type="submit" :disabled="authStore.loading" class="submit-button">
          <span v-if="authStore.loading" class="button-spinner"></span>
          {{ authStore.loading ? 'Регистрация...' : 'Зарегистрироваться' }}
        </button>
      </form>
      <div class="login-link">
        <p>
          Уже есть аккаунт?
          <router-link to="/login">Войти</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.signup-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #fafafa;
  padding: 20px;
}

.signup-card {
  background: white;
  border-radius: 12px;
  padding: 48px;
  width: 100%;
  max-width: 420px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e5e7eb;
}

.logo {
  display: flex;
  justify-content: center;
  margin-bottom: 24px;
  color: #3b82f6;
}

h1 {
  margin: 0 0 8px 0;
  text-align: center;
  color: #111827;
  font-size: 24px;
  font-weight: 600;
}

.subtitle {
  margin: 0 0 32px 0;
  text-align: center;
  color: #6b7280;
  font-size: 14px;
}

.form-group {
  margin-bottom: 20px;
}

label {
  display: block;
  margin-bottom: 8px;
  color: #374151;
  font-weight: 500;
  font-size: 14px;
}

input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.15s;
  box-sizing: border-box;
  background: white;
  color: #111827;
}

input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

input::placeholder {
  color: #9ca3af;
}

.error-message {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #dc2626;
  margin-bottom: 16px;
  padding: 12px;
  background-color: #fef2f2;
  border-radius: 6px;
  font-size: 14px;
  border: 1px solid #fecaca;
}

.error-message svg {
  flex-shrink: 0;
}

.submit-button {
  width: 100%;
  padding: 12px;
  background: #111827;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.submit-button:hover:not(:disabled) {
  background: #1f2937;
}

.submit-button:active:not(:disabled) {
  background: #111827;
}

.submit-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.button-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.login-link {
  margin-top: 24px;
  text-align: center;
  color: #6b7280;
  font-size: 14px;
}

.login-link a {
  color: #3b82f6;
  text-decoration: none;
  font-weight: 500;
}

.login-link a:hover {
  text-decoration: underline;
}

@media (max-width: 480px) {
  .signup-card {
    padding: 32px 24px;
  }
}
</style>



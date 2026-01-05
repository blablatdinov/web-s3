import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { apiService } from '@/services/api'

const AUTH_TOKEN_KEY = 'auth_token'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem(AUTH_TOKEN_KEY))
  const loading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!token.value)

  function setToken(newToken: string | null) {
    token.value = newToken
    if (newToken) {
      localStorage.setItem(AUTH_TOKEN_KEY, newToken)
    } else {
      localStorage.removeItem(AUTH_TOKEN_KEY)
    }
  }

  async function login(username: string, password: string) {
    loading.value = true
    error.value = null
    try {
      const response = await apiService.login(username, password)
      setToken(response.access)
      return { success: true }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Ошибка авторизации'
      error.value = message
      return { success: false, error: message }
    } finally {
      loading.value = false
    }
  }

  async function signUp(username: string, password: string) {
    loading.value = true
    error.value = null
    try {
      await apiService.signUp(username, password)
      return await login(username, password)
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Ошибка регистрации'
      error.value = message
      return { success: false, error: message }
    } finally {
      loading.value = false
    }
  }

  function logout() {
    setToken(null)
  }

  return {
    token,
    loading,
    error,
    isAuthenticated,
    login,
    signUp,
    logout,
  }
})


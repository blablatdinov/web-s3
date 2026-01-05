<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useBucketStore } from '@/stores/bucket'

const router = useRouter()
const authStore = useAuthStore()
const bucketStore = useBucketStore()

function handleLogout() {
  authStore.logout()
  bucketStore.setSelectedBucket(null)
  router.push('/login')
}

function goToFiles() {
  if (!bucketStore.selectedBucketId) {
    router.push('/buckets')
  } else {
    router.push('/files')
  }
}

onMounted(() => {
  // Если бакет не выбран, перенаправляем на выбор бакета
  if (!bucketStore.selectedBucketId) {
    router.push('/buckets')
  }
})
</script>

<template>
  <div class="home-container">
    <div class="header">
      <h1>S3 Storage</h1>
      <button @click="handleLogout" class="logout-button">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M6 14H3C2.44772 14 2 13.5523 2 13V3C2 2.44772 2.44772 2 3 2H6M10 12L14 8M14 8L10 4M14 8H6"/>
        </svg>
        Выйти
      </button>
    </div>
    <div class="content">
      <div class="welcome-section">
        <svg width="64" height="64" viewBox="0 0 64 64" fill="none" stroke="currentColor" stroke-width="1.5" opacity="0.1">
          <path d="M20 16H44C46.2091 16 48 17.7909 48 20V48C48 50.2091 46.2091 52 44 52H20C17.7909 52 16 50.2091 16 48V20C16 17.7909 17.7909 16 20 16Z"/>
          <path d="M24 24H40M24 32H40M24 40H32"/>
        </svg>
        <h2>Добро пожаловать</h2>
        <p>Вы успешно авторизованы в системе управления файлами S3.</p>
      </div>
      <div class="actions">
        <button @click="goToFiles" class="action-button">
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 5C3 4.44772 3.44772 4 4 4H8.58579C8.851 4 9.10536 4.10536 9.29289 4.29289L11.7071 6.70711C11.8946 6.89464 12.149 7 12.4142 7H16C16.5523 7 17 7.44772 17 8V15C17 15.5523 16.5523 16 16 16H4C3.44772 16 3 15.5523 3 15V5Z"/>
          </svg>
          Просмотр файлов
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.home-container {
  min-height: 100vh;
  background: #fafafa;
}

.header {
  background: white;
  border-bottom: 1px solid #e5e7eb;
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

h1 {
  margin: 0;
  color: #111827;
  font-size: 20px;
  font-weight: 600;
}

.logout-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: none;
  color: #6b7280;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
}

.logout-button:hover {
  background-color: #f9fafb;
  border-color: #d1d5db;
  color: #111827;
}

.logout-button svg {
  flex-shrink: 0;
}

.content {
  max-width: 600px;
  margin: 0 auto;
  padding: 80px 24px;
}

.welcome-section {
  text-align: center;
  margin-bottom: 48px;
}

.welcome-section svg {
  margin-bottom: 24px;
  color: #3b82f6;
}

.welcome-section h2 {
  margin: 0 0 12px 0;
  color: #111827;
  font-size: 28px;
  font-weight: 600;
}

.welcome-section p {
  margin: 0;
  color: #6b7280;
  font-size: 16px;
  line-height: 1.5;
}

.actions {
  display: flex;
  justify-content: center;
}

.action-button {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 24px;
  background: #111827;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
}

.action-button:hover {
  background: #1f2937;
}

.action-button:active {
  background: #111827;
}

.action-button svg {
  flex-shrink: 0;
}

@media (max-width: 768px) {
  .content {
    padding: 48px 24px;
  }

  .welcome-section h2 {
    font-size: 24px;
  }
}
</style>


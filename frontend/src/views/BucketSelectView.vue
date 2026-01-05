<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useBucketStore } from '@/stores/bucket'

const router = useRouter()
const authStore = useAuthStore()
const bucketStore = useBucketStore()

const showCreateForm = ref(false)
const creating = ref(false)
const createError = ref<string | null>(null)

const formData = ref({
  bucket_name: '',
  access_key_id: '',
  secret_access_key: '',
  region: 'us-east-1',
  endpoint: '',
})

function handleLogout() {
  authStore.logout()
  bucketStore.setSelectedBucket(null)
  router.push('/login')
}

async function selectBucket(bucketId: number) {
  bucketStore.setSelectedBucket(bucketId)
  router.push('/files')
}

async function handleCreateBucket() {
  if (!formData.value.bucket_name || !formData.value.access_key_id || !formData.value.secret_access_key) {
    createError.value = 'Заполните все обязательные поля'
    return
  }

  creating.value = true
  createError.value = null

  const result = await bucketStore.createBucket({
    bucket_name: formData.value.bucket_name,
    access_key_id: formData.value.access_key_id,
    secret_access_key: formData.value.secret_access_key,
    region: formData.value.region,
    endpoint: formData.value.endpoint || undefined,
  })

  if (result.success) {
    showCreateForm.value = false
    formData.value = {
      bucket_name: '',
      access_key_id: '',
      secret_access_key: '',
      region: 'us-east-1',
      endpoint: '',
    }
    router.push('/files')
  } else {
    createError.value = result.error || 'Ошибка создания бакета'
  }

  creating.value = false
}

function toggleCreateForm() {
  showCreateForm.value = !showCreateForm.value
  createError.value = null
}

onMounted(async () => {
  await bucketStore.loadBuckets()
})
</script>

<template>
  <div class="bucket-select-container">
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
      <div class="title-section">
        <h2>Выберите бакет</h2>
        <p>Выберите бакет для работы с файлами или создайте новый</p>
      </div>

      <div v-if="bucketStore.error" class="error-message">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="8" cy="8" r="7"/>
          <path d="M8 4V8M8 12H8.01"/>
        </svg>
        {{ bucketStore.error }}
      </div>

      <div v-if="bucketStore.loading" class="loading">
        <div class="spinner"></div>
        <span>Загрузка бакетов...</span>
      </div>

      <div v-else class="buckets-section">
        <div v-if="bucketStore.buckets.length === 0 && !showCreateForm" class="empty-state">
          <svg width="64" height="64" viewBox="0 0 64 64" fill="none" stroke="currentColor" stroke-width="1.5" opacity="0.3">
            <path d="M20 16H44C46.2091 16 48 17.7909 48 20V48C48 50.2091 46.2091 52 44 52H20C17.7909 52 16 50.2091 16 48V20C16 17.7909 17.7909 16 20 16Z"/>
            <path d="M24 24H40M24 32H40M24 40H32"/>
          </svg>
          <p>У вас пока нет бакетов</p>
          <button @click="toggleCreateForm" class="create-button">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M8 3V13M3 8H13"/>
            </svg>
            Создать бакет
          </button>
        </div>

        <div v-else class="buckets-list">
          <button
            v-for="bucket in bucketStore.buckets"
            :key="bucket.bucket_id"
            @click="selectBucket(bucket.bucket_id)"
            class="bucket-item"
            :class="{ selected: bucketStore.selectedBucketId === bucket.bucket_id }"
          >
            <div class="bucket-info">
              <div class="bucket-name">{{ bucket.bucket_name }}</div>
              <div class="bucket-details">
                <span>Регион: {{ bucket.region }}</span>
                <span v-if="bucket.endpoint">Endpoint: {{ bucket.endpoint }}</span>
              </div>
            </div>
            <svg width="20" height="20" viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M7 14L13 8M13 8H9M13 8V12"/>
            </svg>
          </button>
        </div>

        <div class="actions">
          <button v-if="!showCreateForm" @click="toggleCreateForm" class="action-button secondary">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M8 3V13M3 8H13"/>
            </svg>
            Создать новый бакет
          </button>
        </div>

        <div v-if="showCreateForm" class="create-form">
          <h3>Создать новый бакет</h3>
          <div v-if="createError" class="error-message">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="8" cy="8" r="7"/>
              <path d="M8 4V8M8 12H8.01"/>
            </svg>
            {{ createError }}
          </div>
          <div class="form-group">
            <label for="bucket_name">Название бакета *</label>
            <input
              id="bucket_name"
              v-model="formData.bucket_name"
              type="text"
              placeholder="my-bucket"
              required
            />
          </div>
          <div class="form-group">
            <label for="access_key_id">Access Key ID *</label>
            <input
              id="access_key_id"
              v-model="formData.access_key_id"
              type="text"
              placeholder="AKIAIOSFODNN7EXAMPLE"
              required
            />
          </div>
          <div class="form-group">
            <label for="secret_access_key">Secret Access Key *</label>
            <input
              id="secret_access_key"
              v-model="formData.secret_access_key"
              type="password"
              placeholder="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
              required
            />
          </div>
          <div class="form-group">
            <label for="region">Регион</label>
            <input
              id="region"
              v-model="formData.region"
              type="text"
              placeholder="us-east-1"
            />
          </div>
          <div class="form-group">
            <label for="endpoint">Endpoint (опционально)</label>
            <input
              id="endpoint"
              v-model="formData.endpoint"
              type="text"
              placeholder="https://s3.example.com"
            />
          </div>
          <div class="form-actions">
            <button @click="handleCreateBucket" :disabled="creating" class="action-button">
              <svg v-if="creating" class="spinner-icon" width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="8" cy="8" r="6" stroke-dasharray="37.7" stroke-dashoffset="9.4"/>
              </svg>
              <svg v-else width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M8 3V13M3 8H13"/>
              </svg>
              {{ creating ? 'Создание...' : 'Создать' }}
            </button>
            <button @click="toggleCreateForm" :disabled="creating" class="action-button secondary">
              Отмена
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.bucket-select-container {
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

.content {
  max-width: 800px;
  margin: 0 auto;
  padding: 48px 24px;
}

.title-section {
  text-align: center;
  margin-bottom: 32px;
}

.title-section h2 {
  margin: 0 0 8px 0;
  color: #111827;
  font-size: 28px;
  font-weight: 600;
}

.title-section p {
  margin: 0;
  color: #6b7280;
  font-size: 16px;
}

.error-message {
  margin-bottom: 24px;
  padding: 12px 16px;
  background: #fef2f2;
  color: #dc2626;
  border-radius: 6px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  border: 1px solid #fecaca;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 80px 20px;
  color: #6b7280;
  font-size: 14px;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #e5e7eb;
  border-top-color: #111827;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.buckets-section {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 24px;
  padding: 80px 20px;
  color: #9ca3af;
}

.empty-state svg {
  flex-shrink: 0;
}

.empty-state p {
  font-size: 16px;
  margin: 0;
  color: #6b7280;
}

.create-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: #111827;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
}

.create-button:hover {
  background: #1f2937;
}

.buckets-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.bucket-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: white;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s;
  text-align: left;
}

.bucket-item:hover {
  border-color: #3b82f6;
  background: #f9fafb;
}

.bucket-item.selected {
  border-color: #3b82f6;
  background: #eff6ff;
}

.bucket-info {
  flex: 1;
  min-width: 0;
}

.bucket-name {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
  margin-bottom: 4px;
}

.bucket-details {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  font-size: 14px;
  color: #6b7280;
}

.bucket-item svg {
  flex-shrink: 0;
  color: #6b7280;
}

.bucket-item:hover svg,
.bucket-item.selected svg {
  color: #3b82f6;
}

.actions {
  display: flex;
  justify-content: center;
}

.action-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: #111827;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
}

.action-button:hover:not(:disabled) {
  background: #1f2937;
}

.action-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.action-button.secondary {
  background: white;
  color: #111827;
  border: 1px solid #e5e7eb;
}

.action-button.secondary:hover:not(:disabled) {
  background: #f9fafb;
  border-color: #d1d5db;
}

.create-form {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 24px;
  margin-top: 24px;
}

.create-form h3 {
  margin: 0 0 20px 0;
  color: #111827;
  font-size: 20px;
  font-weight: 600;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #374151;
  font-size: 14px;
  font-weight: 500;
}

.form-group input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.15s;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 24px;
}

.spinner-icon {
  animation: spin 0.8s linear infinite;
}

@media (max-width: 768px) {
  .content {
    padding: 32px 16px;
  }

  .bucket-item {
    padding: 12px 16px;
  }

  .bucket-details {
    flex-direction: column;
    gap: 4px;
  }

  .create-form {
    padding: 20px;
  }
}
</style>


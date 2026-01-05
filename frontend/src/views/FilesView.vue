<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiService } from '@/services/api'
import { useBucketStore } from '@/stores/bucket'
import type { FilesResponse } from '@/services/api'

const route = useRoute()
const router = useRouter()
const bucketStore = useBucketStore()

const files = ref<string[]>([])
const directories = ref<string[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const downloadingFiles = ref<Set<string>>(new Set())

const currentPath = computed(() => {
  return (route.query.path as string) || ''
})

const pathParts = computed(() => {
  if (!currentPath.value) return []
  const parts = currentPath.value.split('/').filter(Boolean)
  return parts.map((part, index) => ({
    name: part,
    path: parts.slice(0, index + 1).join('/'),
  }))
})

function getItemName(fullPath: string): string {
  const cleanPath = fullPath.replace(/\/$/, '')
  const parts = cleanPath.split('/').filter(Boolean)
  return parts[parts.length - 1] || cleanPath
}

async function loadFiles(path?: string) {
  if (!bucketStore.selectedBucketId) {
    router.push('/buckets')
    return
  }

  loading.value = true
  error.value = null
  try {
    const response: FilesResponse = await apiService.getFiles(bucketStore.selectedBucketId, path)
    files.value = response.files
    directories.value = response.directories
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Ошибка загрузки файлов'
    error.value = message
  } finally {
    loading.value = false
  }
}

function navigateToDirectory(dirPath: string) {
  router.push({ query: { path: dirPath } })
}

function navigateToPath(path: string) {
  router.push({ query: { path } })
}

function goToRoot() {
  router.push({ query: {} })
}

function goUp() {
  if (!currentPath.value) return
  const parts = currentPath.value.split('/').filter(Boolean)
  parts.pop()
  const newPath = parts.join('/')
  if (newPath) {
    router.push({ query: { path: newPath } })
  } else {
    goToRoot()
  }
}

onMounted(() => {
  loadFiles(currentPath.value || undefined)
})

watch(
  () => route.query.path,
  (newPath) => {
    loadFiles((newPath as string) || undefined)
  }
)

async function downloadFile(filePath: string) {
  if (!bucketStore.selectedBucketId) {
    router.push('/buckets')
    return
  }

  if (downloadingFiles.value.has(filePath)) {
    return
  }

  downloadingFiles.value.add(filePath)
  try {
    await apiService.downloadFile(bucketStore.selectedBucketId, filePath)
  } catch (err) {
    const message = err instanceof Error ? err.message : 'Ошибка скачивания файла'
    error.value = message
    // Убираем сообщение об ошибке через 5 секунд
    setTimeout(() => {
      error.value = null
    }, 5000)
  } finally {
    downloadingFiles.value.delete(filePath)
  }
}

function isDownloading(filePath: string): boolean {
  return downloadingFiles.value.has(filePath)
}
</script>

<template>
  <div class="files-container">
    <div class="header">
      <div class="header-left">
        <h1>Файлы</h1>
        <div v-if="bucketStore.selectedBucket" class="bucket-info">
          <span class="bucket-name">{{ bucketStore.selectedBucket.bucket_name }}</span>
        </div>
        <nav class="breadcrumb">
          <button @click="goToRoot" class="breadcrumb-item">Корень</button>
          <span v-if="pathParts.length > 0" class="breadcrumb-separator">/</span>
          <template v-for="(part, index) in pathParts" :key="index">
            <button @click="navigateToPath(part.path)" class="breadcrumb-item">
              {{ part.name }}
            </button>
            <span v-if="index < pathParts.length - 1" class="breadcrumb-separator">/</span>
          </template>
        </nav>
      </div>
      <div class="header-actions">
        <button
          v-if="currentPath"
          @click="goUp"
          class="icon-button"
          title="На уровень вверх"
        >
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M10 15L5 10L10 5M5 10H15"/>
          </svg>
        </button>
        <button
          @click="router.push('/buckets')"
          class="icon-button"
          title="Выбрать бакет"
        >
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 5C3 4.44772 3.44772 4 4 4H8.58579C8.851 4 9.10536 4.10536 9.29289 4.29289L11.7071 6.70711C11.8946 6.89464 12.149 7 12.4142 7H16C16.5523 7 17 7.44772 17 8V15C17 15.5523 16.5523 16 16 16H4C3.44772 16 3 15.5523 3 15V5Z"/>
          </svg>
        </button>
        <router-link to="/" class="icon-button" title="Главная">
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 10L10 3L17 10M5 10V16H8V13H12V16H15V10"/>
          </svg>
        </router-link>
      </div>
    </div>

    <div v-if="error" class="error-message">
      <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="8" cy="8" r="7"/>
        <path d="M8 4V8M8 12H8.01"/>
      </svg>
      {{ error }}
    </div>

    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <span>Загрузка...</span>
    </div>

    <div v-else class="files-content">
      <div v-if="directories !== null && directories.length === 0 && files.length === 0" class="empty-state">
        <svg width="64" height="64" viewBox="0 0 64 64" fill="none" stroke="currentColor" stroke-width="1.5" opacity="0.3">
          <path d="M20 16H44C46.2091 16 48 17.7909 48 20V48C48 50.2091 46.2091 52 44 52H20C17.7909 52 16 50.2091 16 48V20C16 17.7909 17.7909 16 20 16Z"/>
          <path d="M24 24H40M24 32H40M24 40H32"/>
        </svg>
        <p>Директория пуста</p>
      </div>

      <div v-if="directories !== null && directories.length > 0" class="directories-section">
        <div class="items-list">
          <button
            v-for="dir in directories"
            :key="dir"
            @click="navigateToDirectory(dir)"
            class="item directory-item"
          >
            <svg class="item-icon" width="20" height="20" viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M3 5C3 4.44772 3.44772 4 4 4H8.58579C8.851 4 9.10536 4.10536 9.29289 4.29289L11.7071 6.70711C11.8946 6.89464 12.149 7 12.4142 7H16C16.5523 7 17 7.44772 17 8V15C17 15.5523 16.5523 16 16 16H4C3.44772 16 3 15.5523 3 15V5Z"/>
            </svg>
            <span class="item-name">{{ getItemName(dir) }}</span>
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2" opacity="0.4">
              <path d="M6 12L10 8L6 4"/>
            </svg>
          </button>
        </div>
      </div>

      <div v-if="files.length > 0" class="files-section">
        <div class="items-list">
          <div
            v-for="file in files"
            :key="file"
            class="item file-item"
            @dblclick="downloadFile(file)"
            :title="'Двойной клик для скачивания'"
          >
            <svg class="item-icon" width="20" height="20" viewBox="0 0 20 20" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M5 3C4.44772 3 4 3.44772 4 4V16C4 16.5523 4.44772 17 5 17H15C15.5523 17 16 16.5523 16 16V7.41421C16 7.149 15.8946 6.89464 15.7071 6.70711L12.2929 3.29289C12.1054 3.10536 11.851 3 11.5858 3H5Z"/>
              <path d="M12 3V7H16"/>
            </svg>
            <span class="item-name">{{ getItemName(file) }}</span>
            <button
              @click.stop="downloadFile(file)"
              :disabled="isDownloading(file)"
              class="download-button"
              :title="isDownloading(file) ? 'Скачивание...' : 'Скачать файл'"
            >
              <svg v-if="isDownloading(file)" class="spinner-icon" width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="8" cy="8" r="6" stroke-dasharray="37.7" stroke-dashoffset="9.4"/>
              </svg>
              <svg v-else width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M8 12V2M8 12L4 8M8 12L12 8M2 14H14"/>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.files-container {
  min-height: 100vh;
  background: #fafafa;
  padding: 0;
}

.header {
  background: white;
  border-bottom: 1px solid #e5e7eb;
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 10;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 24px;
  flex: 1;
  min-width: 0;
}

.bucket-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 12px;
  background: #f3f4f6;
  border-radius: 6px;
  font-size: 14px;
}

.bucket-name {
  color: #6b7280;
  font-weight: 500;
}

h1 {
  margin: 0;
  color: #111827;
  font-size: 20px;
  font-weight: 600;
  white-space: nowrap;
}

.breadcrumb {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
  min-width: 0;
}

.breadcrumb-item {
  background: none;
  border: none;
  color: #6b7280;
  cursor: pointer;
  font-size: 14px;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.15s;
  white-space: nowrap;
}

.breadcrumb-item:hover {
  background-color: #f3f4f6;
  color: #111827;
}

.breadcrumb-separator {
  color: #d1d5db;
  font-size: 14px;
}

.header-actions {
  display: flex;
  gap: 8px;
  margin-left: 16px;
}

.icon-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: none;
  border: none;
  border-radius: 6px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.15s;
  flex-shrink: 0;
}

.icon-button:hover {
  background-color: #f3f4f6;
  color: #111827;
}

.icon-button svg {
  display: block;
}

.error-message {
  margin: 16px 24px;
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

.error-message svg {
  flex-shrink: 0;
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

.files-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 24px;
}

.directories-section,
.files-section {
  margin-bottom: 32px;
}

.directories-section:last-child,
.files-section:last-child {
  margin-bottom: 0;
}

.items-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 2px;
  background: #f3f4f6;
  border-radius: 8px;
  padding: 2px;
  overflow: hidden;
}

.item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 6px;
  transition: all 0.15s;
  background: white;
}

.directory-item {
  background: none;
  border: none;
  cursor: pointer;
  text-align: left;
  width: 100%;
}

.directory-item:hover {
  background-color: #f9fafb;
}

.directory-item:active {
  background-color: #f3f4f6;
}

.file-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
  user-select: none;
}

.file-item:hover {
  background-color: #f9fafb;
}

.item-icon {
  flex-shrink: 0;
  color: #6b7280;
}

.directory-item .item-icon {
  color: #3b82f6;
}

.item-name {
  flex: 1;
  color: #111827;
  word-break: break-word;
  font-size: 14px;
  font-weight: 400;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 120px 20px;
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

.download-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: none;
  border: none;
  border-radius: 6px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.15s;
  flex-shrink: 0;
  padding: 0;
}

.download-button:hover:not(:disabled) {
  background-color: #f3f4f6;
  color: #111827;
}

.download-button:active:not(:disabled) {
  background-color: #e5e7eb;
}

.download-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.download-button svg {
  display: block;
}

.spinner-icon {
  animation: spin 0.8s linear infinite;
}

@media (max-width: 768px) {
  .header {
    padding: 12px 16px;
    flex-wrap: wrap;
  }

  .header-left {
    width: 100%;
    margin-bottom: 8px;
  }

  .breadcrumb {
    width: 100%;
  }

  .files-content {
    padding: 16px;
  }

  .items-list {
    grid-template-columns: 1fr;
  }
}
</style>


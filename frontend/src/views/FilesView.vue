<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiService } from '@/services/api'
import type { FilesResponse } from '@/services/api'

const route = useRoute()
const router = useRouter()

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
  loading.value = true
  error.value = null
  try {
    const response: FilesResponse = await apiService.getFiles(path)
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
  if (downloadingFiles.value.has(filePath)) {
    return
  }

  downloadingFiles.value.add(filePath)
  try {
    await apiService.downloadFile(filePath)
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
      <h1>Файлы</h1>
      <div class="header-actions">
        <router-link to="/" class="home-link">Главная</router-link>
      </div>
    </div>

    <div class="breadcrumb">
      <button @click="goToRoot" class="breadcrumb-item">Корень</button>
      <span v-if="pathParts.length > 0" class="breadcrumb-separator">/</span>
      <template v-for="(part, index) in pathParts" :key="index">
        <button @click="navigateToPath(part.path)" class="breadcrumb-item">
          {{ part.name }}
        </button>
        <span v-if="index < pathParts.length - 1" class="breadcrumb-separator">/</span>
      </template>
      <button
        v-if="currentPath"
        @click="goUp"
        class="breadcrumb-item breadcrumb-up"
        title="На уровень вверх"
      >
        ↑
      </button>
    </div>

    <div v-if="error" class="error-message">{{ error }}</div>

    <div v-if="loading" class="loading">Загрузка...</div>

    <div v-else class="files-content">
      <div v-if="directories !== null && directories.length === 0 && files.length === 0" class="empty-state">
        <p>Директория пуста</p>
      </div>

      <div v-if="directories !== null && directories.length > 0" class="directories-section">
        <h2>Директории</h2>
        <div class="items-list">
          <button
            v-for="dir in directories"
            :key="dir"
            @click="navigateToDirectory(dir)"
            class="item directory-item"
          >
            <span class="item-icon">📁</span>
            <span class="item-name">{{ getItemName(dir) }}</span>
          </button>
        </div>
      </div>

      <div v-if="files.length > 0" class="files-section">
        <h2>Файлы</h2>
        <div class="items-list">
          <div
            v-for="file in files"
            :key="file"
            class="item file-item"
            @dblclick="downloadFile(file)"
            :title="'Двойной клик для скачивания'"
          >
            <span class="item-icon">📄</span>
            <span class="item-name">{{ getItemName(file) }}</span>
            <button
              @click.stop="downloadFile(file)"
              :disabled="isDownloading(file)"
              class="download-button"
              :title="isDownloading(file) ? 'Скачивание...' : 'Скачать файл'"
            >
              <span v-if="isDownloading(file)" class="download-spinner">⏳</span>
              <span v-else>⬇️</span>
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto 20px;
  padding: 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

h1 {
  margin: 0;
  color: #333;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.home-link {
  padding: 10px 20px;
  background: #667eea;
  color: white;
  text-decoration: none;
  border-radius: 8px;
  font-weight: 500;
  transition: opacity 0.3s;
}

.home-link:hover {
  opacity: 0.9;
}

.breadcrumb {
  max-width: 1200px;
  margin: 0 auto 20px;
  padding: 15px 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 5px;
}

.breadcrumb-item {
  background: none;
  border: none;
  color: #667eea;
  cursor: pointer;
  font-size: 14px;
  padding: 5px 10px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.breadcrumb-item:hover {
  background-color: #f0f0f0;
}

.breadcrumb-separator {
  color: #999;
}

.breadcrumb-up {
  margin-left: auto;
  font-size: 18px;
  font-weight: bold;
}

.files-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 30px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.directories-section,
.files-section {
  margin-bottom: 30px;
}

.directories-section:last-child,
.files-section:last-child {
  margin-bottom: 0;
}

h2 {
  margin: 0 0 15px 0;
  color: #333;
  font-size: 20px;
}

.items-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 10px;
}

.item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px;
  border-radius: 8px;
  transition: background-color 0.2s;
}

.directory-item {
  background: none;
  border: none;
  cursor: pointer;
  text-align: left;
  width: 100%;
}

.directory-item:hover {
  background-color: #f5f5f5;
}

.file-item {
  background-color: #fafafa;
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
  user-select: none;
}

.file-item:hover {
  background-color: #f0f0f0;
}

.item-icon {
  font-size: 20px;
}

.item-name {
  flex: 1;
  color: #333;
  word-break: break-word;
  font-size: 14px;
}

.error-message {
  max-width: 1200px;
  margin: 0 auto 20px;
  padding: 15px;
  background: #fee;
  color: #e74c3c;
  border-radius: 8px;
  text-align: center;
}

.loading {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px;
  background: white;
  border-radius: 12px;
  text-align: center;
  color: #666;
  font-size: 18px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #999;
}

.empty-state p {
  font-size: 18px;
  margin: 0;
}

.download-button {
  background: #667eea;
  color: white;
  border: none;
  border-radius: 6px;
  padding: 6px 12px;
  cursor: pointer;
  font-size: 16px;
  transition: opacity 0.2s, background-color 0.2s;
  flex-shrink: 0;
}

.download-button:hover:not(:disabled) {
  opacity: 0.9;
  background: #5568d3;
}

.download-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.download-spinner {
  display: inline-block;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>


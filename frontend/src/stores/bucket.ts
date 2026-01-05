import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { apiService, type Bucket } from '@/services/api'

const SELECTED_BUCKET_KEY = 'selected_bucket_id'

export const useBucketStore = defineStore('bucket', () => {
  const getInitialBucketId = (): number | null => {
    const stored = localStorage.getItem(SELECTED_BUCKET_KEY)
    return stored ? parseInt(stored, 10) : null
  }
  const selectedBucketId = ref<number | null>(getInitialBucketId())
  const buckets = ref<Bucket[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const selectedBucket = computed(() => {
    if (!selectedBucketId.value) return null
    return buckets.value.find((b) => b.bucket_id === selectedBucketId.value) || null
  })

  function setSelectedBucket(bucketId: number | null) {
    selectedBucketId.value = bucketId
    if (bucketId) {
      localStorage.setItem(SELECTED_BUCKET_KEY, bucketId.toString())
    } else {
      localStorage.removeItem(SELECTED_BUCKET_KEY)
    }
  }

  async function loadBuckets() {
    loading.value = true
    error.value = null
    try {
      const response = await apiService.getBuckets()
      buckets.value = response.buckets
      return { success: true }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Ошибка загрузки бакетов'
      error.value = message
      return { success: false, error: message }
    } finally {
      loading.value = false
    }
  }

  async function createBucket(data: {
    bucket_name: string
    access_key_id: string
    secret_access_key: string
    region: string
    endpoint?: string
  }) {
    loading.value = true
    error.value = null
    try {
      const response = await apiService.createBucket(data)
      await loadBuckets()
      setSelectedBucket(response.bucket_id)
      return { success: true, bucketId: response.bucket_id }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Ошибка создания бакета'
      error.value = message
      return { success: false, error: message }
    } finally {
      loading.value = false
    }
  }

  return {
    selectedBucketId,
    selectedBucket,
    buckets,
    loading,
    error,
    setSelectedBucket,
    loadBuckets,
    createBucket,
  }
})


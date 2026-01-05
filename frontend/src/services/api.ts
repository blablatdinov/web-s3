const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

export interface AuthResponse {
  access: string
}

export interface SignUpResponse {
  user_id: number
  username: string
}

export interface AuthError {
  error?: string
  details?: string
}

export interface FilesResponse {
  files: string[] | null
  directories: string[] | null
}

export interface Bucket {
  bucket_id: number
  user_id: number
  bucket_name: string
  access_key_id: string
  region: string
  endpoint: string | null
  created_at: string
  updated_at: string
}

export interface BucketsResponse {
  buckets: Bucket[]
}

export interface CreateBucketRequest {
  bucket_name: string
  access_key_id: string
  secret_access_key: string
  region: string
  endpoint?: string
}

export interface CreateBucketResponse {
  bucket_id: number
  bucket_name: string
}

class ApiService {
  private baseUrl: string

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string>),
    }

    const token = localStorage.getItem('auth_token')
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    const response = await fetch(url, {
      ...options,
      headers,
    })

    if (!response.ok) {
      const error: AuthError = await response.json().catch(() => ({}))
      throw new Error(error.error || error.details || 'Request failed')
    }

    return response.json()
  }

  async login(username: string, password: string): Promise<AuthResponse> {
    console.log(JSON.stringify({ username, password }))
    return this.request<AuthResponse>('/users/auth', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    })
  }

  async signUp(username: string, password: string): Promise<SignUpResponse> {
    return this.request<SignUpResponse>('/users/sign-up', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    })
  }

  async getFiles(bucketId: number, path?: string): Promise<FilesResponse> {
    let url = `/files?bucket_id=${bucketId}`
    if (path) {
      url += `&path=${encodeURIComponent(path)}`
    }
    const files = this.request<FilesResponse>(url, {
      method: 'GET',
    })
    console.log(files)
    return files
  }

  async getBuckets(): Promise<BucketsResponse> {
    return this.request<BucketsResponse>('/buckets', {
      method: 'GET',
    })
  }

  async createBucket(data: CreateBucketRequest): Promise<CreateBucketResponse> {
    return this.request<CreateBucketResponse>('/buckets', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async downloadFile(bucketId: number, filePath: string): Promise<void> {
    const token = localStorage.getItem('auth_token')
    let url = `${this.baseUrl}/files/${encodeURIComponent(filePath)}/download?bucket_id=${bucketId}`
    if (token) {
      url += `&token=${encodeURIComponent(token)}`
    }

    const headers: Record<string, string> = {}
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    const response = await fetch(url, {
      method: 'GET',
      headers,
    })

    if (!response.ok) {
      const error: AuthError = await response.json().catch(() => ({}))
      throw new Error(error.error || error.details || 'Ошибка скачивания файла')
    }
    const contentDisposition = response.headers.get('Content-Disposition')
    let fileName = filePath.split('/').pop() || 'file'
    if (contentDisposition) {
      const utf8Match = contentDisposition.match(/filename\*=UTF-8''([^;]+)/i)
      if (utf8Match && utf8Match[1]) {
        try {
          fileName = decodeURIComponent(utf8Match[1])
        } catch {
        }
      } else {
        const filenameMatch = contentDisposition.match(/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/i)
        if (filenameMatch && filenameMatch[1]) {
          let extractedName = filenameMatch[1].replace(/^["']|["']$/g, '')
          try {
            fileName = decodeURIComponent(extractedName)
          } catch {
            fileName = extractedName
          }
        }
      }
    }

    try {
      const decoded = decodeURIComponent(fileName)
      if (decoded !== fileName) {
        fileName = decoded
      }
    } catch {
    }

    const blob = await response.blob()
    const downloadUrl = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = fileName
    link.style.display = 'none'
    document.body.appendChild(link)
    
    link.click()
    
    setTimeout(() => {
      if (link.parentNode) {
        document.body.removeChild(link)
      }
      window.URL.revokeObjectURL(downloadUrl)
    }, 100)
  }
}

export const apiService = new ApiService(API_BASE_URL)


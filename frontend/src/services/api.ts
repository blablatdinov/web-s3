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

  async getFiles(path?: string): Promise<FilesResponse> {
    const url = path ? `/files?path=${encodeURIComponent(path)}` : '/files'
    const files = this.request<FilesResponse>(url, {
      method: 'GET',
    })
    console.log(files)
    return files
  }

  async downloadFile(filePath: string): Promise<void> {
    const url = `${this.baseUrl}/files/${encodeURIComponent(filePath)}/download`
    const token = localStorage.getItem('auth_token')

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

    // Получаем имя файла из заголовка Content-Disposition или из пути
    const contentDisposition = response.headers.get('Content-Disposition')
    let fileName = filePath.split('/').pop() || 'file'
    if (contentDisposition) {
      const fileNameMatch = contentDisposition.match(/filename="(.+)"/)
      if (fileNameMatch && fileNameMatch[1]) {
        fileName = fileNameMatch[1]
      }
    }

    // Создаем blob и скачиваем файл
    const blob = await response.blob()
    const downloadUrl = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = fileName
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(downloadUrl)
  }
}

export const apiService = new ApiService(API_BASE_URL)


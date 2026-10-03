// Mock API implementation - will be replaced with real API calls

export interface HeadscaleSettings {
  apiKey: string
  theme: 'light' | 'dark'
  profileName: string
}

export interface UpdateSettingsRequest {
  apiKey?: string
  theme?: 'light' | 'dark'
  profileName?: string
}

export interface ChangePasswordRequest {
  currentPassword: string
  newPassword: string
}

export interface LoginRequest {
  password: string
  apiKey?: string
}

// Mock data store
let mockSettings: HeadscaleSettings = {
  apiKey: 'hs_mock_key_12345',
  theme: 'light',
  profileName: 'Test User',
}

let mockPassword = 'password123'

class HeadscaleAPI {
  private baseUrl: string
  private sessionToken: string | null = null

  constructor(baseUrl: string = '/api/v1') {
    this.baseUrl = baseUrl
  }

  setSessionToken(token: string | null) {
    this.sessionToken = token
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...(this.sessionToken && { Authorization: `Bearer ${this.sessionToken}` }),
      ...options.headers,
    }

    // Mock API - simulate network delay
    await new Promise((resolve) => setTimeout(resolve, 300))

    // Mock responses based on endpoint
    if (endpoint.includes('/headplane/login')) {
      const body = JSON.parse(options.body as string) as LoginRequest
      if (body.password === mockPassword) {
        return { success: true, token: 'mock-session-token-' + Date.now() } as T
      }
      throw new Error('Invalid password')
    }

    if (endpoint.includes('/headplane/settings') && options.method === 'GET') {
      return mockSettings as T
    }

    if (endpoint.includes('/headplane/settings') && options.method === 'POST') {
      const updates = JSON.parse(options.body as string) as UpdateSettingsRequest
      mockSettings = { ...mockSettings, ...updates }
      return { success: true } as T
    }

    if (endpoint.includes('/headplane/change-password')) {
      const body = JSON.parse(options.body as string) as ChangePasswordRequest
      if (body.currentPassword === mockPassword) {
        mockPassword = body.newPassword
        return { success: true } as T
      }
      throw new Error('Invalid current password')
    }

    throw new Error(`Unhandled mock endpoint: ${endpoint}`)
  }

  async login(data: LoginRequest): Promise<{ success: boolean; token: string }> {
    return this.request('/headplane/login', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async getSettings(): Promise<HeadscaleSettings> {
    return this.request('/headplane/settings')
  }

  async updateSettings(
    data: UpdateSettingsRequest
  ): Promise<{ success: boolean }> {
    return this.request('/headplane/settings', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async changePassword(
    data: ChangePasswordRequest
  ): Promise<{ success: boolean }> {
    return this.request('/headplane/change-password', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }
}

export const api = new HeadscaleAPI()

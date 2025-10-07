import { api } from './client'

export const DEMO_TENANT_ID = '00000000-0000-0000-0000-000000000001'

export interface LoginRequest {
  email: string
  password: string
  tenant_id: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  user: {
    id: string
    email: string
    firstName: string
    lastName: string
    roles: string[]
    permissions: string[]
  }
}

export const authApi = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post('/auth/login', data)
    return response.data
  },

  refresh: async (refreshToken: string): Promise<LoginResponse> => {
    const response = await api.post('/auth/refresh', {
      refresh_token: refreshToken,
    })
    return response.data
  },

  me: async (): Promise<LoginResponse> => {
    const response = await api.get('/auth/me')
    return response.data
  },
}

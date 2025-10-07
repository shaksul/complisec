import axios from 'axios'

// Extend Vite's ImportMetaEnv interface
// interface ImportMetaEnv {
//   readonly VITE_API_URL: string
// }

const API_BASE_URL = (import.meta as any).env.VITE_API_URL || 'http://localhost:8080/api'

export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json; charset=utf-8',
    'Accept': 'application/json; charset=utf-8',
  },
})

// Alias for backward compatibility
export const apiClient = api

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle UTF-8 encoding and token refresh
api.interceptors.response.use(
  (response) => {
    // Принудительно устанавливаем UTF-8 кодировку для ответов
    if (response.headers['content-type'] && 
        response.headers['content-type'].includes('application/json') &&
        !response.headers['content-type'].includes('charset=utf-8')) {
      response.headers['content-type'] = response.headers['content-type'].replace(
        'application/json', 
        'application/json; charset=utf-8'
      )
    }
    return response
  },
  async (error) => {
    const originalRequest = error.config

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        const refreshToken = localStorage.getItem('refresh_token')
        if (!refreshToken) {
          // No refresh token available, redirect to login
          localStorage.removeItem('access_token')
          localStorage.removeItem('refresh_token')
          window.location.href = '/login'
          return Promise.reject(error)
        }

        const response = await api.post('/auth/refresh', {
          refresh_token: refreshToken,
        })
        
        const { access_token, refresh_token } = response.data
        localStorage.setItem('access_token', access_token)
        localStorage.setItem('refresh_token', refresh_token)
        
        originalRequest.headers.Authorization = `Bearer ${access_token}`
        return api(originalRequest)
      } catch (refreshError) {
        // Refresh failed, clear tokens and redirect to login
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        window.location.href = '/login'
        return Promise.reject(refreshError)
      }
    }

    return Promise.reject(error)
  }
)

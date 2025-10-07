import { apiClient } from './client'

export interface Tenant {
  id: string
  name: string
  domain?: string
  created_at: string
  updated_at: string
}

export interface CreateTenantDTO {
  name: string
  domain?: string
}

export interface UpdateTenantDTO {
  name: string
  domain?: string
}

export interface TenantListResponse {
  data: Tenant[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
    has_next: boolean
    has_prev: boolean
  }
}

export const tenantsApi = {
  // Получить список организаций
  getTenants: async (page: number = 1, pageSize: number = 20): Promise<TenantListResponse> => {
    const response = await apiClient.get(`/tenants?page=${page}&page_size=${pageSize}`)
    return response.data
  },

  // Получить организацию по ID
  getTenant: async (id: string): Promise<Tenant> => {
    const response = await apiClient.get(`/tenants/${id}`)
    return response.data
  },

  // Получить организацию по домену
  getTenantByDomain: async (domain: string): Promise<Tenant> => {
    const response = await apiClient.get(`/tenants/domain/${domain}`)
    return response.data
  },

  // Создать организацию
  createTenant: async (data: CreateTenantDTO): Promise<Tenant> => {
    const response = await apiClient.post('/tenants', data)
    return response.data
  },

  // Обновить организацию
  updateTenant: async (id: string, data: UpdateTenantDTO): Promise<Tenant> => {
    const response = await apiClient.put(`/tenants/${id}`, data)
    return response.data
  },

  // Удалить организацию
  deleteTenant: async (id: string): Promise<void> => {
    await apiClient.delete(`/tenants/${id}`)
  }
}

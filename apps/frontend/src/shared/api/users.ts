import { api as apiClient } from './client'

export interface User {
  id: string
  email: string
  first_name?: string
  last_name?: string
  is_active: boolean
  created_at?: string
  updated_at?: string
  roles?: string[]
  permissions?: string[]
}

export interface UserCatalog {
  id: string
  email: string
  first_name?: string
  last_name?: string
  is_active: boolean
  roles: string[]
  created_at: string
  updated_at: string
}

export interface UserDetail {
  id: string
  email: string
  first_name?: string
  last_name?: string
  is_active: boolean
  roles: string[]
  created_at: string
  updated_at: string
  last_login?: string
  stats: {
    documents_count: number
    risks_count: number
    incidents_count: number
    assets_count: number
    sessions_count: number
    login_count: number
    activity_score: number
  }
}

export interface UserActivity {
  id: string
  user_id: string
  action: string
  description: string
  ip_address?: string
  user_agent?: string
  created_at: string
  metadata?: Record<string, any>
}

export interface UserActivityStats {
  daily_activity: Array<{
    date: string
    count: number
  }>
  top_actions: Array<{
    action: string
    count: number
  }>
  login_history: Array<{
    ip_address: string
    user_agent: string
    created_at: string
    success: boolean
  }>
}

export interface PaginationResponse {
  page: number;
  page_size: number;
  total: number;
  total_pages: number;
  has_next: boolean;
  has_prev: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: PaginationResponse;
}

export interface UserCatalogParams {
  page?: number
  page_size?: number
  search?: string
  role?: string
  is_active?: boolean
  sort_by?: string
  sort_dir?: 'asc' | 'desc'
}

export const getUsers = async (): Promise<User[]> => {
  console.log('getUsers API called')
  const response = await apiClient.get('/users')
  console.log('getUsers API response:', response.data)
  return response.data.data || []
}

export const getUsersPaginated = async (page: number = 1, pageSize: number = 20): Promise<PaginatedResponse<User>> => {
  const response = await apiClient.get(`/users?page=${page}&page_size=${pageSize}`)
  return {
    data: response.data.data || [],
    pagination: response.data.pagination
  }
}

export const getUserCatalog = async (params: UserCatalogParams = {}): Promise<PaginatedResponse<UserCatalog>> => {
  const searchParams = new URLSearchParams()
  
  if (params.page) searchParams.append('page', params.page.toString())
  if (params.page_size) searchParams.append('page_size', params.page_size.toString())
  if (params.search) searchParams.append('search', params.search)
  if (params.role) searchParams.append('role', params.role)
  if (params.is_active !== undefined) searchParams.append('is_active', params.is_active.toString())
  if (params.sort_by) searchParams.append('sort_by', params.sort_by)
  if (params.sort_dir) searchParams.append('sort_dir', params.sort_dir)

  const response = await apiClient.get(`/users/catalog?${searchParams.toString()}`)
  return {
    data: response.data.data || [],
    pagination: response.data.pagination
  }
}

export const getUser = async (id: string): Promise<User> => {
  const response = await apiClient.get(`/users/${id}`)
  return response.data.data
}

export const getUserDetail = async (id: string): Promise<UserDetail> => {
  const response = await apiClient.get(`/users/${id}/detail`)
  return response.data.data
}

export const getUserActivity = async (id: string, page: number = 1, pageSize: number = 50): Promise<PaginatedResponse<UserActivity>> => {
  const response = await apiClient.get(`/users/${id}/activity?page=${page}&page_size=${pageSize}`)
  return {
    data: response.data.data || [],
    pagination: response.data.pagination
  }
}

export const getUserActivityStats = async (id: string): Promise<UserActivityStats> => {
  const response = await apiClient.get(`/users/${id}/activity/stats`)
  return response.data.data
}

export const createUser = async (userData: {
  email: string
  password: string
  first_name: string
  last_name: string
  department?: string
  role_ids: string[]
}): Promise<User> => {
  const response = await apiClient.post('/users', userData)
  return response.data.data
}

export const updateUser = async (id: string, userData: {
  first_name?: string
  last_name?: string
  is_active?: boolean
  role_ids?: string[]
}): Promise<void> => {
  await apiClient.put(`/users/${id}`, userData)
}

export const deleteUser = async (id: string): Promise<void> => {
  await apiClient.delete(`/users/${id}`)
}

// Экспортируем объект с API методами
export const usersApi = {
  list: getUsers,
  getUsers,
  getUsersPaginated,
  getUserCatalog,
  getUser,
  getUserDetail,
  getUserActivity,
  getUserActivityStats,
  createUser,
  updateUser,
  deleteUser
}


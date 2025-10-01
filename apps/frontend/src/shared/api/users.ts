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

export const getUsers = async (): Promise<User[]> => {
  const response = await apiClient.get('/users')
  return response.data.data || []
}

export const getUsersPaginated = async (page: number = 1, pageSize: number = 20): Promise<PaginatedResponse<User>> => {
  const response = await apiClient.get(`/users?page=${page}&page_size=${pageSize}`)
  return {
    data: response.data.data || [],
    pagination: response.data.pagination
  }
}

export const getUser = async (id: string): Promise<User> => {
  const response = await apiClient.get(`/users/${id}`)
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
  getUsers,
  getUsersPaginated,
  getUser,
  createUser,
  updateUser,
  deleteUser
}

